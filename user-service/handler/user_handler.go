package handler

import (
	"errors"
	"net/http"
	"strconv"
	"userService/model"
	"userService/usecase"
	"userService/utils"

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
	Token string `json:"token"`
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

// Register godoc
// @Summary Register a new user
// @Description Register user with name, email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param user body model.User true "User Register Input"
// @Success 201 {object} map[string]interface{}
// @Failure 400,409,500 {object} map[string]interface{}
// @Router /v1/register [post]
func (h *UserHandler) Register(c echo.Context) error {
	var user model.User

	if err := c.Bind(&user); err != nil {

		logger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Warn("Invalid request body for Register")

		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")

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

		return utils.ErrorResponse(c, statusCode, errorMessage)
	}

	logger.WithField("email", user.Email).Info("User registered successfully")

	return utils.SuccessResponse(c, http.StatusCreated, nil, "User register successful")
}

// Login godoc
// @Summary User login
// @Description Login with email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Login input"
// @Success 200 {object} LoginResponse
// @Failure 400,401,404,500 {object} map[string]interface{}
// @Router /v1/login [post]
func (h *UserHandler) Login(c echo.Context) error {
	var data LoginRequest

	if err := c.Bind(&data); err != nil {

		logger.Warn("Invalid request body for Login")

		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")

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

		return utils.ErrorResponse(c, statusCode, errorMessage)
	}

	logger.WithField("email", data.Email).Info("User logged in successfully")

	return utils.SuccessResponse(c, http.StatusOK, LoginResponse{
		Token: token,
	}, "Login successful")
}

// ForgotPassword godoc
// @Summary Reset user password
// @Description Send email and new password to reset user password
// @Tags Users
// @Accept json
// @Produce json
// @Param body body ForgotPasswordRequest true "Forgot Password Payload"
// @Success 200 {object} utils.APIResponse
// @Failure 400,404,500 {object} utils.APIResponse
// @Router /v1/forgot-password [post]
func (h *UserHandler) ForgotPassword(c echo.Context) error {
	var req ForgotPasswordRequest

	if err := c.Bind(&req); err != nil {
		logger.Warn("Invalid request body for ForgotPassword")
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
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

		return utils.ErrorResponse(c, statusCode, err.Error())

	}

	logger.WithField("email", req.Email).Info("Password reset successfully")

	return utils.SuccessResponse(c, http.StatusOK, nil, "Password has been reset")
}

// GetUserByID godoc
// @Summary Get user by ID
// @Description Retrieve a single user by their ID
// @Tags Users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.APIResponse
// @Failure 400,404,500 {object} utils.APIResponse
// @Router /v1/users/{id} [get]
func (h *UserHandler) GetUserByID(c echo.Context) error {
	idParam := c.Param("id")

	idUint, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		logger.WithField("id", idParam).Warn("Invalid user ID param")
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID")
	}

	user, err := h.userUseCase.GetByID(uint(idUint))
	if err != nil {
		if errors.Is(err, customErr.ErrLoginEmailNotFound) {
			return utils.ErrorResponse(c, http.StatusNotFound, "User not found")
		}
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Internal server error")
	}

	return utils.SuccessResponse(c, http.StatusOK, utils.ConvertToUserResponse(*user), "User fetched successfully")
}

// GetAllUsersPaginated godoc
// @Summary Get all users with pagination
// @Description Retrieve list of users with pagination metadata
// @Tags Users
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} utils.APIResponse
// @Failure 500 {object} utils.APIResponse
// @Router /v1/users [get]
func (h *UserHandler) GetAllUsersPaginated(c echo.Context) error {
	pageQuery := c.QueryParam("page")
	limitQuery := c.QueryParam("limit")

	page, err := strconv.Atoi(pageQuery)
	if err != nil || page <= 0 {
		page = 1
	}

	limit, err := strconv.Atoi(limitQuery)
	if err != nil || limit <= 0 {
		limit = 10
	}

	users, total, err := h.userUseCase.GetAllPaginated(page, limit)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch users")
	}

	userRes := utils.ConvertToUserResponseList(users)

	response := map[string]interface{}{
		"items": userRes,
		"pagination": map[string]interface{}{
			"page":      page,
			"limit":     limit,
			"totalData": total,
			"totalPage": int((total + int64(limit) - 1) / int64(limit)),
		},
	}

	return utils.SuccessResponse(c, http.StatusOK, response, "Users fetched successfully")
}

func (h *UserHandler) UpdateBalance(c echo.Context) error {
	var req UpdateBalanceRequest

	if err := c.Bind(&req); err != nil {
		logger.Warn("Invalid request body for UpdateBalance")
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
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

		return utils.ErrorResponse(c, statusCode, err.Error())
	}

	logger.WithFields(logrus.Fields{
		"email":   req.Email,
		"balance": req.Balance,
	}).Info("User balance updated successfully")

	return utils.SuccessResponse(c, http.StatusOK, nil, "Balance updated successfully")
}
