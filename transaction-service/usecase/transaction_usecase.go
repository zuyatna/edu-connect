package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/zuyatna/edu-connect/transaction-service/model"
)

type ITransactionUsecase interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction) (*model.Transaction, error)
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
