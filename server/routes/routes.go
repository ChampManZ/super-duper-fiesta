package routes

import (
	"server/handlers"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	// Test Health
	e.GET("/", handlers.HealthCheck)

	// User
	e.GET("/users", handlers.GetUsers)

	// Swagger
	e.GET("/swagger/*", handlers.SwaggerHandler)
}
