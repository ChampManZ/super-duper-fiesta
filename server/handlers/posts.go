package handlers

import (
	"net/http"
	"server/config"
	"server/models"

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
	config.DB.Find(&posts)
	return c.JSON(http.StatusOK, posts)
}
