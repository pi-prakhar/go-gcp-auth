package utils

import (
	"crypto/rand"
	"encoding/base64"
	"os"
)

func GetClientId() string {
	clientId := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	return clientId
}

func GetClientSecret() string {
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	return clientSecret
}

func GetCallbackURL() string {
	url := os.Getenv("AUTH_SERVICE_HOST") + "/api/v1/auth/google/callback"
	return url
}

func GenerateRandomString(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func GetJWTKey() []byte {
	return []byte(os.Getenv("AUTH_JWT_KEY"))
}
