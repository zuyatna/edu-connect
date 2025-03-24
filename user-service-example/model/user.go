package model

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	UserID      uuid.UUID      `json:"user_id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string         `json:"name" gorm:"type:varchar(255); not null"`
	Email       string         `json:"email" gorm:"type:varchar(255); not null; unique"`
	Password    string         `json:"password" gorm:"type:varchar(255); not null"`
	DonateCount float64        `json:"donate_count" gorm:"type:float; default:0"`
	CreatedAt   time.Time      `json:"created_at" gorm:"type:timestamp; not null; autoCreateTime"`
	UpdatedAt   time.Time      `json:"updated_at" gorm:"type:timestamp; not null; autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp"`
}

func (u *User) CompareHashAndPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

type UserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserToken struct {
	Token string `json:"token"`
}
