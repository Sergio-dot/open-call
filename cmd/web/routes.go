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
	mux.Get("/room", handlers.Repo.Room)

	// enable static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
