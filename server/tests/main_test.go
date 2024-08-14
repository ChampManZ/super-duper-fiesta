package tests

import (
	"fmt"
	"os"
	"server/config"
	"server/models"
	"testing"
	"time"

	"log"
	"path/filepath"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	gormSQL "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func createJWTTokenTest(t *testing.T, userID uint) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &models.JWTClaims{
		UserID: uint(userID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	})
	tokenString, err := token.SignedString([]byte("testing_mock"))
	if err != nil {
		t.Fatalf("Failed to create JWT token: %v", err)
	}
	return tokenString
}

func RollbackFunc(model interface{}) {
	config.DB.Begin()
	defer config.DB.Rollback()

	if model != nil {
		config.DB.Save(&model)
	}
}

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

	db, err := gorm.Open(gormSQL.Open(dsn), &gorm.Config{
		Logger: logger.New((log.New(os.Stdout, "\r\n", log.LstdFlags)),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:                  true,
			},
		),
	})
	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	return db
}

func createTables() {
	err := config.DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Failed to migrate User table: %v", err)
	}
	err = config.DB.AutoMigrate(&models.Post{})
	if err != nil {
		log.Fatalf("Failed to migrate Post table: %v", err)
	}
	err = config.DB.AutoMigrate(&models.Comment{})
	if err != nil {
		log.Fatalf("Failed to migrate Comment table: %v", err)
	}
	err = config.DB.AutoMigrate(&models.CommentUser{})
	if err != nil {
		log.Fatalf("Failed to migrate CommentUser table: %v", err)
	}

	fmt.Println("Tables migrated successfully")
}

func teardown() {
	migrator := config.DB.Migrator()
	migrator.DropTable(&models.CommentUser{})
	migrator.DropTable(&models.Comment{})
	migrator.DropTable(&models.Post{})
	migrator.DropTable(&models.User{})
}

func TestMain(m *testing.M) {
	config.DB = setupTestDB()
	createTables()
	code := m.Run()
	teardown()
	os.Exit(code)
}
