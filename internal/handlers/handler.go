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
	res := &utils.SuccessResponse[string]{
		Message:    "Successfully authenticated",
		StatusCode: http.StatusOK,
		Data:       username,
	}
	res.Write(w)

}

func (h *AuthHandler) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	oauth2State, err := utils.GenerateRandomString(32)
	var res utils.Responder
	if err != nil {
		res = &utils.ErrorResponse{
			Message:    "Failed to generate oauth2State",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		res.Write(w)
		return
	}
	h.services.SetOAuthStateCookie(&w, oauth2State)
	url := h.services.GetOAuth2Config().AuthCodeURL(oauth2State)

	res = &utils.SuccessResponse[string]{
		Message:    "Successfully generated Redirect URL for Google auth",
		StatusCode: http.StatusOK,
		Data:       url,
	}
	res.Write(w)
}

func (h *AuthHandler) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	// Verify state string
	state := r.FormValue("state")
	var res utils.Responder
	oauth2State, err := h.services.GetOAuthStateFromCookie(r)
	if err != nil {
		res = &utils.ErrorResponse{
			Message:    "Failed to get state from cookie",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		res.Write(w)
		return
	}

	if state != oauth2State {
		res = &utils.ErrorResponse{
			Message:    "State is invalid",
			StatusCode: http.StatusBadRequest,
		}
		res.Write(w)
		return
	}

	// Exchange authorization code for access token
	oauth2Config := h.services.GetOAuth2Config()
	code := r.FormValue("code")
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		res = &utils.ErrorResponse{
			Message:    "Failed to exchange token",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		res.Write(w)
		return
	}

	// Use the token to get user info
	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get(constants.GOOGLE_OAUTH_USER_INFO_ENDPOINT)
	if err != nil {
		res = &utils.ErrorResponse{
			Message:    "Failed to get user info",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		res.Write(w)
		return
	}
	defer resp.Body.Close()

	// Parse and display user info
	var userInfo models.GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		res = &utils.ErrorResponse{
			Message:    "Failed to parse user info",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		res.Write(w)
		return
	}

	// Creating JWT Token
	err = h.services.SetJWTToken(w, userInfo.Email)
	if err != nil {
		res = &utils.ErrorResponse{
			Message:    "Failed to create JWT token",
			StatusCode: http.StatusInternalServerError,
			Error:      err.Error(),
		}
		res.Write(w)
		return
	}

	res = &utils.SuccessResponse[models.GoogleUser]{
		Message:    "Successfully Fetched Google User Info",
		StatusCode: http.StatusOK,
		Data:       userInfo,
	}
	res.Write(w)
}
