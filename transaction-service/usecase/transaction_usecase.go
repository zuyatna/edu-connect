package usecase

import (
	"context"
	"errors"
	"strings"

	"transaction-service/model"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ITransactionUsecase interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
	CreateFundCollect(ctx context.Context, fundCollect *model.FundCollect) (*model.FundCollect, error)
	GetTransactionByID(ctx context.Context, transactionID primitive.ObjectID) (*model.Transaction, error)
	GetPostByID(ctx context.Context, postID uuid.UUID) (*model.Post, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
	AddPostFundAchieved(ctx context.Context, postID uuid.UUID, amount float64) (*model.Post, error)
}

type TransactionUsecase struct {
	transactionRepository ITransactionUsecase
}

func NewTransactionUsecase(transactionRepository ITransactionUsecase) *TransactionUsecase {
	return &TransactionUsecase{
		transactionRepository: transactionRepository,
	}
}

func (u *TransactionUsecase) CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	var e []string

	if transaction.PostID == "00000000-0000-0000-0000-000000000000" {
		e = append(e, "Post ID is required")
	}
	if transaction.Amount <= 0 {
		e = append(e, "Amount must be greater than 0")
	}
	if transaction.AccountNumber == "" {
		e = append(e, "Account Number is required")
	}
	if transaction.AccountName == "" {
		e = append(e, "Account Name is required")
	}

	if len(e) > 0 {
		return nil, errors.New(strings.Join(e, ", "))
	}

	return u.transactionRepository.CreateTransaction(ctx, transaction)
}

func (u *TransactionUsecase) CreateFundCollect(ctx context.Context, fundCollect *model.FundCollect) (*model.FundCollect, error) {
	return u.transactionRepository.CreateFundCollect(ctx, fundCollect)
}

func (u *TransactionUsecase) GetTransactionByID(ctx context.Context, transactionID primitive.ObjectID) (*model.Transaction, error) {
	return u.transactionRepository.GetTransactionByID(ctx, transactionID)
}

func (u *TransactionUsecase) GetPostByID(ctx context.Context, postID uuid.UUID) (*model.Post, error) {
	return u.transactionRepository.GetPostByID(ctx, postID)
}

func (u *TransactionUsecase) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return u.transactionRepository.GetUserByEmail(ctx, email)
}

func (u *TransactionUsecase) UpdateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	return u.transactionRepository.UpdateTransaction(ctx, transaction)
}

func (u *TransactionUsecase) AddPostFundAchieved(ctx context.Context, postID uuid.UUID, amount float64) (*model.Post, error) {
	return u.transactionRepository.AddPostFundAchieved(ctx, postID, amount)
}
