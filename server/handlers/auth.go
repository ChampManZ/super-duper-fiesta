package handlers

import (
	"github.com/labstack/echo/v4"
)

func JWTAPIMiddleware(secret string) echo.MiddlewareFunc {
	return nil
}

// References: https://echo.labstack.com/docs/middleware/key-auth

// Key auth middleware provides a key based authentication.

// For valid key it calls the next handler.
// For invalid key, it sends "401 - Unauthorized" response.
// For missing key, it sends "400 - Bad Request" response.
