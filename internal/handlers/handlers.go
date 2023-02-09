package handlers

import (
	"github.com/Sergio-dot/open-call/internal/auth"
	"golang.org/x/crypto/bcrypt"

	"log"

	"github.com/Sergio-dot/open-call/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	DB    *gorm.DB       // DB is an instance of gorm.DB
	Store *session.Store // Store is an instance of session.Store
)

// Home is the home page handler
func Home(ctx *fiber.Ctx) error {
	return ctx.Render("index", fiber.Map{
		"PageTitle": "OpenCall - Home",
	}, "layouts/main")
}

// Dashboard is the dashboard page handler
func Dashboard(ctx *fiber.Ctx) error {
	sess, err := Store.Get(ctx)
	if err != nil {
		log.Println(err)
		return ctx.Redirect("/")
	}

	return ctx.Render("dashboard", fiber.Map{
		"PageTitle": "OpenCall - Dashboard",
		"UserID":    sess.Get("userID"),
		"Email":     sess.Get("email"),
		"Username":  sess.Get("username"),
		"createdAt": sess.Get("createdAt"),
		"updatedAt": sess.Get("updatedAt"),
	}, "layouts/main")
}

// Login handles logging a user, validating his credentials before granting access
func Login(ctx *fiber.Ctx) error {
	// validate input
	email := ctx.FormValue("loginEmail")
	password := ctx.FormValue("loginPassword")

	// check password length
	if !auth.MinLength(password, 8) {
		log.Println("Password is too short")
		return ctx.Redirect("/")
	}

	// query database for credentials validation
	var user models.User
	err := DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		log.Println(err)
		return ctx.Redirect("/")
	}

	// compare the passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println(err)
		return ctx.Redirect("/")
	}

	// get context session
	sess, err := Store.Get(ctx)
	if err != nil {
		log.Println(err)
		return ctx.Redirect("/")
	}

	// store user ID in the session
	sess.Set("userID", user.ID)
	sess.Set("email", user.Email)
	sess.Set("username", user.Username)
	sess.Set("createdAt", user.CreatedAt)
	sess.Set("updatedAt", user.UpdatedAt)

	err = sess.Save()
	if err != nil {
		return err
	}
	log.Println("Session saved")

	return ctx.Redirect("/dashboard")
}

// SignUp handles user registration, if information input by user pass validation
func SignUp(ctx *fiber.Ctx) error {
	// validate input
	email := ctx.FormValue("signupEmail")
	username := ctx.FormValue("signupUsername")
	password := ctx.FormValue("signupPassword")
	confirmPassword := ctx.FormValue("signupConfirmPassword")

	// check password length
	if !auth.MinLength(password, 8) || !auth.MinLength(confirmPassword, 8) {
		log.Println("Password must be at least 8 characters long")
		return ctx.Redirect("/")
	}

	// check password matching
	if password != confirmPassword {
		log.Println("Passwords do not match")
		return ctx.Redirect("/")
	}

	// check if user already exists
	var count int
	DB.Model(&models.User{}).Where("email = ? or username = ?", email, username).Count(&count)
	if count > 0 {
		log.Println("User already exists")
		return ctx.Redirect("/")
	}

	// hash provided password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		log.Println(err)
		return ctx.Redirect("/")
	}

	// create new user
	user := &models.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
	}
	DB.Create(user)

	return ctx.Redirect("/")
}

// GoogleLogin handles logging the user in through his google account
func GoogleLogin(ctx *fiber.Ctx) error {
	path := auth.ConfigGoogle()
	url := path.AuthCodeURL("state")
	return ctx.Redirect(url)
}

// GoogleCallback handles google's response
func GoogleCallback(ctx *fiber.Ctx) error {
	token, err := auth.ConfigGoogle().Exchange(ctx.Context(), ctx.FormValue("code"))
	if err != nil {
		panic(err)
	}

	email := auth.GetEmail(token.AccessToken)
	return ctx.Status(200).JSON(fiber.Map{"email": email, "login": true})
}

// Logout handles the request to log the user out, destroying the session
func Logout(ctx *fiber.Ctx) error {
	sess, err := Store.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = sess.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	return ctx.Redirect("/")
}
