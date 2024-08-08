package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// HealthCheck godoc
// @Summary Health check
// @Description Check if the server is running
// @Tags healthcheck
// @Accept json
// @Produce json
// @Success 200 {string} string "Server is running"
// @Router / [get]
func HealthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "Server is running")
}
