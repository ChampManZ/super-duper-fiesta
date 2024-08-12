package handlers

import (
	"log"
	"net/http"
	"server/config"
	"server/helpers"
	"server/models"
	"strconv"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// GetUsers godoc
// @Summary Get all users or a specific user by ID
// @Description Retrieve all users or a specific user if User ID is provided in the query parameter
// @Tags Users
// @Accept json
// @Produce json
// @Param uid query int false "User ID"
// @Success 200 {object} []models.User "List of all users"
// @Success 200 {object} models.User "Details of a specific user"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Failed to retrieve users"
// @Router /api/v1/admin/users [get]
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
// @Summary Log in a user
// @Description Authenticate a user and return a JWT token
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.LoginUserRequest true "User login details"
// @Success 200 {object} map[string]string "Login successful, token returned"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 401 {object} map[string]string "Invalid username, email, or password"
// @Failure 500 {object} map[string]string "Failed to generate token"
// @Router /api/v1/login [post]
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

	token, err := helpers.GenerateJWTToken(user)
	if err != nil {
		log.Println("Error creating JWT token:", err)
		return c.String(http.StatusInternalServerError, "Failed to generate token")
	}

	WriteLogInCookie(c, token)

	if result := config.DB.Model(&user).Where("(username = ? OR email = ?) AND password = ?", request.Identifier, request.Identifier, user.Password).Update("cookie_token", token); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update user session token"})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Login successful",
		"token":   token,
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Register a new user with the provided details
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.CreateUserRequest true "User registration details"
// @Success 201 {object} models.User "Newly created user details"
// @Failure 409 {object} map[string]string "Username or email already exists"
// @Failure 500 {object} map[string]string "Failed to create user"
// @Router /api/v1/users [post]
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

	// Set isAdmin to 0 by default.
	// In a real-world application, this should be set to 0 by default and only set to 1 by an admin user.
	// isAdmin should be gain and lost by an admin user only, not through registration.
	isAdmin := "0"

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
		IsAdmin:   isAdmin,
	}

	if result := config.DB.Create(&user); result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Error: Failed to create user. Please try again"})
	}

	return c.JSON(http.StatusCreated, user)
}

// UpdateUser godoc
// @Summary Update an existing user
// @Description Update the details of an existing user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.UpdateUserRequest true "Updated user details"
// @Success 200 {object} models.User "Updated user details"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Failed to update user"
// @Router /api/v1/restricted/users/{uid} [put]
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
	if result := config.DB.Where("user_id = ?", userID).First(&user); result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	// Update user specific fields
	user.Username = request.Username
	user.Firstname = request.Firstname
	user.Surname = request.Surname

	updatedStruct := map[string]interface{}{
		"username":  user.Username,
		"firstname": user.Firstname,
		"surname":   user.Surname,
	}

	// Save changes
	if result := config.DB.Model(&user).Updates(updatedStruct); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update user"})
	}

	return c.JSON(http.StatusOK, user)
}

// ChangePassword godoc
// @Summary Change a user's password
// @Description Change the password of an existing user by ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body models.UpdateUserPasswordRequest true "New password"
// @Success 200 {object} map[string]string "Password updated successfully"
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 404 {object} map[string]string "User not found"
// @Failure 500 {object} map[string]string "Failed to update password"
// @Router /api/v1/restricted/users-update-password/{uid} [put]
func ChangePassword(c echo.Context) error {
	request := new(models.UpdateUserPasswordRequest)
	if err := helpers.BindAndValidateRequest(c, request); err != nil {
		return err
	}

	userIDParam := c.Param("uid")
	userID, err := strconv.Atoi(userIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
	}

	var user models.User
	if result := config.DB.First(&user, userID); result.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	// Handle Password change separately by checking if password is provided in the request
	if request.Password != "" {
		hashedPassword, err := helpers.HashPassword(request.Password)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to hash password"})
		}
		user.Password = hashedPassword
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Password is required"})
	}

	updatedStruct := map[string]interface{}{
		"password": user.Password,
	}

	if result := config.DB.Model(&user).Updates(updatedStruct); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to update password"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Password updated successfully"})
}

// References: https://pkg.go.dev/golang.org/x/crypto/bcrypt
// Overview: Package bcrypt implements Provos and Mazières's bcrypt adaptive hashing algorithm.
// See http://www.usenix.org/event/usenix99/provos/provos.pdf
