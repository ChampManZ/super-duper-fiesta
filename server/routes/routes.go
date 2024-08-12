package routes

import (
	"os"
	"server/handlers"
	"server/models"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupRoutes(e *echo.Echo) {
	// Versioning
	api := e.Group("/api/v1")

	// Public Routes for maintaining health check and swagger
	e.GET("/", handlers.HealthCheck)             // GET / (Health check endpoint)
	e.GET("/swagger/*", handlers.SwaggerHandler) // GET /swagger/* (Swagger documentation)

	// Public API Routes
	api.POST("/login", handlers.LoggedInUser) // POST /api/v1/login
	api.POST("/logout", handlers.Logout)      // POST /api/v1/logout
	api.POST("/users", handlers.CreateUser)   // POST /api/v1/users
	api.GET("/posts", handlers.GetPosts)      // GET /api/v1/posts

	//------------------------ Admin routes ------------------------//
	admin := api.Group("/admin")
	loggerConfig := middleware.LoggerConfig{
		Format: `${time_rfc3339} ${status} ${method} ${host}${path} ${latency_human}` + "\n",
	}
	admin.Use(middleware.LoggerWithConfig(loggerConfig))
	admin.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == os.Getenv("ADMIN_USERNAME") && password == os.Getenv("ADMIN_PASSWORD") {
			return true, nil
		}
		return false, nil
	}))
	admin.GET("/main", handlers.MainAdminPage) // GET /api/v1/admin/main (Main admin page)

	// Users
	admin.GET("/users", handlers.GetUsers)               // GET /api/v1/admin/users (Retrieve all users)
	admin.POST("/run-migrations", handlers.RunMigration) // POST /api/v1/admin/run-migrations (Run migrations)

	// Comments
	admin.GET("/comments", handlers.GetComments) // GET /api/v1/admin/comments

	// ComentUser
	admin.GET("/commentuser", handlers.GetCommentUser) // GET /api/v1/admin/commentuser (Retrieve comment-user associations)

	//------------------------ Cookie (For debug) ------------------------//
	cookie := api.Group("/cookie")
	cookie.Use(handlers.CookieChecker)
	cookie.GET("/main", handlers.MainAdminPage) // GET /api/v1/cookie/main

	//------------------------ JWT Protected Routes (Need authentication routes) ------------------------//
	jwt_protected := api.Group("/restricted")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(models.JWTClaims)
		},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}
	jwt_protected.Use(echojwt.WithConfig(config))
	jwt_protected.GET("/main", handlers.RestrictedHandler) // GET /api/v1/restricted/main

	// User routes
	jwt_protected.PUT("/users/:uid", handlers.UpdateUser) // PUT /api/v1/restricted/users/:uid (Update a user by ID)

	// Post routes
	jwt_protected.POST("/posts", handlers.CreatePost) // POST /api/v1/restricted/posts (Create a new post)

	// Comment routes
	jwt_protected.POST("/comments", handlers.PostComment) // POST /api/v1/restricted/comments (Create a new comment)
}
