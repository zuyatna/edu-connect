package repository

import (
	"errors"
	"userService/model"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Register(user *model.User) error
	Login(username, password string) (*model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

var logger = logrus.New()

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{
		db: db,
	}
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (r *userRepository) Register(user *model.User) error {

	processHash, err := hashPassword(user.Password)
	if err != nil {
		logger.WithError(err).Error("Failed to hash password")
		return err
	}

	user.Password = processHash

	res := r.db.Create(&user)
	if res.Error != nil {
		logger.WithFields(logrus.Fields{
			"email": user.Email,
			"error": res.Error,
		}).Error("Failed to register user")
		return res.Error
	}

	logger.WithField("email", user.Email).Info("User registered successfully")
	return nil

}

func (r *userRepository) Login(email, password string) (*model.User, error) {

	var user model.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.WithField("email", email).Warn("Login failed: Email not found")
			return nil, errors.New("email doesn't exist")
		}
		logger.WithError(err).Error("Database error during login")
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		logger.WithField("email", email).Warn("Login failed: Wrong password")
		return nil, errors.New("wrong password")
	}

	logger.WithField("email", email).Info("User logged in successfully")
	return &user, nil

}
