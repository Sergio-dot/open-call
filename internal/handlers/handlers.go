package handlers

import (
	"github.com/Sergio-dot/open-call/internal/auth"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
	"time"

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
	sess, err := Store.Get(ctx)
	if err != nil {
		log.Println(err)
		return ctx.Redirect("/")
	}

	// get messages from session
	successMessage, _ := sess.Get("success-message").(string)
	errorMessage, _ := sess.Get("error-message").(string)

	// remove message from the session
	sess.Delete("success-message")
	sess.Delete("error-message")
	sess.Save()

	return ctx.Render("index", fiber.Map{
		"PageTitle":    "OpenCall - Home",
		"ToastSuccess": successMessage,
		"ToastError":   errorMessage,
	}, "layouts/main")
}

// Dashboard is the dashboard page handler
func Dashboard(ctx *fiber.Ctx) error {
	sess, err := Store.Get(ctx)
	if err != nil {
		log.Println(err)
		return ctx.Redirect("/")
	}

	// get messages from session
	successMessage, _ := sess.Get("success-message").(string)
	errorMessage, _ := sess.Get("error-message").(string)

	// remove message from the session
	sess.Delete("success-message")
	sess.Delete("error-message")
	sess.Save()

	// Fix the loss of session data
	sess, err = Store.Get(ctx)
	if err != nil {
		log.Println(err)
		return ctx.Redirect("/")
	}

	return ctx.Render("dashboard", fiber.Map{
		"PageTitle":    "OpenCall - Dashboard",
		"SessionID":    sess.Get("sessionID"),
		"UserID":       sess.Get("userID"),
		"Email":        sess.Get("email"),
		"Username":     sess.Get("username"),
		"CreatedAt":    sess.Get("createdAt").(time.Time).Format("02-01-2006"),
		"UpdatedAt":    sess.Get("updatedAt").(time.Time).Format("02-01-2006"),
		"ToastSuccess": successMessage,
		"ToastError":   errorMessage,
	}, "layouts/main")
}

// Login handles logging a user, validating his credentials before granting access
func Login(ctx *fiber.Ctx) error {
	// get context session
	sess, err := Store.Get(ctx)
	if err != nil {
		log.Println(err)
		return ctx.Redirect("/")
	}

	// validate input
	email := ctx.FormValue("loginEmail")
	password := ctx.FormValue("loginPassword")

	// check password length
	if !auth.MinLength(password, 8) {
		sess.Set("error-message", "Password too short")
		sess.Save()
		return ctx.Redirect("/")
	}

	// query database for credentials validation
	var user models.User
	err = DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		sess.Set("error-message", "Wrong email or password")
		sess.Save()
		return ctx.Redirect("/")
	}

	// compare the passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		sess.Set("error-message", "Wrong email or password")
		sess.Save()
		return ctx.Redirect("/")
	}

	sess.Set("userID", user.ID)
	sess.Set("email", user.Email)
	sess.Set("username", user.Username)
	sess.Set("createdAt", user.CreatedAt)
	sess.Set("updatedAt", user.UpdatedAt)
	sess.Set("success-message", "Logged in")

	sess.Save()

	return ctx.Redirect("/dashboard")
}

