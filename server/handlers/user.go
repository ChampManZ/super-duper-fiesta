package handlers

import (
	"net/http"
	"server/config"
	"server/models"

	"github.com/labstack/echo/v4"
)

// GetUsers godoc
// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Router /users [get]
func GetUsers(c echo.Context) error {
	var users []models.User
	config.DB.Find(&users)
	return c.JSON(http.StatusOK, users)
}
