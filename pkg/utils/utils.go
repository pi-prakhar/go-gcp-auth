package utils

import (
	"crypto/rand"
	"encoding/base64"
	"os"

	"log"
)

func GetClientId() string {
	clientIdFile := os.Getenv("GOOGLE_OAUTH_CLIENT_ID")
	clientId, err := os.ReadFile(clientIdFile)
	if err != nil {
		log.Println("Failed to load client id : ", err.Error())
	}
	return string(clientId)
}

func GetClientSecret() string {
	clientSecretFile := os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET")
	clientSecret, err := os.ReadFile(clientSecretFile)
	if err != nil {
		log.Println("Failed to load client secret : ", err.Error())
	}
	return string(clientSecret)
}

func GetCallbackURL() string {
	authServiceHostFile := os.Getenv("AUTH_SERVICE_HOST")
	authServiceHost, err := os.ReadFile(authServiceHostFile)
	if err != nil {
		log.Println("Failed to load auth service host : ", err.Error())
	}
	url := string(authServiceHost) + "/api/v1/auth/google/callback"
	log.Println("callback url:" + url)
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
	authJWTKeyFile := os.Getenv("AUTH_JWT_KEY")
	authJWTKey, err := os.ReadFile(authJWTKeyFile)
	if err != nil {
		log.Println("Failed to load auth jwt key : ", err.Error())
	}
	return authJWTKey
}
