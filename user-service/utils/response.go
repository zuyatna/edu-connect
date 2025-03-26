package utils

import (
	"github.com/labstack/echo/v4"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

func SuccessResponse(c echo.Context, statusCode int, data interface{}, message string) error {
	return c.JSON(statusCode, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, statusCode int, errMsg string) error {
	return c.JSON(statusCode, APIResponse{
		Status:  "error",
		Message: errMsg,
		Data:    nil,
	})
}
