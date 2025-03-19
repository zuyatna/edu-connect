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

type UserInfoAPIResponse struct {
	Data UserInfoResponse `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
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
