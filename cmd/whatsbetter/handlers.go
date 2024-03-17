package main

import "net/http"

func health(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func (app *app) createItem(w http.ResponseWriter, r *http.Request) {

}
