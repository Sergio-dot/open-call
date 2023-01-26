package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/Sergio-dot/open-call/internal/config"
)

var app *config.AppConfig

// NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}

// ClientError generates a client error
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error: ", status)
	http.Error(w, http.StatusText(status), status)
}

// ServerError generate a server error
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// IsAuthenticated checks if the user is authenticated
func IsAuthenticated(r *http.Request) bool {
	exists := app.Session.Exists(r.Context(), "user_id")
	return exists
}
