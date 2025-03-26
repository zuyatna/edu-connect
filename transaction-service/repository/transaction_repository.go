package repository

import (
	"context"
	"fmt"
	"os"
	"time"

	"transaction-service/model"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ITransactionRepository interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
	CreateFundCollect(ctx context.Context, fundCollect *model.FundCollect) (*model.FundCollect, error)
	GetTransactionByID(ctx context.Context, transactionID primitive.ObjectID) (*model.Transaction, error)
	GetPostByID(ctx context.Context, postID uuid.UUID) (*model.Post, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
	AddPostFundAchieved(ctx context.Context, postID uuid.UUID, amount float64) (*model.Post, error)
}

type TransactionRepository struct {
	transactionCollection *mongo.Collection
	gormClient            *gorm.DB
}

func NewTransactionRepository(
	mongos *mongo.Database,
	gormClient *gorm.DB,
) *TransactionRepository {
	return &TransactionRepository{
		transactionCollection: mongos.Collection("transactions"),
		gormClient:            gormClient,
	}
}

func (r *TransactionRepository) CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	doc := bson.D{
		{Key: "user_id", Value: transaction.UserID},
		{Key: "post_id", Value: transaction.PostID},
		{Key: "user_email", Value: transaction.UserEmail},
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

func (r *TransactionRepository) CreateFundCollect(ctx context.Context, fundCollect *model.FundCollect) (*model.FundCollect, error) {
	if err := r.gormClient.Create(fundCollect).Error; err != nil {
		return nil, err
	}

	return fundCollect, nil
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

func (r *TransactionRepository) GetPostByID(ctx context.Context, PostID uuid.UUID) (*model.Post, error) {
	var post model.Post
	if err := r.gormClient.Where("post_id = ?", PostID).First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *TransactionRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	dsn := os.Getenv("POSTGRES_URI_EXTERNAL")
	userDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user database: %w", err)
	}

	sqlDB, err := userDB.DB()
	if err != nil {
		return nil, err
	}
	defer sqlDB.Close()

	if err := userDB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
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

func (r *TransactionRepository) AddPostFundAchieved(ctx context.Context, postID uuid.UUID, amount float64) (*model.Post, error) {
	var post model.Post
	if err := r.gormClient.Where("post_id = ?", postID).First(&post).Error; err != nil {
		return nil, err
	}

	post.FundAchieved += amount

	if err := r.gormClient.Save(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}
