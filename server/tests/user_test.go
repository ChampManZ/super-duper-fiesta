package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/mail"
	"os"
	"reflect"
	"server/handlers"
	"server/helpers"
	"server/models"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func GenerateNewUser(t *testing.T) {
	e := echo.New()
	e.Validator = helpers.NewValidator()

	userJSON := `{"username":"testuser","firstname":"Test","surname":"User","email":"test@example.com","password":"password123"}`
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(userJSON))
	createReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	createRec := httptest.NewRecorder()
	createCtx := e.NewContext(createReq, createRec)

	if assert.NoError(t, handlers.CreateUser(createCtx)) {
		assert.Equal(t, http.StatusCreated, createRec.Code)
	}

	var createdUser models.User
	err := json.Unmarshal(createRec.Body.Bytes(), &createdUser)
	if err != nil {
		fmt.Println("Failed to unmarshal created user")
	} else {
		fmt.Printf("Created User ID: %d\n", createdUser.UserID) // Created User ID: 1 | Gorm auto increments the ID starting from 1
	}
}

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

	GenerateNewUser(t)

	loginJSON := `{"identifier":"testuser","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(loginJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.LoggedInUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "token")
	}
}

func TestLoginUserByEmail(t *testing.T) {
	createTables()
	defer teardown()

	e := echo.New()
	e.Validator = helpers.NewValidator()

	GenerateNewUser(t)

	loginJSON := `{"identifier":"testuser","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/login", strings.NewReader(loginJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.LoggedInUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestLogout(t *testing.T) {
	createTables()
	defer teardown()

	e := echo.New()
	e.Validator = helpers.NewValidator()

	GenerateNewUser(t)

	logoutReq := httptest.NewRequest(http.MethodPost, "/api/v1/logout", nil)
	logoutReq.Header.Set("auth", "username=testuser, password=password123")
	logoutRec := httptest.NewRecorder()
	logoutC := e.NewContext(logoutReq, logoutRec)

	if assert.NoError(t, handlers.Logout(logoutC)) {
		assert.Equal(t, http.StatusOK, logoutRec.Code)
		assert.Contains(t, logoutRec.Body.String(), "Logout successful")
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

func TestGetUsersByID(t *testing.T) {
	mockUser := models.User{
		Username:  "testuser",
		Firstname: "Test",
		Surname:   "User",
		Email:     "test@example.com",
		IsAdmin:   "0",
	}

	createTables()
	defer teardown()

	e := echo.New()
	e.Validator = helpers.NewValidator()

	GenerateNewUser(t)

	uid := "1"

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users/:uid", nil)
	getReq.Header.Set("auth", "username="+os.Getenv("ADMIN_USERNAME")+", password="+os.Getenv("ADMIN_PASSWORD")) // Set admin credentials
	getRec := httptest.NewRecorder()
	getCtx := e.NewContext(getReq, getRec)
	getCtx.SetParamNames("uid")
	getCtx.SetParamValues(uid)

	if assert.NoError(t, handlers.GetUsers(getCtx)) {
		assert.Equal(t, http.StatusOK, getRec.Code)

		var responseBody interface{}
		err := json.Unmarshal(getRec.Body.Bytes(), &responseBody)
		assert.NoError(t, err)

		switch val := responseBody.(type) {
		case []interface{}:
			for _, element := range val {
				if userMap, ok := element.(map[string]interface{}); ok {
					assert.Equal(t, mockUser.Username, userMap["username"])
					assert.Equal(t, mockUser.Firstname, userMap["firstname"])
					assert.Equal(t, mockUser.Surname, userMap["surname"])
					assert.Equal(t, mockUser.Email, userMap["Email"])
					assert.Equal(t, mockUser.IsAdmin, userMap["is_admin"])
				} else {
					panic("Failed test")
				}
			}
		case map[string]interface{}:
			t.Errorf("Unexpected JSON structure: expected object, got array")
		default:
			t.Errorf("Unknown JSON structure")
		}
	}
}

func TestUpdateUser(t *testing.T) {

	e := echo.New()
	e.Validator = helpers.NewValidator()

	mockUpdateUserRequest := models.UpdateUserRequest{
		Username:  "newuser",
		Firstname: "New",
		Surname:   "User",
	}

	testCase := []struct {
		name           string
		userID         string
		body           models.UpdateUserRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Invalid User ID",
			userID:         "abc",
			body:           mockUpdateUserRequest,
			expectedStatus: http.StatusBadRequest,
			expectedError:  "Invalid input",
		},
		{
			name:           "Valid User",
			userID:         "1",
			body:           mockUpdateUserRequest,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "User not found",
			userID:         "2",
			body:           mockUpdateUserRequest,
			expectedStatus: http.StatusNotFound,
			expectedError:  `User not found`,
		},
		{
			name:   "Empty Username",
			userID: "1",
			body: models.UpdateUserRequest{
				Username:  "",
				Firstname: "New",
				Surname:   "User",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  `Error:Field validation for 'Username' failed on the 'required' tag`,
		},
		{
			name:   "Empty Firstname",
			userID: "1",
			body: models.UpdateUserRequest{
				Username:  "NewUser",
				Firstname: "",
				Surname:   "User",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  `Error:Field validation for 'Firstname' failed on the 'required' tag`,
		},
		{
			name:   "Empty Surname",
			userID: "1",
			body: models.UpdateUserRequest{
				Username:  "NewUser",
				Firstname: "New",
				Surname:   "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  `Error:Field validation for 'Surname' failed on the 'required' tag`,
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {

			createTables()
			defer teardown()

			GenerateNewUser(t)

			reqBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/restricted/users/"+tt.userID, strings.NewReader(string(reqBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/restricted/users/:uid")
			c.SetParamNames("uid")
			c.SetParamValues(tt.userID)

			err := handlers.UpdateUser(c)

			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedError != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedError)
			}
		})
	}
}

func TestUpdateUserPassword(t *testing.T) {
	e := echo.New()
	e.Validator = helpers.NewValidator()

	mockUpdatePasswordRequest := models.UpdateUserPasswordRequest{
		Password: "newpassword123",
	}

	testCase := []struct {
		name           string
		userID         string
		body           models.UpdateUserPasswordRequest
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Invalid User ID",
			userID:         "abc",
			body:           mockUpdatePasswordRequest,
			expectedStatus: http.StatusBadRequest,
			expectedError:  `Invalid input`,
		},
		{
			name:           "User not found",
			userID:         "2",
			body:           mockUpdatePasswordRequest,
			expectedStatus: http.StatusNotFound,
			expectedError:  `User not found`,
		},
		{
			name:           "Valid User",
			userID:         "1",
			body:           mockUpdatePasswordRequest,
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Empty String Password",
			userID: "1",
			body: models.UpdateUserPasswordRequest{
				Password: "",
			},
			expectedStatus: http.StatusBadRequest,
			expectedError:  `Error:Field validation for 'Password' failed on the 'required' tag`,
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {

			createTables()
			defer teardown()

			GenerateNewUser(t)

			reqBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPut, "/api/v1/restricted/users-update-password/"+tt.userID, strings.NewReader(string(reqBody)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api/v1/restricted/users-update-password/:uid")
			c.SetParamNames("uid")
			c.SetParamValues(tt.userID)

			err := handlers.ChangePassword(c)

			if err != nil {
				e.HTTPErrorHandler(err, c)
			}

			assert.Equal(t, tt.expectedStatus, rec.Code)
			if tt.expectedError != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedError)
			}
		})
	}
}
