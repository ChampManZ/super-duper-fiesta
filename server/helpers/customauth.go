package helpers

import (
	"net/http"
	"os"
	"sync"
	"time"

	"crypto/subtle"

	"github.com/labstack/echo/v4"
)

type Session struct {
	Expiry time.Time
}

var sessions = struct {
	sync.RWMutex // Read Write mutex, guards the session store
	store        map[string]Session
}{
	store: make(map[string]Session),
}

func CustomBasicAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		username, password, ok := c.Request().BasicAuth()
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		// Check if username and password are correct
		if subtle.ConstantTimeCompare([]byte(username), []byte(os.Getenv("ADMIN_USERNAME"))) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte(os.Getenv("ADMIN_PASSWORD"))) == 1 {

			// Check if session exists and is not expired
			sessions.RLock()
			session, exists := sessions.store[username]
			sessions.RUnlock()

			if exists && time.Now().Before(session.Expiry) {
				// Session is valid, allow access
				return next(c)
			}

			// Create a new session with an expiry time in 24 hours
			sessions.Lock()
			sessions.store[username] = Session{Expiry: time.Now().Add(24 * time.Hour)}
			sessions.Unlock()

			// Continue to the next handler
			return next(c)
		}

		return echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
	}
}
