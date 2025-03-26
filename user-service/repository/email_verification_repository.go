package repository

import (
	"time"
	"userService/model"

	"gorm.io/gorm"
)

type IVerificationRepository interface {
	CreateToken(email, token string, expiresAt time.Time) error
	ValidateToken(token string) (*model.EmailVerification, error)
	MarkTokenUsed(token string) error
	GetActiveVerificationByEmail(email string) (*model.EmailVerification, error)
}

type verificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) IVerificationRepository {
	return &verificationRepository{db: db}
}

func (r *verificationRepository) CreateToken(email, token string, expiresAt time.Time) error {
	data := &model.EmailVerification{
		Email:     email,
		Token:     token,
		ExpiresAt: expiresAt,
		Used:      false,
	}
	return r.db.Create(data).Error
}

func (r *verificationRepository) ValidateToken(token string) (*model.EmailVerification, error) {
	var ev model.EmailVerification
	err := r.db.Where("token = ? AND used = false AND expires_at > ?", token, time.Now()).
		First(&ev).Error
	if err != nil {
		return nil, err
	}
	return &ev, nil
}

func (r *verificationRepository) MarkTokenUsed(token string) error {
	return r.db.Model(&model.EmailVerification{}).
		Where("token = ?", token).
		Update("used", true).Error
}

func (r *verificationRepository) GetActiveVerificationByEmail(email string) (*model.EmailVerification, error) {
	var ev model.EmailVerification
	err := r.db.Where("email = ? AND used = false AND expires_at > ?", email, time.Now()).First(&ev).Error
	if err != nil {
		return nil, err
	}
	return &ev, nil
}
