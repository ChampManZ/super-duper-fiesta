package handlers

import (
	"net/http"
	"server/config"
	"server/helpers"
	"server/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func PostComment(c echo.Context) error {
	request := new(models.Comment)
	if err := helpers.BindAndValidateRequest(c, request); err != nil {
		return err
	}

	userToken := c.Get(("user")).(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	userID := claims["user_id"]

	var user models.User
	if res := config.DB.First(&user, userID); res.Error != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "User not found"})
	}

	comment := models.Comment{
		CommentMSG: request.CommentMSG,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Users:      []models.User{user},
	}

	if result := config.DB.Create(&comment); result.Error != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to create comment. Please try again"})
	}

	return c.JSON(http.StatusCreated, comment)
}
