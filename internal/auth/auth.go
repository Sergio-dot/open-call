package auth

import (
	"encoding/json"
	"fmt"
	"github.com/Sergio-dot/open-call/internal/config"
	"github.com/Sergio-dot/open-call/internal/models"
	"io"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// ConfigGoogle to set oauth config
func ConfigGoogle() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     config.Config("Client"),
		ClientSecret: config.Config("Secret"),
		RedirectURL:  config.Config("redirect_url"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	return conf
}

// GetEmail retrieve user's email
func GetEmail(token string) string {
	reqURL, err := url.Parse("https://www.googleapis.com/oauth2/v1/userinfo")

	if err != nil {
		log.Println(err)
	}

	ptoken := fmt.Sprintf("Bearer %s", token)
	res := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {ptoken}},
	}
	req, err := http.DefaultClient.Do(res)
	if err != nil {
		panic(err)

	}
	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var data models.GoogleResponse
	errors := json.Unmarshal(body, &data)
	if errors != nil {

		panic(errors)
	}
	return data.Email
}

// MinLength returns true if passed 'field' parameter is equal or longer than 'length', otherwise returns false
func MinLength(field string, length int) bool {
	if len(field) < length {
		return false
	}
	return true
}
