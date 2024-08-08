package handlers

import (
	"net/http"
	"server/config"
	"server/models"

	"github.com/labstack/echo/v4"
)

// GetComments godoc
// @Summary Get all comments
// @Description Get all comments
// @Tags comments
// @Accept json
// @Produce json
// @Success 200 {object} models.Comment
// @Router /comments [get]
func GetComments(c echo.Context) error {
	var comments []models.Comment
	config.DB.Find(&comments)
	return c.JSON(http.StatusOK, comments)
}

// GetCommentUser godoc
// @Summary Get all user comments
// @Description Get all user comments
// @Tags commentUser
// @Accept json
// @Produce json
// @Success 200 {object} models.CommentUser
// @Router /commentuser [get]
func GetCommentUser(c echo.Context) error {
	var commentUser []models.CommentUser
	config.DB.Find(&commentUser)
	return c.JSON(http.StatusOK, commentUser)
}
