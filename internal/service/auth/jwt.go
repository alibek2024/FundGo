package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func generateAccessToken(userID int) (string, error) {
	privateKeyPem, err := os.ReadFile("certs/private_key.pem")
	if err != nil {
		return "", fmt.Errorf("error opening file: %w", err)
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyPem)
	if err != nil {
		return "", fmt.Errorf("parsing error: %w", err)
	}

	claims := jwt.MapClaims{
		"iss": "auth-service",
		"sub": fmt.Sprintf("%d", userID),
		"exp": time.Now().Add(time.Minute * 15).Unix(),
		"iat": time.Now().Unix(),
		
	}
}
