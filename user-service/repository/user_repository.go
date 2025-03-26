package repository

import (
	"errors"
	"fmt"
	"userService/model"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Register(user *model.User) error
	Login(username, password string) (*model.User, error)
	UpdateIsVerified(email string, verified bool) error
	GetByEmail(email string) (*model.User, error)
	UpdatePasswordByEmail(email, newPassword string) error
	UpdateBalanceByEmail(email string, balance float64) error
	GetByID(id uint) (*model.User, error)
	GetAllPaginated(page int, limit int) ([]model.User, int64, error)
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

func (r *userRepository) GetAllPaginated(page int, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	offset := (page - 1) * limit

	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		logger.WithError(err).Error("Failed to count users")
		return nil, 0, err
	}

	if err := r.db.Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		logger.WithError(err).Error("Failed to fetch users with pagination")
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.WithField("user_id", id).Warn("Get user failed: ID not found")
			return nil, errors.New("user not found")
		}
		logger.WithField("user_id", id).Error("Failed to get user by ID")
		return nil, err
	}

	logger.WithField("user_id", id).Info("User retrieved successfully by ID")
	return &user, nil
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

func (r *userRepository) UpdateIsVerified(email string, verified bool) error {
	result := r.db.Model(&model.User{}).Where("email = ?", email).Update("is_verified", verified)
	if result.Error != nil {
		logger.WithFields(logrus.Fields{
			"email": email,
			"error": result.Error,
		}).Error("Failed to update is_verified status")
		return result.Error
	}
	if result.RowsAffected == 0 {
		logger.WithField("email", email).Warn("No user found to update verification")
		return gorm.ErrRecordNotFound
	}

	logger.WithField("email", email).Info("User verification status updated")
	return nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) UpdatePasswordByEmail(email, newPassword string) error {
	hashed, err := hashPassword(newPassword)
	if err != nil {
		logger.WithError(err).Error("Failed to hash new password")
		return err
	}

	result := r.db.Model(&model.User{}).Where("email = ?", email).Update("password", hashed)
	if result.Error != nil {
		logger.WithFields(logrus.Fields{
			"email": email,
			"error": result.Error,
		}).Error("Failed to update password")
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	logger.WithField("email", email).Info("User password updated successfully")
	return nil
}

func (r *userRepository) UpdateBalanceByEmail(email string, balance float64) error {

	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	if user.Balance+balance < 0 {
		return fmt.Errorf("insufficient balance")
	}

	result := r.db.Model(&model.User{}).
		Where("email = ?", email).
		Update("balance", gorm.Expr("balance + ?", balance))

	if result.Error != nil {
		logger.WithFields(logrus.Fields{
			"email": email,
			"error": result.Error,
		}).Error("Failed to update balance")
		return result.Error
	}

	if result.RowsAffected == 0 {
		logger.WithField("email", email).Warn("No user found to update balance")
		return gorm.ErrRecordNotFound
	}

	logger.WithFields(logrus.Fields{
		"email":   email,
		"balance": balance,
	}).Info("User balance updated successfully")
	return nil
}
