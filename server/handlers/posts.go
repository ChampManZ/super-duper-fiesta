package handlers

import (
	"net/http"
	"server/config"
	"server/helpers"
	"server/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// GetPosts godoc
// @Summary Get all posts
// @Description Get all posts
// @Tags posts
// @Accept json
// @Produce json
// @Success 200 {object} models.Post
// @Router /posts [get]
func GetPosts(c echo.Context) error {
	var posts []models.GetPublicPostsRequest
	if result := config.DB.Table("posts").Select("posts.post_id, users.username, users.firstname, users.surname, posts.message, posts.created_at, posts.updated_at").Joins("inner join users on users.user_id = posts.user_id").Scan(&posts); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get posts"})
	}

	return c.JSON(http.StatusOK, posts)
}

func CreatePost(c echo.Context) error {
	request := new(models.Post)
	if err := helpers.BindAndValidateRequest(c, request); err != nil {
		return err
	}

	// Reference: https://github.com/labstack/echo/issues/1504
	user := c.Get(("user"))
	token := user.(*jwt.Token)
	claims := token.Claims.(*models.JWTClaims)

	userID := claims.UserID

	post := models.Post{
		Message: request.Message,
		UserID:  userID,
	}

	if result := config.DB.Create(&post); result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to create post. Please try again"})
	}

	return c.JSON(http.StatusCreated, post)
}
