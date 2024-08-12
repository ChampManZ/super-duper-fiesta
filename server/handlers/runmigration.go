package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/labstack/echo/v4"
)

func RunMigration(c echo.Context) error {
	migrationDir := "./migrations"

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")
	DBURL := fmt.Sprintf("mysql://%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)

	m, err := migrate.New(
		"file://"+migrationDir,
		DBURL,
	)
	if err != nil {
		log.Println("Error creating migration instance:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error creating migration instance"})
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Println("Migration failed:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Migration failed"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Migration successful"})
}
