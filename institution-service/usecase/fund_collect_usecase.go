package usecase

import (
	"context"
	"errors"
	"strings"

	"institution-service/model"
	"institution-service/repository"
)

type IFundCollectUsecase interface {
	CreateFundCollect(ctx context.Context, fundCollect *model.FundCollect) (*model.FundCollect, error)
	GetFundCollectByPostID(ctx context.Context, postID string) ([]model.FundCollect, error)
}

type FundCollectUsecase struct {
	fundCollectRepository repository.IFundCollectRepository
}

func NewFundCollectUsecase(fundCollectRepository repository.IFundCollectRepository) *FundCollectUsecase {
	return &FundCollectUsecase{
		fundCollectRepository: fundCollectRepository,
	}
}

func (u *FundCollectUsecase) CreateFundCollect(ctx context.Context, fundCollect *model.FundCollect) (*model.FundCollect, error) {
	var e []string

	if fundCollect.PostID.String() == "00000000-0000-0000-0000-000000000000" {
		e = append(e, "Post ID is required")
	}
	if fundCollect.UserID == "00000000-0000-0000-0000-000000000000" {
		e = append(e, "User ID is required")
	}
	if fundCollect.UserName == "" {
		e = append(e, "User Name is required")
	}
	if fundCollect.Amount <= 0 {
		e = append(e, "Amount must be greater than 0")
	}

	if len(e) > 0 {
		return nil, errors.New(strings.Join(e, ", "))
	}

	return u.fundCollectRepository.CreateFundCollect(ctx, fundCollect)
}

func (u *FundCollectUsecase) GetFundCollectByPostID(ctx context.Context, postID string) ([]model.FundCollect, error) {
	return u.fundCollectRepository.GetFundCollectByPostID(ctx, postID)
}
