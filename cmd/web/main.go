package main

import (
	"encoding/gob"
	"fmt"
	"github.com/Sergio-dot/open-call/internal/config"
	"github.com/Sergio-dot/open-call/internal/driver"
	"github.com/Sergio-dot/open-call/internal/handlers"
	"github.com/Sergio-dot/open-call/internal/helpers"
	"github.com/Sergio-dot/open-call/internal/models"
	"github.com/Sergio-dot/open-call/internal/render"
	"github.com/alexedwards/scs/v2"
	"log"
	"net/http"
	"os"
	"time"
)

const port = ":8080"

var (
	app      config.AppConfig
	session  *scs.SessionManager
	infoLog  *log.Logger
	errorLog *log.Logger
)

// Main is the main application function
func main() {
	db, err := run()
	if err != nil {
		log.Fatal("error running application", err)
	}
	defer db.SQL.Close()

	fmt.Println(fmt.Sprintf("Starting application on port %s", port))
	srv := http.Server{
		Addr:    port,
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() (*driver.DB, error) {
	// values to put in the session
	gob.Register(models.User{})

	// true = Production, false = Development
	app.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// session management settings
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // set to 'true' in production (https)

	// store session to AppConfig
	app.Session = session

	// connect to database
	log.Println("Connecting to database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=opencall user=postgres password=root")
	if err != nil {
		log.Fatal("Could not reach database. Dying...")
	}
	log.Println("Connected to database")

	// create the template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("could not create template cache")
		return nil, err
	}

	// store template cache in AppConfig
	app.TemplateCache = tc

	// cache setting - set UseCache to 'true' in production
	app.UseCache = false

	// creates a new repository, giving access to AppConfig
	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	// give to the render package access to AppConfig
	render.NewRenderer(&app)

	// give to helpers package the access to AppConfig
	helpers.NewHelpers(&app)

	return db, nil
}
