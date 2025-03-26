package repository

import (
	"time"
	"userService/model"

	"gorm.io/gorm"
)

type IPasswordResetRepository interface {
	CreateResetToken(email, token string, expiresAt time.Time) error
	ValidateResetToken(token string) (*model.PasswordReset, error)
	MarkResetTokenUsed(token string) error
	GetActivePasswordResetByEmail(email string) (*model.PasswordReset, error)
}

type passwordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) IPasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) CreateResetToken(email, token string, expiresAt time.Time) error {
	data := &model.PasswordReset{
		Email:     email,
		Token:     token,
		ExpiresAt: expiresAt,
		Used:      false,
	}
	return r.db.Create(data).Error
}

func (r *passwordResetRepository) ValidateResetToken(token string) (*model.PasswordReset, error) {
	var reset model.PasswordReset
	err := r.db.Where("token = ? AND used = false AND expires_at > ?", token, time.Now()).
		Limit(1).
		First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}

func (r *passwordResetRepository) MarkResetTokenUsed(token string) error {
	return r.db.Model(&model.PasswordReset{}).
		Where("token = ?", token).
		Update("used", true).Error
}

func (r *passwordResetRepository) GetActivePasswordResetByEmail(email string) (*model.PasswordReset, error) {
	var reset model.PasswordReset
	err := r.db.Where("email = ? AND used = ? AND expires_at > ?", email, false, time.Now()).First(&reset).Error
	if err != nil {
		return nil, err
	}
	return &reset, nil
}
