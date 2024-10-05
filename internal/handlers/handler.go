package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pi-prakhar/go-gcp-auth/internal/constants"
	"github.com/pi-prakhar/go-gcp-auth/internal/models"
	"github.com/pi-prakhar/go-gcp-auth/internal/services"
	"github.com/pi-prakhar/go-gcp-auth/pkg/utils"
)

type AuthHandler struct {
	services *services.GoogleAuthService
}

func NewAuthHandler(services *services.GoogleAuthService) *AuthHandler {
	handler := &AuthHandler{
		services: services,
	}

	return handler
}

func (h *AuthHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	var html = `<html><body><a href="/api/v1/auth/google/login">Google Login</a></body></html>`
	fmt.Fprint(w, html)
}

func (h *AuthHandler) HandleProtected(w http.ResponseWriter, r *http.Request) {
	username := r.Header.Get("username")
	fmt.Fprintf(w, "Welcome %s! You are authenticated.", username)
}

func (h *AuthHandler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauth2State, err := utils.GenerateRandomString(32)
	if err != nil {
		http.Error(w, "Failed to generate oauth state", http.StatusInternalServerError)
		return
	}
	h.services.SetOAuthStateCookie(&w, oauth2State)
	url := h.services.GetOAuth2Config().AuthCodeURL(oauth2State)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (h *AuthHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state string
	state := r.FormValue("state")
	oauth2State, err := h.services.GetOAuthStateFromCookie(r)
	if err != nil {
		http.Error(w, "Failed to get state from cookie : "+err.Error(), http.StatusInternalServerError)
	}

	if state != oauth2State {
		http.Error(w, "State is invalid", http.StatusBadRequest)
		return
	}

	// Exchange authorization code for access token
	oauth2Config := h.services.GetOAuth2Config()
	code := r.FormValue("code")
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Use the token to get user info
	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get(constants.GOOGLE_OAUTH_USER_INFO_ENDPOINT)
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

	err = h.services.SetJWTToken(w, userInfo.Email)
	if err != nil {
		http.Error(w, "Internal Server Error - Failed to set token", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "User Info: %v\n", userInfo)
}
