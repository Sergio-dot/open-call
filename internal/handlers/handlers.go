package handlers

import (
	"github.com/Sergio-dot/open-call/internal/config"
	"github.com/Sergio-dot/open-call/internal/driver"
	"github.com/Sergio-dot/open-call/internal/models"
	"github.com/Sergio-dot/open-call/internal/render"
	"github.com/Sergio-dot/open-call/internal/repository"
	"github.com/Sergio-dot/open-call/internal/repository/dbrepo"
	"net/http"
)

// Repo is the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepo creates a new repository with access to AppConfig
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		return
	}
}

// SignIn is the handler to log the user in
func (m *Repository) SignIn(w http.ResponseWriter, r *http.Request) {
	// TODO: user login
}

// Room is the room page handler
func (m *Repository) Room(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "room.page.tmpl", &models.TemplateData{})
	if err != nil {
		return
	}
}
