package helpers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func NewValidator() *CustomValidator {
	return &CustomValidator{validator: validator.New()}
}

func BindAndValidateRequest(c echo.Context, req interface{}) error {
	// Bind request data
	// Note: Should return an error instead of JSON failed response
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Invalid request data"})
	}

	// Validate request
	if err := c.Validate(req); err != nil {
		return err
	}

	return nil
}
