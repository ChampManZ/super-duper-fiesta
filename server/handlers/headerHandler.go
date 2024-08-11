package handlers

import "github.com/labstack/echo/v4"

func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "DWTakeHome/1.0")
		return next(c)
	}
}
