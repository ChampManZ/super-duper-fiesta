package routes

import (
	"server/handlers"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	// Versioning
	api := e.Group("/api/v1")

	// Test Health (Public)
	e.GET("/", handlers.HealthCheck)

	// User
	users := api.Group("/users")
	users.GET("", handlers.GetUsers)        // GET /api/v1/users
	users.POST("", handlers.CreateUser)     // POST /api/v1/users
	users.PUT("/:uid", handlers.UpdateUser) // PUT /api/v1/users/:uid

	// Post
	posts := api.Group("/posts")
	posts.GET("", handlers.GetPosts) // GET /api/v1/posts

	// Comment
	comment := api.Group("/comments")
	comment.GET("/comments", handlers.GetComments) // GET /api/v1/comments

	// CommentUser
	commentUsers := api.Group("/commentuser")
	commentUsers.GET("/commentuser", handlers.GetCommentUser) // GET /api/v1/commentuser

	// Swagger
	e.GET("/swagger/*", handlers.SwaggerHandler)
}
