package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pi-prakhar/go-gcp-auth/internal/models"
	"github.com/pi-prakhar/go-gcp-auth/internal/services"
	"github.com/pi-prakhar/go-gcp-auth/pkg/utils"
)

func HandleHome(w http.ResponseWriter, r *http.Request) {
	var html = `<html><body><a href="/api/v1/auth/google/login">Google Login</a></body></html>`
	fmt.Fprint(w, html)
}

func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauth2State, err := utils.GenerateRandomString(32)
	if err != nil {
		http.Error(w, "Failed to generate oauth state", http.StatusInternalServerError)
		return
	}
	services.SetOAuthStateCookie(&w, oauth2State)
	url := services.GetOAuth2Config().AuthCodeURL(oauth2State)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state string
	state := r.FormValue("state")
	oauth2State, err := services.GetOAuthStateFromCookie(r)
	if err != nil {
		http.Error(w, "Failed to get state from cookie : "+err.Error(), http.StatusInternalServerError)
	}

	if state != oauth2State {
		http.Error(w, "State is invalid", http.StatusBadRequest)
		return
	}

	// Exchange authorization code for access token
	oauth2Config := services.GetOAuth2Config()
	code := r.FormValue("code")
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Use the token to get user info
	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		http.Error(w, "Failed to get user info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Parse and display user info
	var userInfo models.GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		http.Error(w, "Failed to parse user info: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = services.SetJWTToken(w, userInfo.Email)
	if err != nil {
		http.Error(w, "Internal Server Error - Failed to set token", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User Info: %v\n", userInfo)
}

func HandleProtected(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	fmt.Fprintf(w, "Welcome %s! You are authenticated.", username)
}
