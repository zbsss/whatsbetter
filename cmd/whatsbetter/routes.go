package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/form/v4"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"github.com/zbsss/whatsbetter/internal/models"
	"github.com/zbsss/whatsbetter/internal/models/mocks"
	"github.com/zbsss/whatsbetter/internal/rating"
	"github.com/zbsss/whatsbetter/internal/validator"
	"github.com/zbsss/whatsbetter/ui"
	"github.com/zbsss/whatsbetter/ui/components"
	"github.com/zbsss/whatsbetter/ui/pages"
)

func (app *app) decodePostForm(r *http.Request, dst any) error {
	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecodeError *form.InvalidDecoderError
		if errors.As(err, &invalidDecodeError) {
			panic(err)
		}
		return err
	}

	return nil
}

func (app *app) home(w http.ResponseWriter, r *http.Request) {
	props := pages.IndexProps{
		Title:     "What's better?",
		Items:     mocks.MockItems,
		HTMXDebug: app.config.debug,
		FormData:  newFormData(r),
	}

	err := pages.Index(props).Render(context.Background(), w)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *app) createRestaurant(w http.ResponseWriter, r *http.Request) {
	var fd components.FormData
	err := app.decodePostForm(r, &fd)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	fd.CheckField(validator.NotBlank(fd.Name), "Name", "Name is required")
	fd.CheckField(validator.PositiveInt(fd.NumReviews), "NumReviews", "Number of reviews must be a positive integer")

	fd.CheckField(validator.IsFloat64(fd.Rating), "Rating", "Rating must be a decimal number between 0.0 and 5.0")

	ra, _ := strconv.ParseFloat(fd.Rating, 64)
	fd.CheckField(validator.InRangeFloat64(ra, 0.0, 5.0), "Rating", "Rating must be a decimal number between 0.0 and 5.0")

	if !fd.Valid() {
		cmp := components.Form(fd)
		err = cmp.Render(context.Background(), w)
		if err != nil {
			app.serverError(w, r, err)
			return
		}
		return
	}

	numReviews, _ := strconv.Atoi(fd.NumReviews)

	item := models.Item{
		Name:         fd.Name,
		Rating:       ra,
		NumOfReviews: numReviews,
	}
	item.CalculatedRating = rating.Calculate(item.Rating, item.NumOfReviews)

	cmp := components.Form(newFormData(r))
	err = cmp.Render(context.Background(), w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	cmp1 := components.ItemListOOB(item)
	err = cmp1.Render(context.Background(), w)
	if err != nil {
		app.serverError(w, r, err)
		return
	}
}

func newFormData(r *http.Request) components.FormData {
	return components.FormData{
		CSRFToken: nosurf.Token(r),
	}
}

func (app *app) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fs := http.FileServer(http.FS(ui.Files))
	router.Handler(http.MethodGet, "/static/*filepath", fs)

	router.HandlerFunc(http.MethodGet, "/healthz", health)

	dynamic := alice.New(noSurf)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodPost, "/restaurants", dynamic.ThenFunc(app.createRestaurant))

	// standard middleware for all requests
	standard := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return standard.Then(router)
}
