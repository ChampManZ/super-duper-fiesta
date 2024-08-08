package handlers

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func SwaggerHandler(c echo.Context) error {
	echoSwagger.WrapHandler(c)
	return nil
}
