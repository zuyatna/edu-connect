package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Institution struct {
	InstitutionID uuid.UUID      `json:"institution_id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name          string         `json:"name" gorm:"type:varchar(255); not null"`
	Email         string         `json:"email" gorm:"type:varchar(255); not null; unique"`
	Password      string         `json:"password" gorm:"type:varchar(255); not null"`
	Address       string         `json:"address" gorm:"type:varchar(255); not null"`
	Phone         string         `json:"phone" gorm:"type:varchar(255); not null"`
	Website       string         `json:"website" gorm:"type:varchar(255)"`
	CreatedAt     time.Time      `json:"created_at" gorm:"type:timestamp; not null; autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"type:timestamp; not null; autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp"`
}

func (u *Institution) CompareHashAndPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

type InstitutionLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type InstitutionRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Address  string `json:"address"`
	Phone    string `json:"phone"`
	Website  string `json:"website"`
}

type InstitutionResponse struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Website string `json:"website"`
}

type InstitutionToken struct {
	Token string `json:"token"`
}

type InstitutionDeleteResponse struct {
	Message string `json:"message"`
}
