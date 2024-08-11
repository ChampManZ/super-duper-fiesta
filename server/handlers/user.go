package handlers

import (
	"log"
	"net/http"
	"os"
	"server/config"
	"server/helpers"
	"server/models"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
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

// LoggedInUser godoc
// @Summary Login user
// @Description Login user
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.LoginUserRequest true "User object that needs to be logged in"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Invalid username or email"
// @Failure 401 {object} map[string]string "Invalid password"
// @Failure 500 {object} map[string]string "Failed to generate token"
// @Router /login [post]
func LoggedInUser(c echo.Context) error {
	request := new(models.LoginUserRequest)
	if err := helpers.BindAndValidateRequest(c, request); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	var user models.User
	if result := config.DB.Where("username = ? OR email = ?", request.Identifier, request.Identifier).First(&user); result.Error != nil {
		return c.JSON((http.StatusUnauthorized), map[string]string{"message": "Invalid username or email"})
	}
	// Reference: CompareHashAndPassword compares a bcrypt hashed password with its possible plaintext equivalent. Returns nil on success, or an error on failure.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Invalid password"})
	}

	userSessionToken := uuid.New().String()

	if result := config.DB.Model(&user).Where("username = ? OR email = ?", request.Identifier, request.Identifier).Update("cookie_token", userSessionToken); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update user session token"})
	}

	// In this project, we are not focusing on the security of the cookie token much
	// We will focus on main function instead, but we still have this to set priority on security
	cookie := &http.Cookie{}
	cookie.Name = "sessionID"
	cookie.Value = userSessionToken
	cookie.Expires = time.Now().Add(72 * time.Hour)
	cookie.HttpOnly = true
	cookie.Secure = true

	c.SetCookie(cookie)

	token, err := helpers.GenerateJWTToken(user)
	if err != nil {
		log.Println("Error creating JWT token:", err)
		return c.String(http.StatusInternalServerError, "Failed to generate token")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Login successful",
		"token":   token,
	})
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

	// Check if user info already exists
	var existingUser models.User
	if result := config.DB.Where("username = ? OR email = ?", request.Username, request.Email).First(&existingUser); result.Error == nil {
		if existingUser.Username == request.Username {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Username already exists"})
		} else if existingUser.Email == request.Email {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Email already exists"})
		} else {
			return c.JSON(http.StatusConflict, map[string]string{"message": "Failed to create user. Please try again"})
		}
	}

	// Hash Password
	hashedPassword, err := helpers.HashPassword(request.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to hash password"})
	}
	request.Password = hashedPassword

	var isAdmin string
	if request.Username == os.Getenv("ADMIN_USERNAME") && request.Password == os.Getenv("ADMIN_PASSWORD") {
		isAdmin = "1"
	}

	user := models.User{
		Username:  request.Username,
		Firstname: request.Firstname,
		Surname:   request.Surname,
		Email:     request.Email,
		Password:  hashedPassword,
		IsAdmin:   isAdmin,
	}

	if result := config.DB.Create(&user); result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Error: Failed to create user. Please try again"})
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
