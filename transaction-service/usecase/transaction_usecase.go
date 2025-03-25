package usecase

import (
	"context"
	"errors"
	"strings"

	"transaction-service/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ITransactionUsecase interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
	GetTransactionByID(ctx context.Context, transactionID primitive.ObjectID) (*model.Transaction, error)
	UpdateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
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

func (u *TransactionUsecase) GetTransactionByID(ctx context.Context, transactionID primitive.ObjectID) (*model.Transaction, error) {
	return u.transactionRepository.GetTransactionByID(ctx, transactionID)
}

func (u *TransactionUsecase) UpdateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error) {
	return u.transactionRepository.UpdateTransaction(ctx, transaction)
}
