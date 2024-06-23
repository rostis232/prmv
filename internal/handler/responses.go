package handler

import (
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(c echo.Context, status int, err string) error {
	return c.JSON(status, ErrorResponse{Error: err})
}
