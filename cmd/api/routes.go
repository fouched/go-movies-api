package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {
	// create router mux
	mux := chi.NewRouter()
	// if something goes wrong log error and don't panic/halt
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.Home)
	mux.Get("/authenticate", app.authenticate)
	mux.Get("/movies", app.AllMovies)

	return mux
}
