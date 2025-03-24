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
	GetTransactionByID(ctx context.Context, transactionID primitive.ObjectID) (*model.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
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

func (r *TransactionRepository) GetTransactionByID(ctx context.Context, transactionID primitive.ObjectID) (*model.Transaction, error) {
	var transaction model.Transaction

	filter := bson.D{
		{Key: "_id", Value: transactionID},
	}

	err := r.transactionCollection.FindOne(ctx, filter).Decode(&transaction)
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

func (r *TransactionRepository) UpdateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	filter := bson.D{
		{Key: "_id", Value: transaction.TransactionID},
	}

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "payment_id", Value: transaction.PaymentID},
			{Key: "payment_url", Value: transaction.PaymentURL},
			{Key: "payment_status", Value: transaction.PaymentStatus},
			{Key: "updated_at", Value: time.Now().Format(time.RFC3339)},
		}},
	}

	_, err := r.transactionCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}
