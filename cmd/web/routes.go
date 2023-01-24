package main

import (
	"github.com/Sergio-dot/open-call/pkg/config"
	"github.com/Sergio-dot/open-call/pkg/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	// create new router
	mux := chi.NewRouter()

	// set middlewares to use
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	// routes
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	return mux
}
