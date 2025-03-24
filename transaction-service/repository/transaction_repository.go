package repository

import (
	"context"
	"time"

	"github.com/zuyatna/edu-connect/transaction-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ITransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
}

type TransactionRepository struct {
	transactionCollection *mongo.Collection
}

func NewTransactionRepository(db *mongo.Database) *TransactionRepository {
	return &TransactionRepository{
		transactionCollection: db.Collection("transactions"),
	}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	doc := bson.D{
		{Key: "user_id", Value: transaction.UserID},
		{Key: "post_id", Value: transaction.PostID},
		{Key: "payment_id", Value: transaction.PaymentID},
		{Key: "amount", Value: transaction.Amount},
		{Key: "account_number", Value: transaction.AccountNumber},
		{Key: "account_name", Value: transaction.AccountName},
		{Key: "created_at", Value: time.Now().Format(time.RFC3339)},
	}

	result, err := r.transactionCollection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	transaction.TransactionID = result.InsertedID.(primitive.ObjectID)

	return transaction, nil
}
