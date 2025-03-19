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
}

type userUseCase struct {
	userRepo repository.IUserRepository
}

var logger = logrus.New()

func NewUserUseCase(userRepo repository.IUserRepository) IUserUseCase {
	return &userUseCase{
		userRepo: userRepo,
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

func GenerateJWTToken(email string) (string, error) {

	JWTSecret := os.Getenv("JWT_SECRET")

	claims := jwt.MapClaims{
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

	token, err := GenerateJWTToken(user.Email)
	if err != nil {
		logger.WithField("email", email).Error("Failed to generate JWT token")
		return "", err
	}

	logger.WithField("email", email).Info("User logged in successfully")
	return token, nil
}
