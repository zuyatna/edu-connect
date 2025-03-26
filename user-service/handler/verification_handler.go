package handler

import (
	"net/http"
	"userService/usecase"
	"userService/utils"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	customErr "userService/error"
)

type requestResendVerify struct {
	Email string `json:"email"`
}

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
		return utils.ErrorResponse(c, http.StatusBadRequest, "Token is required")
	}

	logrus.WithField("token", token).Info("Verification request received")

	err := h.verificationUC.VerifyToken(token)
	if err != nil {
		if err == customErr.ErrVerificationTokenInvalid {
			logrus.WithField("token", token).Warn("Verification failed: invalid or expired token")
			return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid or expired token")
		}

		logrus.WithField("token", token).Error("Verification failed: internal error")
		return utils.ErrorResponse(c, http.StatusBadRequest, "Internal server error")
	}

	logrus.WithField("token", token).Info("User verified successfully")
	return utils.SuccessResponse(c, http.StatusOK, nil, "Email verified successfully")

}

func (h *VerificationHandler) ResendVerification(c echo.Context) error {
	var req requestResendVerify
	if err := c.Bind(&req); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
	}

	if req.Email == "" {
		logrus.Warn("Resend verification failed: email is missing")
		return utils.ErrorResponse(c, http.StatusBadRequest, "Email is required")
	}

	logrus.WithField("email", req.Email).Info("Resend verification request received")

	err := h.verificationUC.ResendVerification(req.Email)
	if err != nil {
		if err == customErr.ErrLoginEmailNotFound {
			logrus.WithField("email", req.Email).Warn("Resend verification failed: Email not found")
			return utils.ErrorResponse(c, http.StatusBadRequest, "Email not found")
		}
		if err == customErr.ErrVerificationTokenStillValid {
			logrus.WithField("email", req.Email).Warn("Resend verification denied: Active token exists")
			return utils.ErrorResponse(c, http.StatusBadRequest, "Existing verification token still valid")
		}

		logrus.WithField("email", req.Email).Error("Resend verification failed: internal error")
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
	}

	logrus.WithField("email", req.Email).Info("Verification email resent successfully")
	return utils.SuccessResponse(c, http.StatusOK, nil, "Verification email resent successfully")
}
