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

	// Public Routes for maintaining health check and swagger
	e.GET("/", handlers.HealthCheck)
	e.GET("/swagger/*", handlers.SwaggerHandler)

	// Public API Routes
	api.POST("/login", handlers.LoggedInUser) // POST /api/v1/login
	api.POST("/users", handlers.CreateUser)   // POST /api/v1/users

	// Public Posts is available to all users even if they are not logged in or registered
	// In the future, we can add a function to restrict access to certain posts
	// e.g. only friends can see the post (in case we have a friends system)
	api.GET("/posts", handlers.GetPosts) // GET /api/v1/posts

	// Same to comment, public comments are available to all users
	api.GET("/comments", handlers.GetComments) // GET /api/v1/comments

	// Protected routes (require JWT)
	protected := api.Group("")
	protected.Use(echojwt.JWT([]byte(os.Getenv("JWT_SECRET"))))

	// User routes
	protected.GET("/users", handlers.GetUsers)        // GET /api/v1/users
	protected.PUT("/users/:uid", handlers.UpdateUser) // PUT /api/v1/users/:uid

	// Post routes
	protected.POST("/posts", handlers.CreatePost) // POST /api/v1/posts

	// Comment routes

	// CommentUser routes
	protected.GET("/commentuser", handlers.GetCommentUser) // GET /api/v1/commentuser
}
