package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
