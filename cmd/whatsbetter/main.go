package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/form/v4"
)

type config struct {
	debug bool
	addr  string
}

type app struct {
	config      config
	logger      *slog.Logger
	formDecoder *form.Decoder
}

func main() {
	var cfg config
	flag.BoolVar(&cfg.debug, "debug", false, "enable debug mode")
	flag.StringVar(&cfg.addr, "addr", ":8080", "HTTP network address")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	app := &app{
		config:      cfg,
		logger:      logger,
		formDecoder: form.NewDecoder(),
	}

	srv := &http.Server{
		Addr:         cfg.addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  1 * time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	app.logger.Info("starting server", slog.String("addr", cfg.addr))

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				app.logger.Info("Server shut down gracefully")
			} else {
				app.logger.Error("Server shut down unexpectedly:", err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}
