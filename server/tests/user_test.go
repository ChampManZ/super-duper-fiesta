package tests

import (
	"os"
	"strings"
	"testing"

	"net/http"
	"net/http/httptest"
	"net/mail"
	"reflect"
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
func TestCreateUser_Valid(t *testing.T) {
	createTables()
	defer teardown()

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
		assert.NotContains(t, rec.Body.String(), `"password"`)
		assert.NotContains(t, rec.Body.String(), `"email"`)
	}
}

func TestCreateUser_Invalid(t *testing.T) {
	e := echo.New()
	e.Validator = helpers.NewValidator()

	tests := []struct {
		name          string
		userJSON      string
		statusCode    int
		expectedError string
	}{
		{
			name:          "Short Username",
			userJSON:      `{"username":"te","firstname":"Test","surname":"User","email":"test@example.com","password":"password123"}`,
			statusCode:    http.StatusBadRequest,
			expectedError: `Error:Field validation for 'Username' failed on the 'min' tag`,
		},
		{
			name:          "Long Username",
			userJSON:      `{"username":"gasdhjgshfgasdhjfgasdfghjasdghfasdghjfgasdfghsfsfghaasdghjkfhjgsdafasdfsdff","firstname":"Test","surname":"User","email":"test@example.com","password":"password123"}`,
			statusCode:    http.StatusBadRequest,
			expectedError: `Error:Field validation for 'Username' failed on the 'max' tag`,
		},
		{
			name:          "Empty Username",
			userJSON:      `{"username":"","firstname":"Test","surname":"User","email":"test@example.com","password":"password123"}`,
			statusCode:    http.StatusBadRequest,
			expectedError: `Error:Field validation for 'Username' failed on the 'required' tag`,
		},
		{
			name:          "No Username Provided",
			userJSON:      `{"firstname":"Test","surname":"User","email":"test@example.com","password":"password123"}`,
			statusCode:    http.StatusBadRequest,
			expectedError: `Error:Field validation for 'Username' failed on the 'required' tag`,
		},
		{
			name:          "Wrong Type Username",
			userJSON:      `{"username":123,"firstname":"Test","surname":"User","email":"test@example.com","password":"password123"}`,
			statusCode:    http.StatusBadRequest,
			expectedError: `Invalid request data`,
		},
		{
			name:          "Invalid Email",
			userJSON:      `{"username":"testuser","firstname":"Test","surname":"User","email":"te","password":"password123"}`,
			statusCode:    http.StatusBadRequest,
			expectedError: `Error:Field validation for 'Email' failed on the 'email' tag`,
		},
		{
			name:          "Empty Email",
			userJSON:      `{"username":"testuser","firstname":"Test","surname":"User","email":"","password":"password123"}`,
			statusCode:    http.StatusBadRequest,
			expectedError: `Error:Field validation for 'Email' failed on the 'required' tag`,
		},
		{
			name:          "Invalid Password",
			userJSON:      `{"username":"testuser","firstname":"Test","surname":"User","email":test@example.com","password":"we"}`,
			statusCode:    http.StatusBadRequest,
			expectedError: `Invalid request data`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			createTables()
			defer teardown()

			req := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(test.userJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			actual := handlers.CreateUser(c)
			assert.Error(t, actual)
			httpError, ok := actual.(*echo.HTTPError)
			if ok {
				assert.Equal(t, test.statusCode, httpError.Code)

				switch message := httpError.Message.(type) {
				case string:
					assert.Contains(t, message, test.expectedError)
				case map[string]string:
					assert.Contains(t, message["message"], test.expectedError)
				default:
					t.Errorf("unexpected error message type: %T", message)
				}
			}
		})
	}
}

// Integration Testing
func TestLoginUserByUsername(t *testing.T) {
	createTables()
	defer teardown()

	e := echo.New()
	e.Validator = helpers.NewValidator()

	userJSON := `{"username":"testuser","firstname":"Test","surname":"User","email":"test@example.com","password":"password123"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(userJSON))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	createRec := httptest.NewRecorder()
	createC := e.NewContext(createReq, createRec)

	if assert.NoError(t, handlers.CreateUser(createC)) {
		assert.Equal(t, http.StatusCreated, createRec.Code)
	}

	loginJSON := `{"identifier":"testuser","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(loginJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.LoggedInUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestLoginUserByEmail(t *testing.T) {
	createTables()
	defer teardown()

	e := echo.New()
	e.Validator = helpers.NewValidator()

	userJSON := `{"username":"testuser","firstname":"Test","surname":"User","email":"test@example.com","password":"password123"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(userJSON))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	createRec := httptest.NewRecorder()
	createC := e.NewContext(createReq, createRec)

	if assert.NoError(t, handlers.CreateUser(createC)) {
		assert.Equal(t, http.StatusCreated, createRec.Code)
	}

	loginJSON := `{"identifier":"testuser","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(loginJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.LoggedInUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetUsers(t *testing.T) {
	createTables()
	defer teardown()

	e := echo.New()
	e.Validator = helpers.NewValidator()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
	req.Header.Set("auth", "username="+os.Getenv("ADMIN_USERNAME")+", password="+os.Getenv("ADMIN_PASSWORD")) // Set admin credentials
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.GetUsers(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

	}
}
