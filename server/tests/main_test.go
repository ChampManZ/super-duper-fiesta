package tests

import (
	"fmt"
	"os"
	"server/config"
	"server/models"
	"testing"

	"log"
	"path/filepath"

	"github.com/joho/godotenv"
	gormSQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	rootDir := filepath.Join("../../.env.test.local")
	err := godotenv.Load(rootDir)

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT_MOCK")
	database := os.Getenv("DB_NAME_MOCK")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, database)

	db, err := gorm.Open(gormSQL.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	return db
}

func createTables() {
	config.DB.AutoMigrate(&models.User{})
}

func teardown() {
	migrator := config.DB.Migrator()
	migrator.DropTable(&models.CommentUser{})
	migrator.DropTable(&models.Post{})
	migrator.DropTable(&models.Comment{})
	migrator.DropTable(&models.User{})
}

func TestMain(m *testing.M) {
	config.DB = setupTestDB()
	createTables()
	code := m.Run()
	teardown()
	os.Exit(code)
}
