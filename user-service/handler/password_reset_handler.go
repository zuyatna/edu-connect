package handler

import (
	"errors"
	"net/http"
	"userService/usecase"

	customErr "userService/error"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type PasswordResetHandler struct {
	resetUC usecase.IPasswordResetUseCase
}

func NewPasswordResetHandler(resetUC usecase.IPasswordResetUseCase) *PasswordResetHandler {
	return &PasswordResetHandler{
		resetUC: resetUC,
	}
}

type RequestResetPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	NewPassword string `json:"new_password"`
}

func (h *PasswordResetHandler) RequestResetPassword(c echo.Context) error {
	var req RequestResetPasswordRequest

	if err := c.Bind(&req); err != nil {
		logrus.Warn("Invalid request body for RequestResetPassword")
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
	}

	logrus.WithField("email", req.Email).Info("Password reset request received")

	err := h.resetUC.RequestReset(req.Email)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, customErr.ErrRegisterEmailRequired) || errors.Is(err, customErr.ErrLoginEmailNotFound) {
			statusCode = http.StatusBadRequest
		}
		logrus.WithError(err).WithField("email", req.Email).Error("Password reset request failed")
		return c.JSON(statusCode, ErrorResponse{Error: err.Error()})
	}

	logrus.WithField("email", req.Email).Info("Password reset token sent successfully")
	return c.JSON(http.StatusOK, map[string]string{"message": "Password reset link sent to email"})
}

func (h *PasswordResetHandler) ResetPassword(c echo.Context) error {

	token := c.QueryParam("token")

	var req ResetPasswordRequest

	if err := c.Bind(&req); err != nil {
		logrus.Warn("Invalid request body for ResetPassword")
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
	}

	logrus.WithField("token", token).Info("Reset password execution started")

	err := h.resetUC.ResetPassword(token, req.NewPassword)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if errors.Is(err, customErr.ErrVerificationTokenInvalid) || errors.Is(err, customErr.ErrRegisterInvalidPassword) {
			statusCode = http.StatusBadRequest
		}
		logrus.WithError(err).WithField("token", token).Error("Reset password failed")
		return c.JSON(statusCode, ErrorResponse{Error: err.Error()})
	}

	logrus.WithField("token", token).Info("Password reset successfully")
	return c.JSON(http.StatusOK, map[string]string{"message": "Password has been reset successfully"})
}
