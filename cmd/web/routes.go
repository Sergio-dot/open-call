package main

import (
	"github.com/Sergio-dot/open-call/internal/config"
	"github.com/Sergio-dot/open-call/internal/handlers"
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
	mux.Post("/user/signin", handlers.Repo.SignIn)
	mux.Post("/user/signup", handlers.Repo.SignUp)

	// routes protected by authentication middleware
	mux.Route("/", func(mux chi.Router) {
		mux.Use(Auth)

		mux.Get("/user/signout", handlers.Repo.SignOut)
		mux.Get("/user/update/{id}/do", handlers.Repo.UpdateUser)
		mux.Get("/dashboard", handlers.Repo.Dashboard)
		mux.Get("/room", handlers.Repo.Room)
	})

	// enable static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
