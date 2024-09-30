package services

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pi-prakhar/go-gcp-auth/internal/models"
	"github.com/pi-prakhar/go-gcp-auth/pkg/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	oauth2Config *oauth2.Config
)

func init() {
	oauth2Config = &oauth2.Config{
		ClientID:     utils.GetClientId(),
		ClientSecret: utils.GetClientSecret(),
		RedirectURL:  utils.GetCallbackURL(),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email",
		},
		Endpoint: google.Endpoint,
	}
}

func GetOAuth2Config() *oauth2.Config {
	return oauth2Config
}

// Set cookie on the user's browser
func SetOAuthStateCookie(w *http.ResponseWriter, state string) {

	cookie := &http.Cookie{
		Name:     "oauthState",
		Value:    state,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Minute),
	}
	http.SetCookie(*w, cookie)
}

// Retrieve cookie from the request
func GetOAuthStateFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("oauthState")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func SetAuthCookie(w *http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "auth_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,                           // Prevents JavaScript access
		Secure:   true,                           // Ensures cookie is only sent over HTTPS
		Expires:  time.Now().Add(24 * time.Hour), // Set cookie expiration
	}
	http.SetCookie(*w, cookie)
}

func generateAuthJWTToken(username string) (string, error) {

	expirationTime := time.Now().Add(1 * time.Hour)

	claims := &models.Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Declare the token with the algorithm used for signing and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(utils.GetJWTKey())
	if err != nil {
		return "", err
	}

	return tokenString, nil

}

func SetJWTToken(w http.ResponseWriter, username string) error {
	// Generate the JWT token
	tokenString, err := generateAuthJWTToken(username)
	if err != nil {
		return err
	}

	// Set the JWT as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,                    // Ensure the cookie cannot be accessed by JavaScript
		Secure:   true,                    // Set to true if using HTTPS
		Path:     "/",                     // The path for which the cookie is valid
		SameSite: http.SameSiteStrictMode, // Ensure cookie is sent only for same-site requests
	})
	return nil
}
