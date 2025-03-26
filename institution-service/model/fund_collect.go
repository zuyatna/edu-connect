package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FundCollect struct {
	FundCollectID uuid.UUID      `json:"fund_collect_id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	PostID        uuid.UUID      `json:"post_id" gorm:"type:uuid; not null"`
	UserID        string         `json:"user_id" gorm:"type:varchar(255); not null"`
	UserName      string         `json:"user_name" gorm:"type:varchar(255); not null"`
	Amount        float64        `json:"amount" gorm:"type:float; not null"`
	TransactionID string         `json:"transaction_id" gorm:"type:varchar(255); not null"`
	CreatedAt     time.Time      `json:"created_at" gorm:"type:timestamp; not null; autoCreateTime"`
	UpdatedAt     time.Time      `json:"updated_at" gorm:"type:timestamp; not null; autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"type:timestamp"`
}

type FundCollectRequest struct {
	PostID        string  `json:"post_id"`
	UserID        string  `json:"user_id"`
	UserName      string  `json:"user_name"`
	Amount        float64 `json:"amount"`
	TransactionID string  `json:"transaction_id"`
}

type FundCollectResponse struct {
	FundCollectID string  `json:"fund_collect_id"`
	PostID        string  `json:"post_id"`
	UserID        string  `json:"user_id"`
	UserName      string  `json:"user_name"`
	Amount        float64 `json:"amount"`
	TransactionID string  `json:"transaction_id"`
}
