package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"type:varchar(50);not null" json:"name"`
	Email      string         `gorm:"type:varchar(100);unique;not null" json:"email"`
	Password   string         `gorm:"type:varchar(255);not null" json:"password"`
	Balance    float64        `json:"balance"`
	IsVerified bool           `json:"is_verified"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
