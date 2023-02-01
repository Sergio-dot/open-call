package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Sergio-dot/open-call/internal/config"
	"github.com/Sergio-dot/open-call/internal/driver"
	"github.com/Sergio-dot/open-call/internal/forms"
	"github.com/Sergio-dot/open-call/internal/models"
	"github.com/Sergio-dot/open-call/internal/render"
	"github.com/Sergio-dot/open-call/internal/repository"
	"github.com/Sergio-dot/open-call/internal/repository/dbrepo"
	"github.com/go-chi/chi"
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
	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)

		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// stores user id and username in the session
	m.App.Session.Put(r.Context(), "user_id", id)

	// prompt a flash message
	m.App.Session.Put(r.Context(), "toast", "Logged in")

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
	m.App.Session.Put(r.Context(), "toast", "Disconnected")
}

// UpdateUser handles update request by a user
func (m *Repository) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// convert id parameter to appropriate type
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	// retrieve user from database using the provided id
	u, err := m.DB.GetUserByID(id)
	if err != nil {
		log.Println("could not find user", err)
	}

	// get information to update
	uname := r.URL.Query().Get("uname")
	email := r.URL.Query().Get("email")

	u.Username = uname
	u.Email = email

	// update user data into database
	err = m.DB.UpdateUser(u)
	if err != nil {
		log.Fatal("could not update user", err)
	}

	// redirect user
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	m.App.Session.Put(r.Context(), "toast", "Updated successfully")
}

// Dashboard is the dashboard page handler
func (m *Repository) Dashboard(w http.ResponseWriter, r *http.Request) {
	// get user id from session
	id := m.App.Session.Get(r.Context(), "user_id")

	// pull user info from database
	user, err := m.DB.GetUserByID(id.(int))
	if err != nil {
		log.Println(err)
	}

	// store user info into template data
	data := make(map[string]interface{})
	data["id"] = user.ID
	data["username"] = user.Username
	data["email"] = user.Email
	data["createdAt"] = user.CreatedAt.Format("02-01-2006")

	err = render.Template(w, r, "dashboard.page.tmpl", &models.TemplateData{
		Data: data,
	})
	if err != nil {
		return
	}
}
