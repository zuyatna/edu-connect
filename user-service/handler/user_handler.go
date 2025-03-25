package handler

import (
	"errors"
	"net/http"
	"userService/model"
	"userService/usecase"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"

	customErr "userService/error"
)

type UserHandler struct {
	userUseCase usecase.IUserUseCase
}

func NewUserHandler(userUseCase usecase.IUserUseCase) UserHandler {
	return UserHandler{
		userUseCase: userUseCase,
	}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

type UserInfoResponse struct {
	Username string `json:"username"`
}

type VerifyTokenResponse struct {
	Valid bool   `json:"valid"`
	Error string `json:"error,omitempty"`
}

type UserInfoAPIResponse struct {
	Data UserInfoResponse `json:"data"`
}

type ForgotPasswordRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type UpdateBalanceRequest struct {
	Email   string  `json:"email"`
	Balance float64 `json:"balance"`
}

var logger = logrus.New()

func (h *UserHandler) Register(c echo.Context) error {
	var user model.User

	if err := c.Bind(&user); err != nil {
		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Invalid request body for Register")
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "Invalid request body"},
		)
	}

	logger.WithField("email", user.Email).Info("Register request received")

	if err := h.userUseCase.Register(user); err != nil {

		var statusCode int
		var errorMessage string

		if errors.Is(err, customErr.ErrRegisterDuplicatedEmail) {

			statusCode = http.StatusConflict
			errorMessage = err.Error()

		} else if errors.Is(err, customErr.ErrRegisterEmailRequired) {

			statusCode = http.StatusBadRequest
			errorMessage = err.Error()

		} else if errors.Is(err, customErr.ErrRegisterNameRequired) {

			statusCode = http.StatusBadRequest
			errorMessage = err.Error()

		} else if errors.Is(err, customErr.ErrRegisterPasswordRequired) {

			statusCode = http.StatusBadRequest
			errorMessage = err.Error()

		} else if errors.Is(err, customErr.ErrRegisterInvalidPassword) {

			statusCode = http.StatusBadRequest
			errorMessage = err.Error()

		} else {

			statusCode = http.StatusInternalServerError
			errorMessage = err.Error()

		}

		logger.WithFields(logrus.Fields{
			"email": user.Email,
			"error": errorMessage,
		}).Error("Register failed")

		return c.JSON(
			statusCode,
			map[string]string{"error": errorMessage},
		)
	}

	logger.WithField("email", user.Email).Info("User registered successfully")

	return c.JSON(
		http.StatusCreated,
		map[string]string{"message": "User register successful"},
	)
}

func (h *UserHandler) Login(c echo.Context) error {
	var data LoginRequest

	if err := c.Bind(&data); err != nil {
		logger.Warn("Invalid request body for Login")
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{"error": "Invalid request body"},
		)
	}

	logger.WithField("email", data.Email).Info("Login request received")

	token, err := h.userUseCase.Login(data.Email, data.Password)

	if err != nil {

		var statusCode int
		var errorMessage string

		if errors.Is(err, customErr.ErrRegisterEmailRequired) {

			statusCode = http.StatusBadRequest
			errorMessage = err.Error()

		} else if errors.Is(err, customErr.ErrRegisterPasswordRequired) {

			statusCode = http.StatusBadRequest
			errorMessage = err.Error()

		} else if errors.Is(err, customErr.ErrLoginEmailNotFound) {

			statusCode = http.StatusNotFound
			errorMessage = err.Error()

		} else if errors.Is(err, customErr.ErrLoginInvalidPassword) {

			statusCode = http.StatusUnauthorized
			errorMessage = err.Error()

		} else {

			statusCode = http.StatusInternalServerError
			errorMessage = err.Error()

		}

		logger.WithFields(logrus.Fields{
			"email": data.Email,
			"error": errorMessage,
		}).Warn("Login failed")

		return c.JSON(statusCode, ErrorResponse{Error: errorMessage})
	}

	logger.WithField("email", data.Email).Info("User logged in successfully")

	return c.JSON(http.StatusOK, LoginResponse{
		Message: "Login successful",
		Token:   token,
	})
}

func (h *UserHandler) VerifyEmail(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		logger.Warn("Verify email failed: email query param is empty")
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Email is required"})
	}

	logger.WithField("email", email).Info("Verification request received")

	err := h.userUseCase.UpdateIsVerified(email)
	if err != nil {
		var statusCode int
		if errors.Is(err, customErr.ErrLoginEmailNotFound) {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusInternalServerError
		}

		logger.WithFields(logrus.Fields{
			"email": email,
			"error": err.Error(),
		}).Error("Email verification failed")

		return c.JSON(statusCode, ErrorResponse{Error: err.Error()})
	}

	logger.WithField("email", email).Info("User email verified successfully")
	return c.JSON(http.StatusOK, map[string]string{"message": "Email verified successfully"})
}

func (h *UserHandler) ForgotPassword(c echo.Context) error {
	var req ForgotPasswordRequest

	if err := c.Bind(&req); err != nil {
		logger.Warn("Invalid request body for ForgotPassword")
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
	}

	logger.WithField("email", req.Email).Info("Forgot password request received")

	err := h.userUseCase.ForgotPassword(req.Email, req.NewPassword)
	if err != nil {
		var statusCode int
		if errors.Is(err, customErr.ErrRegisterInvalidEmail) || errors.Is(err, customErr.ErrRegisterInvalidPassword) {
			statusCode = http.StatusBadRequest
		} else if errors.Is(err, customErr.ErrLoginEmailNotFound) {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusInternalServerError
		}

		logger.WithFields(logrus.Fields{
			"email": req.Email,
			"error": err.Error(),
		}).Error("Forgot password failed")

		return c.JSON(statusCode, ErrorResponse{Error: err.Error()})
	}

	logger.WithField("email", req.Email).Info("Password reset successfully")
	return c.JSON(http.StatusOK, map[string]string{"message": "Password has been reset"})
}

func (h *UserHandler) UpdateBalance(c echo.Context) error {
	var req UpdateBalanceRequest

	if err := c.Bind(&req); err != nil {
		logger.Warn("Invalid request body for UpdateBalance")
		return c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
	}

	logger.WithFields(logrus.Fields{
		"email":   req.Email,
		"balance": req.Balance,
	}).Info("Update balance request received")

	err := h.userUseCase.UpdateBalance(req.Email, req.Balance)
	if err != nil {
		var statusCode int
		if errors.Is(err, customErr.ErrRegisterInvalidEmail) {
			statusCode = http.StatusBadRequest
		} else if errors.Is(err, customErr.ErrLoginEmailNotFound) {
			statusCode = http.StatusNotFound
		} else {
			statusCode = http.StatusInternalServerError
		}

		logger.WithFields(logrus.Fields{
			"email":   req.Email,
			"balance": req.Balance,
			"error":   err.Error(),
		}).Error("Update balance failed")

		return c.JSON(statusCode, ErrorResponse{Error: err.Error()})
	}

	logger.WithFields(logrus.Fields{
		"email":   req.Email,
		"balance": req.Balance,
	}).Info("User balance updated successfully")

	return c.JSON(http.StatusOK, map[string]string{"message": "Balance updated successfully"})
}
