package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pi-prakhar/go-gcp-auth/internal/constants"
	"github.com/pi-prakhar/go-gcp-auth/internal/models"
	"github.com/pi-prakhar/go-gcp-auth/pkg/utils"
)

type AuthMiddleware struct{}

func (m *AuthMiddleware) IsAuthenticated(next func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(constants.GOOGLE_AUTH_TOKEN_COOKIE_NAME)
		var res utils.Responder
		if err != nil {

			if err == http.ErrNoCookie {
				res = &utils.ErrorResponse{
					Message:    "Unauthorized - No Token",
					StatusCode: http.StatusUnauthorized,
					Error:      err.Error(),
				}
				res.Write(w)
				return
			}
			res = &utils.ErrorResponse{
				Message:    "Bad Request",
				StatusCode: http.StatusBadRequest,
				Error:      err.Error(),
			}
			res.Write(w)
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
				res = &utils.ErrorResponse{
					Message:    "Unauthorized - Invalid Token",
					StatusCode: http.StatusUnauthorized,
					Error:      err.Error(),
				}
				res.Write(w)
				return
			}
			res = &utils.ErrorResponse{
				Message:    "Bad Request",
				StatusCode: http.StatusBadRequest,
				Error:      err.Error(),
			}
			res.Write(w)
			return
		}

		// Check if the token is valid
		if !token.Valid {
			res = &utils.ErrorResponse{
				Message:    "Unauthorized - Token Expired or Invalid",
				StatusCode: http.StatusUnauthorized,
			}
			res.Write(w)
			return
		}

		// If valid, set the username in the request context for further handlers
		r.Header.Set("username", claims.Username)

		// Call the next handler in the chain
		next(w, r)
	}
}
