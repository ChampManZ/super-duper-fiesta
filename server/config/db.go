package config

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"server/helpers"
	"server/models"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"strings"

	"github.com/labstack/echo/v4"
	gormSQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Variable to store the environment variables
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_NAME")

	// DSN : data source name, used to open a database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)

	db, err := gorm.Open(gormSQL.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	DB = db
}

func RunMigration(c echo.Context) error {
	var req models.RunMigrationRequest
	if err := helpers.BindAndValidateRequest(c, &req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	migrationPath := "./migrations/" + req.MigrationID + ".sql"
	if _, err := os.Stat(migrationPath); os.IsNotExist(err) {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "Migration not found"})
	}

	sqlBytes, err := os.ReadFile(migrationPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error reading migration file"})
	}
	sqlQuery := string(sqlBytes)

	if res := DB.Exec(sqlQuery); res.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error running migration"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Migration ran successfully"})
}

func GetMigration(c echo.Context) error {
	entries, err := os.ReadDir("./migrations")
	if err != nil {
		log.Fatalf("Could not read migrations directory: %v", err)
	}

	var migrations []models.GetMigrationListRequest
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".sql" {
			migrationID := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			migrations = append(migrations, models.GetMigrationListRequest{
				MigrationID: migrationID,
				Title:       entry.Name(),
			})
		}
	}
	return c.JSON(http.StatusOK, migrations)
}
