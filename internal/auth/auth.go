package auth

import (
	"encoding/json"
	"fmt"
	"github.com/Sergio-dot/open-call/internal/config"
	"github.com/Sergio-dot/open-call/internal/models"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"io"
	"log"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// ConfigGoogle to set oauth config for Google login
func ConfigGoogle() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     config.Config("GOOGLE_CLIENT_ID"),
		ClientSecret: config.Config("GOOGLE_SECRET_KEY"),
		RedirectURL:  config.Config("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	return conf
}

// ConfigFacebook to set oauth config for Facebook login
func ConfigFacebook() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     config.Config("FACEBOOK_APP_ID"),
		ClientSecret: config.Config("FACEBOOK_SECRET_KEY"),
		RedirectURL:  config.Config("FACEBOOK_REDIRECT_URL"),
		Scopes:       []string{"email"},
		Endpoint:     facebook.Endpoint,
	}

	return conf
}

// ConfigGithub to set oauth config for GitHub login
func ConfigGithub() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     config.Config("GITHUB_CLIENT_ID"),
		ClientSecret: config.Config("GITHUB_SECRET_KEY"),
		RedirectURL:  config.Config("GITHUB_REDIRECT_URL"),
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
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
