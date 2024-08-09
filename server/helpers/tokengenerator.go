package helpers

import (
	"os"
	"server/models"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWTToken(user models.User) (string, error) {
	claims := &jwt.StandardClaims{
		Id:        strconv.Itoa(int(user.UserID)),
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
