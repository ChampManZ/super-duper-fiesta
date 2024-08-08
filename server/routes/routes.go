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

	// Post
	e.GET("/posts", handlers.GetPosts)

	// Comment
	e.GET("/comments", handlers.GetComments)

	// CommentUser
	e.GET("/commentuser", handlers.GetCommentUser)

	// Swagger
	e.GET("/swagger/*", handlers.SwaggerHandler)
}
