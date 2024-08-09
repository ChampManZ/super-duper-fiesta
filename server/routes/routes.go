package routes

import (
	"os"
	"server/handlers"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	// Versioning
	api := e.Group("/api/v1")

	// Public routes
	e.GET("/", handlers.HealthCheck)
	api.POST("/login", handlers.LoggedInUser) // POST /api/v1/login
	api.POST("/users", handlers.CreateUser)   // POST /api/v1/users
	e.GET("/swagger/*", handlers.SwaggerHandler)

	// Protected routes (require JWT)
	protected := api.Group("")
	protected.Use(echojwt.JWT([]byte(os.Getenv("JWT_SECRET"))))

	// User routes
	protected.GET("/users", handlers.GetUsers)        // GET /api/v1/users
	protected.PUT("/users/:uid", handlers.UpdateUser) // PUT /api/v1/users/:uid

	// Post routes
	protected.GET("/posts", handlers.GetPosts) // GET /api/v1/posts

	// Comment routes
	protected.GET("/comments", handlers.GetComments) // GET /api/v1/comments

	// CommentUser routes
	protected.GET("/commentuser", handlers.GetCommentUser) // GET /api/v1/commentuser
}
