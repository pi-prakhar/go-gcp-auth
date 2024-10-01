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

type GoogleAuthService struct{}

func NewGoogleAuthService() *GoogleAuthService {
	googleAuthService := &GoogleAuthService{}
	googleAuthService.initConfig()
	return googleAuthService
}

func (g *GoogleAuthService) initConfig() {
	//TODO : create constants file to store the scopes
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

func (g *GoogleAuthService) GetOAuth2Config() *oauth2.Config {
	return oauth2Config
}

func (g *GoogleAuthService) SetOAuthStateCookie(w *http.ResponseWriter, state string) {
	//TODO : create constants file to store the cookie name
	cookie := &http.Cookie{
		Name:     "oauthState",
		Value:    state,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(10 * time.Minute),
	}
	http.SetCookie(*w, cookie)
}

func (g *GoogleAuthService) GetOAuthStateFromCookie(r *http.Request) (string, error) {
	//TODO : create constants file to store the cookie name
	cookie, err := r.Cookie("oauthState")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func (g *GoogleAuthService) SetAuthCookie(w *http.ResponseWriter, token string) {
	//TODO : create constants file to store the cookie name
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

// TODO: getAuthCookie function
func (g *GoogleAuthService) generateAuthJWTToken(username string) (string, error) {

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

func (g *GoogleAuthService) SetJWTToken(w http.ResponseWriter, username string) error {
	// Generate the JWT token
	tokenString, err := g.generateAuthJWTToken(username)
	if err != nil {
		return err
	}

	// Set the JWT as a cookie
	//TODO : create constants file to store the cookie name
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
