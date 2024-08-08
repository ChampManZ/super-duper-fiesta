package handlers

import (
	"net/http"
	"server/config"
	"server/models"

	"github.com/labstack/echo/v4"
)

func GetUsers(c echo.Context) error {
	var users []models.User
	config.DB.Find(&users)
	return c.JSON(http.StatusOK, users)
}
