package helpers

import (
	"os"
	"server/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWTToken(user models.User) (string, error) {
	// Best standard is to have a standard claim as another object
	// Reference: https://pkg.go.dev/github.com/golang-jwt/jwt/v5#NewWithClaims
	claims := models.JWTClaims{
		UserID:   user.UserID,
		Username: user.Username,
		Admin:    user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
