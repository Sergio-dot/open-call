package handlers

import (
	"github.com/Sergio-dot/open-call/internal/config"
	"github.com/Sergio-dot/open-call/internal/driver"
	"github.com/Sergio-dot/open-call/internal/forms"
	"github.com/Sergio-dot/open-call/internal/models"
	"github.com/Sergio-dot/open-call/internal/render"
	"github.com/Sergio-dot/open-call/internal/repository"
	"github.com/Sergio-dot/open-call/internal/repository/dbrepo"
	"log"
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
	err := render.Template(w, r, "home.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
	if err != nil {
		return
	}
}

// SignIn handles signing the user in
func (m *Repository) SignIn(w http.ResponseWriter, r *http.Request) {
	// refresh the session token
	_ = m.App.Session.RenewToken(r.Context())

	// parse the login form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	// retrieve login credentials input by user
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// form validation
	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// authenticate user
	id, username, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)

		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// stores user id and username in the session
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "username", username)

	// prompt a flash message
	m.App.Session.Put(r.Context(), "flash", "Logged in")

	// redirect user to dashboard
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// SignUp handles the user registration
func (m *Repository) SignUp(w http.ResponseWriter, r *http.Request) {
	// refresh the session token
	_ = m.App.Session.RenewToken(r.Context())

	// parse the sign-up form
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	// retrieve registration credentials input by user - password is hashed directly, for security reasons
	username := r.Form.Get("registrationUsername")
	email := r.Form.Get("registrationEmail")
	password := r.Form.Get("registrationPassword")
	repeatPassword := r.Form.Get("registrationRepeatPassword")

	// check if 'password' matches 'repeatPassword'
	if password != repeatPassword {
		m.App.Session.Put(r.Context(), "warning", "Password doesn't match")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	newUser := models.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	// form validation
	form := forms.New(r.PostForm)
	form.Required("registrationUsername", "registrationEmail", "registrationPassword", "registrationRepeatPassword")
	form.IsEmail("email")

	// insert user into the database
	err = m.DB.CreateUser(newUser)
	if err != nil {
		return
	}

	// prompt a flash message
	m.App.Session.Put(r.Context(), "flash", "Registered successfully. Now you can sign in with your credentials")

	// redirect user to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// SignOut handles signing the user out
func (m *Repository) SignOut(w http.ResponseWriter, r *http.Request) {
	// destroy the session
	_ = m.App.Session.Destroy(r.Context())

	// refresh the session token
	_ = m.App.Session.RenewToken(r.Context())

	// redirect user to home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Dashboard is the dashboard page handler
func (m *Repository) Dashboard(w http.ResponseWriter, r *http.Request) {
	// get username from session
	username := m.App.Session.Get(r.Context(), "username")

	data := make(map[string]interface{})
	data["username"] = username

	err := render.Template(w, r, "dashboard.page.tmpl", &models.TemplateData{
		Data: data,
	})
	if err != nil {
		return
	}
}

// Room is the room page handler
func (m *Repository) Room(w http.ResponseWriter, r *http.Request) {
	err := render.Template(w, r, "room.page.tmpl", &models.TemplateData{})
	if err != nil {
		return
	}
}
