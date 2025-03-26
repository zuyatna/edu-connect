package usecase

import (
	"os"
	"regexp"
	"strings"
	"time"
	"userService/model"
	"userService/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"

	customErr "userService/error"
)

type IUserUseCase interface {
	Register(user model.User) error
	Login(email, password string) (string, error)
	ForgotPassword(email, newPassword string) error
	UpdateIsVerified(email string) error
	GetByEmail(email string) (*model.User, error)
	UpdateBalance(email string, balance float64) error
	GetByID(id uint) (*model.User, error)
	GetAllPaginated(page int, limit int) ([]model.User, int64, error)
}

type userUseCase struct {
	userRepo            repository.IUserRepository
	verificationUsecase IVerificationUseCase
}

var logger = logrus.New()

func NewUserUseCase(userRepo repository.IUserRepository, verificationUC IVerificationUseCase) IUserUseCase {
	return &userUseCase{
		userRepo:            userRepo,
		verificationUsecase: verificationUC,
	}
}

func isValidPassword(password string) bool {
	regex := `^.{8,}$`
	matched, _ := regexp.MatchString(regex, password)
	return matched
}

func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(regex, email)
	return matched
}

func (u *userUseCase) GetAllPaginated(page int, limit int) ([]model.User, int64, error) {
	users, total, err := u.userRepo.GetAllPaginated(page, limit)
	if err != nil {
		logger.Error("GetAllPaginated failed")
		return nil, 0, customErr.ErrInternalServer
	}

	return users, total, nil
}

func (u *userUseCase) GetByID(id uint) (*model.User, error) {
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		if err.Error() == "user not found" {
			logger.WithField("id", id).Warn("Get user by ID failed: not found")
			return nil, customErr.ErrLoginEmailNotFound
		}

		logger.WithField("id", id).Error("Get user by ID failed: internal error")
		return nil, customErr.ErrInternalServer
	}

	logger.WithField("id", id).Info("User fetched successfully by ID")
	return user, nil
}

func GenerateJWTToken(name, email string) (string, error) {

	JWTSecret := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
		"name":  name,
		"email": email,
		"exp":   time.Now().Add(time.Minute * 60).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(JWTSecret))
	if err != nil {
		logger.WithError(err).Error("Failed to generate JWT token")
		return "", err
	}

	logger.WithField("username", email).Info("JWT token generated successfully")

	return signedToken, nil

}

func (u *userUseCase) Register(user model.User) error {

	if user.Email == "" {
		logger.Warn("Register failed: Email is empty")
		return customErr.ErrRegisterEmailRequired
	}

	if user.Name == "" {
		logger.Warn("Register failed: Name is empty")
		return customErr.ErrRegisterNameRequired
	}

	if user.Password == "" {
		logger.Warn("Register failed: Password is empty")
		return customErr.ErrRegisterPasswordRequired
	}

	if !isValidEmail(user.Email) {
		logger.WithField("email", user.Email).Warn("Register failed: Invalid email format")
		return customErr.ErrRegisterInvalidEmail
	}

	if !isValidPassword(user.Password) {
		logger.WithField("email", user.Email).Warn("Register failed: Password does not meet criteria")
		return customErr.ErrRegisterInvalidPassword
	}

	err := u.userRepo.Register(&user)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			logger.WithField("email", user.Email).Warn("Register failed: Duplicate email")
			return customErr.ErrRegisterDuplicatedEmail
		}

		logger.WithFields(logrus.Fields{
			"email": user.Email,
			"error": err.Error(),
		}).Error("Register failed: Internal server error")
		return customErr.ErrInternalServer
	}

	logger.WithField("email", user.Email).Info("User registered successfully")

	_ = u.verificationUsecase.GenerateVerification(user.Email)

	return nil
}

