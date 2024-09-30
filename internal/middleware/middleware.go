package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pi-prakhar/go-gcp-auth/internal/models"
	"github.com/pi-prakhar/go-gcp-auth/pkg/utils"
)

// Middleware to check if the user is authenticated by verifying the JWT
func AuthMiddleware(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the JWT from the cookie
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Unauthorized - No Token", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Get the JWT string from the cookie
		tokenString := cookie.Value

		// Initialize a new instance of `Claims`
		claims := &models.Claims{}

		// Parse the JWT token, validating the signature
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return utils.GetJWTKey(), nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Unauthorized - Invalid Token", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		// Check if the token is valid
		if !token.Valid {
			http.Error(w, "Unauthorized - Token Expired or Invalid", http.StatusUnauthorized)
			return
		}

		// If valid, set the username in the request context for further handlers
		r.Header.Set("username", claims.Username)

		// Call the next handler in the chain
		next(w, r)
	}
}
