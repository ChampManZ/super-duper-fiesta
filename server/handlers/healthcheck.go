package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck godoc
// @Summary Perform a health check on the server
// @Description Check if the server is running and responsive
// @Tags HealthCheck
// @Accept json
// @Produce json
// @Success 200 {string} string "Server is running"
// @Router / [get]
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Server is running")
}
