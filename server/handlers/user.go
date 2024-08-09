package handlers

import (
	"net/http"
	"server/config"
	"server/helpers"
	"server/models"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetUsers godoc
// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Param userid query int false "User ID"
// @Success 200 {object} []models.User
// @Success 200 {object} models.User
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Failed to retrieve users"
// @Router /users [get]
func GetUsers(c echo.Context) error {
	userID := c.QueryParam("uid")

	if userID != "" {
		userID, err := strconv.Atoi(userID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
		}

		var user models.User
		if result := config.DB.First(&user, userID); result.Error != nil {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
		}

		return c.JSON(http.StatusOK, user)
	}
	var users []models.User
	if result := config.DB.Find(&users); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get users"})
	}

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
	// Bind input data and validate request
	request := new(models.CreateUserRequest)
	if err := helpers.BindAndValidateRequest(c, request); err != nil {
		return err
	}

	// Hash Password
	hashedPassword, err := helpers.HashPassword(request.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to hash password"})
	}
	request.Password = hashedPassword

	user := models.User{
		Username:  request.Username,
		Firstname: request.Firstname,
		Surname:   request.Surname,
		Email:     request.Email,
		Password:  hashedPassword,
	}

	if result := config.DB.Create(&user); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, user)
}

// UpdateUser godoc
// @Summary Update single user
// @Description Update single user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Router /users/{id} [put]
func UpdateUser(c echo.Context) error {
	request := new(models.UpdateUserRequest)
	if err := helpers.BindAndValidateRequest(c, request); err != nil {
		return err
	}

	// Check if user already exists
	userIDParam := c.Param("uid")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	var user models.User
	if result := config.DB.First(&user, userID); result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	// Update user specific fields
	user.Username = request.Username
	user.Firstname = request.Firstname
	user.Surname = request.Surname

	// Handle Password change separately by checking if password is provided in the request
	if request.Password != "" {
		hashedPassword, err := helpers.HashPassword(request.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to hash password"})
		}
		user.Password = hashedPassword
	}

	// Save changes
	if result := config.DB.Save(&user); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update user"})
	}

	return c.JSON(http.StatusOK, user)
}

// References: https://pkg.go.dev/golang.org/x/crypto/bcrypt
// Overview: Package bcrypt implements Provos and Mazières's bcrypt adaptive hashing algorithm.
// See http://www.usenix.org/event/usenix99/provos/provos.pdf
