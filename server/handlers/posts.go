package handlers

import (
	"net/http"
	"server/config"
	"server/helpers"
	"server/models"
	"strconv"

	"github.com/golang-jwt/jwt"
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
	var posts []models.Post
	if result := config.DB.Find(&posts); result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to get posts"})
	}

	return c.JSON(http.StatusOK, posts)
}

func CreatePost(c echo.Context) error {
	request := new(models.Post)
	if err := helpers.BindAndValidateRequest(c, request); err != nil {
		return err
	}

	user := c.Get(("user")).(*jwt.Token)
	claims := user.Claims.(jwt.StandardClaims)
	userID, err := strconv.ParseUint(claims.Id, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to get user ID"})
	}

	post := models.Post{
		Message: request.Message,
		UserID:  uint(userID),
	}

	if result := config.DB.Create(&post); result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to create post. Please try again"})
	}

	return c.JSON(http.StatusCreated, post)
}
