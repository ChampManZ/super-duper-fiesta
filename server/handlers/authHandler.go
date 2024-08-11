package handlers

import (
	"log"
	"net/http"

	"server/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func AccessibleHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func RestrictedHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	return c.String(http.StatusOK, "Welcome "+username+"!")
}

// In case we want to check the cookie token, we can use this middleware
// Industry level will also hash the cookie token to be shorter as well
func CookieChecker(next echo.HandlerFunc) echo.HandlerFunc {
	var userID uint

	return func(c echo.Context) error {
		cookie, err := c.Cookie("sessionID")
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Session unauthorized"})
		}

		sessionToken := cookie.Value
		if res := config.DB.Table("users").Select("user_id").Where("session_token = ?", sessionToken).Scan(&userID); res.Error != nil {
			return c.String(http.StatusUnauthorized, "Session unauthorized or expired")
		}

		c.Set("userID", userID)

		return next(c)
	}
}

// For debug
func MainAdminPage(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the admin page")
}

func CookiePage(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the cookie page")
}

func JWTPage(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the JWT page")
}
