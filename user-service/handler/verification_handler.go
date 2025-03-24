package handler

import (
	"net/http"
	"userService/usecase"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	customErr "userService/error"
)

type VerificationHandler struct {
	verificationUC usecase.IVerificationUseCase
}

func NewVerificationHandler(verificationUC usecase.IVerificationUseCase) *VerificationHandler {
	return &VerificationHandler{
		verificationUC: verificationUC,
	}
}

func (h *VerificationHandler) Verify(c echo.Context) error {

	token := c.QueryParam("token")

	if token == "" {
		logrus.Warn("Verification failed: token is missing")
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Token is required"})
	}

	logrus.WithField("token", token).Info("Verification request received")

	err := h.verificationUC.VerifyToken(token)
	if err != nil {
		if err == customErr.ErrVerificationTokenInvalid {
			logrus.WithField("token", token).Warn("Verification failed: invalid or expired token")
			return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid or expired token"})
		}

		logrus.WithField("token", token).Error("Verification failed: internal error")
		return c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Internal server error"})
	}

	logrus.WithField("token", token).Info("User verified successfully")
	return c.JSON(http.StatusOK, map[string]string{"message": "Email verified successfully"})

}
