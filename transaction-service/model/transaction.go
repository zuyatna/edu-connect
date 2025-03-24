package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Transaction struct {
	TransactionID primitive.ObjectID `json:"transaction_id" bson:"_id,omitempty"`
	UserID        string             `json:"user_id" bson:"user_id"`
	PostID        string             `json:"post_id" bson:"post_id"`
	PaymentID     string             `json:"payment_id" bson:"payment_id"`
	Amount        float64            `json:"amount" bson:"amount"`
	AccountNumber string             `json:"account_number" bson:"account_number"`
	AccountName   string             `json:"account_name" bson:"account_name"`
	CreatedAt     string             `json:"created_at" bson:"created_at"`
}

type TransactionRequest struct {
	PostID        string  `json:"post_id" bson:"post_id"`
	Amount        float64 `json:"amount"`
	AccountNumber string  `json:"account_number"`
	AccountName   string  `json:"account_name"`
	HideName      bool    `json:"hide_name"`
}
