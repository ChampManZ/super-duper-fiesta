package main

import (
	"log"
	"server/config"
	"server/routes"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func init() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	// Connect to database
	config.ConnectDatabase()

	// Start server
	e := echo.New()
	routes.SetupRoutes(e)
	e.Logger.Fatal(e.Start(":1323"))
}
