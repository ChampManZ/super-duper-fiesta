package main

import (
	"fmt"
	"log"
	"path/filepath"
	"server/config"
	"server/handlers"
	"server/helpers"
	"server/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "server/docs"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		rootDir := filepath.Join("..", ".env")
		godotenv.Load(rootDir)
		fmt.Println("Loading .env file from root directory")
	} else {
		log.Fatal("Error loading .env file")
	}
}

// @title Simple Social Feed API
// @version 1.0
// @description This is a simple social feed API for DataWow Take Home Assignment
// @host localhost:1323
// @BasePath /
func main() {
	// Connect to database
	config.ConnectDatabase()

	// Start server
	e := echo.New()
	e.Use(handlers.ServerHeader)

	// Register Validator for request binding
	e.Validator = helpers.NewValidator()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"}, // Since this application is for demo purposes only, we will allow only localhost:3000
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	routes.SetupRoutes(e)
	e.Logger.Fatal(e.Start(":1323"))
}
