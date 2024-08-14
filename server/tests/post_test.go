package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/config"
	"server/handlers"
	"server/helpers"
	"server/models"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func createTestPost(t *testing.T, db *gorm.DB, user *models.User) *models.Post {
	post := models.Post{
		UserID:    user.UserID,
		Message:   "This is a test post",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := db.Create(&post).Error; err != nil {
		t.Fatalf("Failed to create test post: %v", err)
	}
	return &post
}

// ----------- Model Testing ----------- //
// Requirement:
// UserID must not be negative | uint is always positive and Golang is type strict
// Message must not be empty
// CreatedAt and UpdatedAt must not allow future dates
func TestPost(t *testing.T) {
	validPost := models.Post{
		Message:   "This is a test post",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	inValidPost := models.Post{
		Message:   "",
		CreatedAt: time.Now().AddDate(0, 0, 1),
		UpdatedAt: time.Now().AddDate(0, 0, 1),
	}

	assert.NotEmpty(t, validPost.Message)
	assert.False(t, validPost.CreatedAt.After(time.Now()))
	assert.False(t, validPost.UpdatedAt.After(time.Now()))

	assert.Empty(t, inValidPost.Message)
	assert.True(t, inValidPost.CreatedAt.After(time.Now()))
	assert.True(t, inValidPost.UpdatedAt.After(time.Now()))
}

// ----------- API Testing ----------- //
func TestGetPosts(t *testing.T) {
	createTables()
	defer teardown()

	e := echo.New()
	e.Validator = helpers.NewValidator()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/posts", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, handlers.GetPosts(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestGetPostsByID(t *testing.T) {
	createTables()
	defer teardown()

	mockUser := createTestUser(t, config.DB)
	mockPost := createTestPost(t, config.DB, mockUser)

	e := echo.New()
	e.Validator = helpers.NewValidator()

	pid := "1"

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/posts/:pid", nil)
	getReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	getRec := httptest.NewRecorder()
	getCtx := e.NewContext(getReq, getRec)
	getCtx.SetParamNames("pid")
	getCtx.SetParamValues(pid)

	if assert.NoError(t, handlers.GetPosts(getCtx)) {
		assert.Equal(t, http.StatusOK, getRec.Code)

		var resBody interface{}
		err := json.Unmarshal(getRec.Body.Bytes(), &resBody)
		assert.NoError(t, err)

		switch val := resBody.(type) {
		case []interface{}:
			for _, element := range val {
				if postMap, ok := element.(map[string]interface{}); ok {
					assert.Equal(t, mockPost.Message, postMap["post_message"])
				} else {
					panic("Failed test")
				}
			}
		default:
			t.Errorf("Unknown JSON structure")
			t.Fatal(resBody)
		}
	}
}

func TestCreatePost(t *testing.T) {
	createTables()
	defer teardown()

	userMock := createTestUser(t, config.DB)

	e := echo.New()
	e.Validator = helpers.NewValidator()

	tokenString := createJWTTokenTest(t, userMock.UserID)

	postJSON := `{"message": "This is a test post"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/restricted/posts", strings.NewReader(postJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	access_config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(models.JWTClaims)
		},
		SigningKey: []byte("testing_mock"),
	}
	jwtMiddleware := echojwt.WithConfig(access_config)

	if assert.NoError(t, jwtMiddleware(handlers.CreatePost)(c)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var createdPost models.Post
		err := json.Unmarshal(rec.Body.Bytes(), &createdPost)
		assert.NoError(t, err)
		assert.Equal(t, "This is a test post", createdPost.Message)
		assert.Equal(t, userMock.UserID, createdPost.UserID)
	}
}
