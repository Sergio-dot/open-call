package main

import (
	"fmt"
	"github.com/Sergio-dot/open-call/internal/config"
	"github.com/Sergio-dot/open-call/internal/handlers"
	"github.com/Sergio-dot/open-call/internal/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"time"
)

const port = ":8080"

var (
	app     config.AppConfig
	session *scs.SessionManager
)

// Main is the main application function
func main() {
	// true = Production, false = Development
	app.InProduction = false

	// session management settings
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // set to 'true' in production (https)

	// store session to AppConfig
	app.Session = session

	// create the template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("could not generate template cache")
	}

	// store template cache in AppConfig
	app.TemplateCache = tc

	// cache setting - set UseCache to 'true' in production
	app.UseCache = false

	// creates a new repository, giving access to AppConfig
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// give to render package access to the AppConfig
	render.NewTemplates(&app)

	fmt.Println(fmt.Sprintf("Starting application on port %s", port))
	srv := http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}
