package handler

import (
	"github.com/labstack/echo/v4"
)

type errorResponse struct {
	Error string `json:"error"`
}

func newErrorResponse(c echo.Context, status int, err string) error {
	return c.JSON(status, errorResponse{Error: err})
}
