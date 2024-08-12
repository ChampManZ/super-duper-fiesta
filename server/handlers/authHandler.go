package handlers

import (
	"log"
	"net/http"

	"server/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// AccessibleHandler godoc
// @Summary Accessible route without authentication
// @Description This route is accessible to everyone without authentication
// @Tags accessible
// @Accept json
// @Produce json
// @Success 200 {string} string "Accessible"
func AccessibleHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

// RestrictedHandler godoc
// @Summary Restricted route with JWT authentication
// @Description This route is restricted and requires a valid JWT token to access
// @Tags restricted
// @Accept json
// @Produce json
// @Success 200 {string} string "Welcome [username]!"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /api/v1/restricted/main [get]
func RestrictedHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	username := claims["username"].(string)
	return c.String(http.StatusOK, "Welcome "+username+"!")
}

// CookieChecker godoc
// @Summary Middleware to check cookie session
// @Description This middleware checks if a valid session cookie is present
// @Tags middleware
// @Accept json
// @Produce json
// @Failure 401 {object} map[string]string "Session unauthorized"
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

// WriteLogInCookie godoc
// @Summary Write JWT token in cookie
// @Description This function writes a JWT token into a secure HttpOnly cookie
// @Tags cookie
// @Accept json
// @Produce json
// @Success 200 {string} string "Cookie set"
func WriteLogInCookie(c echo.Context, sessionToken string) {
	cookie := new(http.Cookie)
	cookie.Name = "JWTCookie"
	cookie.Value = sessionToken
	cookie.Expires = time.Now().Add(72 * time.Hour)
	cookie.HttpOnly = true
	cookie.Secure = true
	c.SetCookie(cookie)
}

// Logout godoc
// @Summary Logout user by clearing cookie
// @Description This route logs out the user by clearing the JWT cookie
// @Tags logout
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Logout successful"
// @Router /api/v1/logout [post]
func Logout(c echo.Context) error {
	// Clear cookie
	cookie := new(http.Cookie)
	cookie.Name = "JWTCookie"
	cookie.Value = ""
	cookie.Expires = time.Now()
	cookie.HttpOnly = true
	cookie.Secure = true
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, map[string]string{"message": "Logout successful"})
}

// For debug

// MainAdminPage godoc
// @Summary Admin main page
// @Description This is the main admin page accessible only to authenticated users
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {string} string "Welcome to the admin page"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Router /admin/main [get]
func MainAdminPage(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the admin page")
}

// CookiePage godoc
// @Summary Cookie debug page
// @Description This page is used for debugging cookies
// @Tags debug
// @Accept json
// @Produce json
// @Success 200 {string} string "Welcome to the cookie page"
// @Router /cookie-page [get]
func CookiePage(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the cookie page")
}

// JWTPage godoc
// @Summary JWT debug page
// @Description This page is used for debugging JWT tokens
// @Tags debug
// @Accept json
// @Produce json
// @Success 200 {string} string "Welcome to the JWT page"
// @Router /jwt-page [get]
func JWTPage(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to the JWT page")
}