func (u *userUseCase) Login(email, password string) (string, error) {

	if email == "" {
		logger.Warn("Login failed: Email or password is empty")
		return "", customErr.ErrRegisterEmailRequired
	}

	if password == "" {
		logger.Warn("Login failed: Email or password is empty")
		return "", customErr.ErrRegisterPasswordRequired
	}

	user, err := u.userRepo.Login(email, password)
	if err != nil {
		if strings.Contains(err.Error(), "email doesn't exist") {

			logger.WithField("email", email).Warn("Login failed: Email not found")
			return "nil", customErr.ErrLoginEmailNotFound

		} else if strings.Contains(err.Error(), "wrong password") {

			logger.WithField("email", email).Warn("Login failed: Wrong password")
			return "nil", customErr.ErrLoginInvalidPassword

		}

		logger.WithFields(logrus.Fields{
			"email": email,
			"error": err.Error(),
		}).Error("Login failed: Internal server error")

		return "nil", customErr.ErrInternalServer
	}

	token, err := GenerateJWTToken(user.Name, user.Email)
	if err != nil {
		logger.WithField("email", email).Error("Failed to generate JWT token")
		return "", err
	}

	logger.WithField("email", email).Info("User logged in successfully")
	return token, nil
}

func (u *userUseCase) UpdateIsVerified(email string) error {
	err := u.userRepo.UpdateIsVerified(email, true)
	if err != nil {
		if err == customErr.ErrLoginEmailNotFound {
			logger.WithField("email", email).Warn("Verification failed: Email not found")
			return customErr.ErrLoginEmailNotFound
		}
		logger.WithField("email", email).Error("Verification failed: Internal error")
		return customErr.ErrInternalServer
	}

	logger.WithField("email", email).Info("User verification updated successfully")
	return nil
}

func (u *userUseCase) ForgotPassword(email, newPassword string) error {
	if !isValidEmail(email) {
		logger.WithField("email", email).Warn("Forgot password failed: Invalid email")
		return customErr.ErrRegisterInvalidEmail
	}

	if !isValidPassword(newPassword) {
		logger.WithField("email", email).Warn("Forgot password failed: Weak password")
		return customErr.ErrRegisterInvalidPassword
	}

	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		logger.WithField("email", email).Warn("Forgot password failed: Email not found")
		return customErr.ErrLoginEmailNotFound
	}

	err = u.userRepo.UpdatePasswordByEmail(user.Email, newPassword)
	if err != nil {
		logger.WithField("email", email).Error("Failed to update password")
		return customErr.ErrInternalServer
	}

	logger.WithField("email", email).Info("User password updated successfully")
	return nil
}

func (u *userUseCase) UpdateBalance(email string, balance float64) error {

	if !isValidEmail(email) {
		logger.WithField("email", email).Warn("Update balance failed: Invalid email format")
		return customErr.ErrRegisterInvalidEmail
	}

	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		logger.WithField("email", email).Warn("Update balance failed: Email not found")
		return customErr.ErrLoginEmailNotFound
	}

	err = u.userRepo.UpdateBalanceByEmail(user.Email, balance)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"email":   email,
			"balance": balance,
			"error":   err.Error(),
		}).Error("Failed to update balance")
		return customErr.ErrInternalServer
	}

	logger.WithFields(logrus.Fields{
		"email":   email,
		"balance": balance,
	}).Info("User balance updated successfully")
	return nil
}

func (u *userUseCase) GetByEmail(email string) (*model.User, error) {

	if email == "" {
		logger.Warn("Get user failed: Email is empty")
		return nil, customErr.ErrRegisterEmailRequired
	}

	if !isValidEmail(email) {
		logger.WithField("email", email).Warn("Get user failed: Invalid email format")
		return nil, customErr.ErrRegisterInvalidEmail
	}

	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		if err == customErr.ErrLoginEmailNotFound {
			logger.WithField("email", email).Warn("Get user failed: Email not found")
			return nil, customErr.ErrLoginEmailNotFound
		}

		logger.WithFields(logrus.Fields{
			"email": email,
			"error": err.Error(),
		}).Error("Get user failed: Internal server error")

		return nil, customErr.ErrInternalServer
	}

	logger.WithField("email", email).Info("User retrieved successfully")
	return user, nil

}
