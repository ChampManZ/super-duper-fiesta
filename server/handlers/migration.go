package handlers

import (
	"server/config"

	"github.com/labstack/echo/v4"
)

func RunMigration(c echo.Context) error {
	return config.RunMigration(c)
}

func GetMigration(c echo.Context) error {
	return config.GetMigration(c)
}
