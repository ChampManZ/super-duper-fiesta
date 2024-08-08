package handlers

import (
	"net/http"
	"server/config"
	"server/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/labstack/echo/v4"
)

// GetUsers godoc
// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Router /users [get]
func GetUsers(c echo.Context) error {
	var users []models.User
	config.DB.Find(&users)
	return c.JSON(http.StatusOK, users)
}

// CreateUser godoc
// @Summary Create single user
// @Description Create single user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User object that needs to be created"
// @Success 201 {object} models.User
// @Router /users [post]
func CreateUser(c echo.Context) error {
	request := new(models.CreateUserRequest)
	if err := c.Bind(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}
	if err := c.Validate(request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	// Hash Password
	hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	request.Password = string(hash)

	user := models.User{
		Username:  request.Username,
		Firstname: request.Firstname,
		Surname:   request.Surname,
		Email:     request.Email,
		Password:  string(hash),
	}

	if result := config.DB.Create(&user); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, user)
}

// References: https://pkg.go.dev/golang.org/x/crypto/bcrypt
// Overview: Package bcrypt implements Provos and Mazières's bcrypt adaptive hashing algorithm.
// See http://www.usenix.org/event/usenix99/provos/provos.pdf
