package handlers

import (
	"net/http"
	"server/config"
	"server/helpers"
	"server/models"
	"time"

	"github.com/labstack/echo/v4"
)

// GetComments godoc
// @Summary Get all comments
// @Description Get all comments
// @Tags comments
// @Accept json
// @Produce json
// @Success 200 {array} models.Comment
// @Failure 500 {object} map[string]string "Failed to retrieve comments"
// @Router /api/v1/admin/comments [get]
func GetComments(c echo.Context) error {
	var comments []models.Comment
	config.DB.Find(&comments)
	return c.JSON(http.StatusOK, comments)
}

// PostComment godoc
// @Summary Create a comment
// @Description Create a new comment
// @Tags comments
// @Accept json
// @Produce json
// @Param comment body models.Comment true "Comment object that needs to be created"
// @Success 201 {object} models.Comment
// @Failure 400 {object} map[string]string "Invalid input or failed to create comment"
// @Router /api/v1/restricted/comments [post]
func PostComment(c echo.Context) error {
	request := new(models.Comment)
	if err := helpers.BindAndValidateRequest(c, request); err != nil {
		return err
	}

	// Response to test: { "CommentMSG": "something" }
	comment := models.Comment{
		CommentMSG: request.CommentMSG,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if result := config.DB.Create(&comment); result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to create comment. Please try again"})
	}

	return c.JSON(http.StatusCreated, comment)
}
