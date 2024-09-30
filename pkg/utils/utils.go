package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

func GetClientId() string {
	clientId := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	fmt.Println(clientId)
	return clientId
}

func GetClientSecret() string {
	clientSecret := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	fmt.Println(clientSecret)
	return clientSecret
}

func GetCallbackURL() string {
	url := os.Getenv("AUTH_SERVICE_HOST") + "/api/v1/auth/google/callback"
	fmt.Println(os.Getenv("AUTH_SERVICE_HOST"))
	fmt.Println(url)
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
	fmt.Println(os.Getenv("AUTH_JWT_KEY"))
	return []byte(os.Getenv("AUTH_JWT_KEY"))
}
