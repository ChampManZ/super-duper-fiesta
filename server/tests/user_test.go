package tests

import (
	"os"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"
	"net/mail"
	"reflect"
	"server/config"
	"server/handlers"
	"server/helpers"
	"server/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// ----------- Model Testing ----------- //
func TestUserModel(t *testing.T) {
	cases := []struct {
		name          string
		user          models.User
		expectedError bool
	}{
		{
			name: "Valid User",
			user: models.User{
				Username:    "testuser",
				Firstname:   "Test",
				Surname:     "User",
				Email:       "test@test.com",
				Password:    "password123",
				IsAdmin:     "0",
				CookieToken: "cookie",
			},
			expectedError: false,
		},
		{
			name: "Invalid Username (Too short)",
			user: models.User{
				Username:    "te",
				Firstname:   "Test",
				Surname:     "User",
				Email:       "test@test.com",
				Password:    "password123",
				IsAdmin:     "0",
				CookieToken: "cookie",
			},
			expectedError: true,
		},
		{
			name: "Invalid Username (Too long)",
			user: models.User{
				Username:    "gasdhjgshfgasdhjfgasdfghjasdghfasdghjfgasdfghsfsfghaasdghjkfhjgsdafasdfsdff",
				Firstname:   "Test",
				Surname:     "User",
				Email:       "test@test.com",
				Password:    "password123",
				IsAdmin:     "0",
				CookieToken: "cookie",
			},
			expectedError: true,
		},
		{
			name: "Invalid Username (Int)",
			user: models.User{
				Firstname:   "Test",
				Surname:     "User",
				Email:       "test@test.com",
				Password:    "password123",
				IsAdmin:     "0",
				CookieToken: "cookie",
			},
			expectedError: true,
		},
		{
			name: "Invalid Email",
			user: models.User{
				Username:    "testuser",
				Firstname:   "Test",
				Surname:     "User",
				Email:       "test",
				Password:    "password123",
				IsAdmin:     "0",
				CookieToken: "cookie",
			},
			expectedError: true,
		},
		{
			name: "Invalid Password",
			user: models.User{
				Username:    "testuser",
				Firstname:   "Test",
				Surname:     "User",
				Email:       "test@test.com",
				Password:    "pass",
				IsAdmin:     "0",
				CookieToken: "cookie",
			},
			expectedError: true,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			err := ValidateUser(testCase.user)
			if (err != nil) != testCase.expectedError {
				t.Errorf("Expected error is %v but got %v", testCase.expectedError, err)
			}
		})
	}
}

// Requirements:
// A username must be minimum 3 to 32 characters maxmimum
// Firstname is require and must be a string
// Surname is require and must be a string
// Email is require and must be a valid email
// Password is require and must be a string with minimum 8 characters
func ValidateUser(u models.User) error {
	if reflect.TypeOf(u.Username) != reflect.TypeOf("") || len(u.Username) < 3 || len(u.Username) > 32 {
		return assert.AnError
	}
	if reflect.TypeOf(u.Firstname) != reflect.TypeOf("") {
		return assert.AnError
	}
	if reflect.TypeOf(u.Surname) != reflect.TypeOf("") {
		return assert.AnError
	}
	if !ValidEmail(u.Email) {
		return assert.AnError
	}
	if reflect.TypeOf(u.Password) != reflect.TypeOf("") || len(u.Password) < 8 {
		return assert.AnError
	}
	return nil
}

func ValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// ----------- API Testing ----------- //
func TestCreateUser(t *testing.T) {
	e := echo.New()
	e.Validator = helpers.NewValidator()

	userJSON := `{"username":"testuser","firstname":"Test","surname":"User","email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(userJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON) // Set the request content type to JSON
	rec := httptest.NewRecorder()                                    // Create a response recorder
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.CreateUser(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Contains(t, rec.Body.String(), `"username":"testuser"`)
		assert.Contains(t, rec.Body.String(), `"firstname":"Test"`)
		assert.Contains(t, rec.Body.String(), `"surname":"User"`)
	}
}

func TestLoginUserByUsername(t *testing.T) {
	e := echo.New()
	e.Validator = helpers.NewValidator()

	mockUser := models.User{
		Username:  "testuser",
		Firstname: "Test",
		Surname:   "User",
		Email:     "test@example.com",
		Password:  "password123",
	}

	config.DB.Create(&mockUser)

	loginJSON := `{"identifier":"testuser","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(loginJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.LoggedInUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"message": "Login successful"`)
		assert.Contains(t, rec.Body.String(), `"token"`)
	}
}

func TestLoginUserByEmail(t *testing.T) {
	e := echo.New()
	e.Validator = helpers.NewValidator()

	mockUser := models.User{
		Username:  "testuser",
		Firstname: "Test",
		Surname:   "User",
		Email:     "test@example.com",
		Password:  "password123",
	}

	config.DB.Create(&mockUser)

	loginJSON := `{"identifier":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(loginJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.LoggedInUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"message": "Login successful"`)
		assert.Contains(t, rec.Body.String(), `"token"`)
	}
}

func TestGetUsers(t *testing.T) {
	e := echo.New()
	e.Validator = helpers.NewValidator()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
	req.Header.Set("auth", "username="+os.Getenv("ADMIN_USERNAME")+", password="+os.Getenv("ADMIN_PASSWORD")) // Set admin credentials
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.GetUsers(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"id"`)
		assert.Contains(t, rec.Body.String(), `"username"`)
		assert.Contains(t, rec.Body.String(), `"firstname"`)
		assert.Contains(t, rec.Body.String(), `"surname"`)
		assert.Contains(t, rec.Body.String(), `"email"`)
		assert.Contains(t, rec.Body.String(), `"is_admin"`)
		assert.Contains(t, rec.Body.String(), `"Posts"`)
		assert.Contains(t, rec.Body.String(), `"Comments"`)
	}
}
