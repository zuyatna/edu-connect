package utils

import (
	"userService/model"

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

func ConvertToUserResponse(user model.User) model.UserResponse {
	return model.UserResponse{
		UserID:     user.UserID,
		Name:       user.Name,
		Email:      user.Email,
		IsVerified: user.IsVerified,
	}
}

func ConvertToUserResponseList(users []model.User) []model.UserResponse {
	var res []model.UserResponse
	for _, user := range users {
		res = append(res, ConvertToUserResponse(user))
	}
	return res
}
