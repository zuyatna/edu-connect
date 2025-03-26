package model

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type Transaction struct {
	TransactionID primitive.ObjectID `json:"transaction_id" bson:"_id,omitempty"`
	UserID        string             `json:"user_id" bson:"user_id"`
	PostID        string             `json:"post_id" bson:"post_id"`
	UserEmail     string             `json:"user_email" bson:"user_email"`
	PaymentID     string             `json:"payment_id" gorm:"not null"`
	PaymentURL    string             `json:"payment_url" gorm:""`
	PaymentStatus string             `json:"payment_status" gorm:"default:'PENDING'"`
	Amount        float64            `json:"amount" gorm:"not null"`
	AccountNumber string             `json:"account_number" gorm:"not null"`
	AccountName   string             `json:"account_name" gorm:"not null"`
	CreatedAt     time.Time          `json:"created_at" gorm:"default:current_timestamp"`
	UpdatedAt     time.Time          `json:"updated_at" gorm:"default:current_timestamp"`
}

type TransactionRequest struct {
	PostID        string  `json:"post_id" bson:"post_id"`
	Amount        float64 `json:"amount"`
	AccountNumber string  `json:"account_number"`
	AccountName   string  `json:"account_name"`
}

type TransactionResponse struct {
	TransactionID string  `json:"transaction_id"`
	UserID        string  `json:"user_id"`
	PostID        string  `json:"post_id"`
	UserEmail     string  `json:"user_email"`
	PaymentID     string  `json:"payment_id"`
	PaymentURL    string  `json:"payment_url"`
	PaymentStatus string  `json:"payment_status"`
	Amount        float64 `json:"amount"`
	AccountNumber string  `json:"account_number"`
	AccountName   string  `json:"account_name"`
}

type User struct {
	ID    string `gorm:"primaryKey" json:"id"`
	Name  string `gorm:"type:varchar(50);not null" json:"name"`
	Email string `gorm:"type:varchar(100);unique;not null" json:"email"`
}

type Post struct {
	PostID       uuid.UUID `json:"post_id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title        string    `json:"title" gorm:"type:varchar(255); not null"`
	FundAchieved float64   `json:"fund_achieved" gorm:"type:float; default:0"`
}

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