// SignUp handles user registration, if information input by user pass validation
func SignUp(ctx *fiber.Ctx) error {
	// get context session
	sess, err := Store.Get(ctx)
	if err != nil {
		log.Println(err)
		return ctx.Redirect("/")
	}

	// validate input
	email := ctx.FormValue("signupEmail")
	username := ctx.FormValue("signupUsername")
	password := ctx.FormValue("signupPassword")
	confirmPassword := ctx.FormValue("signupConfirmPassword")

	// check password length
	if !auth.MinLength(password, 8) || !auth.MinLength(confirmPassword, 8) {
		sess.Set("error-message", "Password should be at least 8 characters long")
		sess.Save()
		return ctx.Redirect("/")
	}

	// check password matching
	if password != confirmPassword {
		sess.Set("error-message", "Passwords do not match")
		sess.Save()
		return ctx.Redirect("/")
	}

	// check if email already exists in database
	var count int
	DB.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		sess.Set("error-message", "Email is already taken")
		sess.Save()
		return ctx.Redirect("/")
	}
	// check if username already exists in database
	DB.Model(&models.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		sess.Set("error-message", "Username is already taken")
		sess.Save()
		return ctx.Redirect("/")
	}

	// hash provided password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		sess.Set("error-message", "Error while processing request. Try again")
		sess.Save()
		return ctx.Redirect("/")
	}

	// create new user
	user := &models.User{
		Email:    email,
		Username: username,
		Password: string(hashedPassword),
	}
	DB.Create(user)

	sess.Set("success-message", "Account created")
	sess.Save()

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
	// get context session
	sess, err := Store.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// exchanging the authorization code obtained from the Google OAuth2 flow for an access token
	token, err := auth.ConfigGoogle().Exchange(ctx.Context(), ctx.FormValue("code"))
	if err != nil {
		sess.Set("error-message", "Something went wrong. Try again")
		sess.Save()
		return ctx.Redirect("/")
	}

	// get email from token
	email := auth.GetEmail(token.AccessToken)
	username := strings.Split(email, "@")[0]

	// check if user is already in database
	var user models.User
	DB.Where("email = ?", email).First(&user)
	if user.ID == 0 {
		// create new user
		user = models.User{
			Email:    email,
			Username: username,
		}
		DB.Create(&user)
	}

	// log in the user
	sess.Set("userID", user.ID)
	sess.Set("email", user.Email)
	sess.Set("username", user.Username)
	sess.Set("createdAt", user.CreatedAt)
	sess.Set("updatedAt", user.UpdatedAt)
	sess.Set("success-message", "Logged in")
	sess.Save()

	return ctx.Redirect("/dashboard")
}

// UpdateUser handles the request to update user information such as username and email
func UpdateUser(ctx *fiber.Ctx) error {
	passwordChange := true

	// get context session
	sess, err := Store.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// get id from parameters and convert to int
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	// retrieve user by id from database
	var user models.User
	if err = DB.First(&user, id).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).SendString("User not found")
	}

	// get username to update
	if uname := ctx.Query("uname"); uname != "" {
		user.Username = uname
		passwordChange = false
	}

	// get email to update
	if email := ctx.Query("email"); email != "" {
		user.Email = email
		passwordChange = false
	}

	// update user data into database
	if !passwordChange {
		if err = DB.Save(&user).Error; err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString("Error updating user")
		}

		sess.Set("username", user.Username)
		sess.Set("email", user.Email)
		sess.Set("success-message", "Updated successfully")
		sess.Save()

		return ctx.Redirect("/dashboard")
	}

	if passwordChange {
		// get password fields
		op := ctx.Query("op") // old password
		np := ctx.Query("np") // new password
		cp := ctx.Query("cp") // confirm new password

		// check for empty fields
		if op == "" || np == "" || cp == "" {
			sess.Set("error-message", "Fields can't be empty")
			sess.Save()
			return ctx.Redirect("/dashboard")
		}

		// check new password length
		if len(np) < 8 {
			sess.Set("error-message", "Password should be at least 8 characters long")
			sess.Save()
			return ctx.Redirect("/dashboard")
		}

		// check if new password and confirm password are the same
		if cp != np {
			sess.Set("error-message", "'Confirm password' must match 'New password'")
			sess.Save()
			return ctx.Redirect("/dashboard")
		}

		// check if old password is the same of the stored password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(op))
		if err != nil {
			log.Println("Password do not match")
			sess.Set("error-message", "Password doesn't match")
			sess.Save()
			return ctx.Redirect("/dashboard")
		}

		// hash new password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(np), 12)
		if err != nil {
			return err
		}

		// update password field of user model
		user.Password = string(hashedPassword)

		if err = DB.Save(&user).Error; err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString("Error updating user")
		}

		sess.Set("success-message", "Changed password successfully")
		sess.Save()
		return ctx.Redirect("/dashboard")
	}

	return nil
}

// Logout handles the request to log the user out, destroying the session
func Logout(ctx *fiber.Ctx) error {
	// get context session
	sess, err := Store.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// destroy session
	err = sess.Destroy()
	if err != nil {
		log.Fatal(err)
	}

	// message to show a toast
	sess.Set("success-message", "Disconnected")
	sess.Save()
	return ctx.Redirect("/")
}
