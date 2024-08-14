package tests

import (
	"encoding/json"
	"fmt"
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
)

// Requirement:
// UserID must not be negative | uint is always positive and Golang is type strict
// CommentMSG must not be empty
// CreatedAt and UpdatedAt must not allow future dates
func TestComment(t *testing.T) {
	validComment := models.Comment{
		CommentMSG: "This is a test comment",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	assert.NotEmpty(t, validComment.CommentMSG)
	assert.False(t, validComment.CreatedAt.After(time.Now()))
	assert.False(t, validComment.UpdatedAt.After(time.Now()))
}

func TestCreateComment(t *testing.T) {
	createTables()
	defer teardown()

	userMock := createTestUser(t, config.DB)
	postMock := createTestPost(t, config.DB, userMock)

	e := echo.New()
	e.Validator = helpers.NewValidator()

	tokenString := createJWTTokenTest(t, userMock.UserID)

	commentJSON := `{"post_id":` + fmt.Sprint(postMock.PostID) + `,"comment_msg":"This is a test comment"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/restricted/comments", strings.NewReader(commentJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, "Bearer "+tokenString)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	access_config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(models.JWTClaims)
		},
		SigningKey: []byte("testing_mock"),
	}
	jwtMiddleware := echojwt.WithConfig(access_config)

	if assert.NoError(t, jwtMiddleware(handlers.CreateComment)(ctx)) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var createComment models.Comment
		err := json.Unmarshal(rec.Body.Bytes(), &createComment)

		assert.NoError(t, err)
		assert.Equal(t, postMock.PostID, createComment.PostID)
		assert.Equal(t, postMock.Message, createComment.CommentMSG)
	}
}
