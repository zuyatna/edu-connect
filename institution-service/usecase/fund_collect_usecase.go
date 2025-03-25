package usecase

import (
	"context"
	"errors"
	"strings"

	"institution-service/model"
	"institution-service/repository"
)

type IFundCollectUsecase interface {
	CreateFundCollect(ctx context.Context, fund_collect *model.FundCollect) (*model.FundCollect, error)
	GetFundCollectByPostID(ctx context.Context, post_id string) ([]model.FundCollect, error)
}

type FundCollectUsecase struct {
	fundCollectRepository repository.IFundCollectRepository
}

func NewFundCollectUsecase(fundCollectRepository repository.IFundCollectRepository) *FundCollectUsecase {
	return &FundCollectUsecase{
		fundCollectRepository: fundCollectRepository,
	}
}

func (u *FundCollectUsecase) CreateFundCollect(ctx context.Context, fund_collect *model.FundCollect) (*model.FundCollect, error) {
	var e []string

	if fund_collect.PostID.String() == "00000000-0000-0000-0000-000000000000" {
		e = append(e, "Post ID is required")
	}
	if fund_collect.UserID == "00000000-0000-0000-0000-000000000000" {
		e = append(e, "User ID is required")
	}
	if fund_collect.UserName == "" {
		e = append(e, "User Name is required")
	}
	if fund_collect.Amount <= 0 {
		e = append(e, "Amount must be greater than 0")
	}

	if len(e) > 0 {
		return nil, errors.New(strings.Join(e, ", "))
	}

	return u.fundCollectRepository.CreateFundCollect(ctx, fund_collect)
}

func (u *FundCollectUsecase) GetFundCollectByPostID(ctx context.Context, post_id string) ([]model.FundCollect, error) {
	return u.fundCollectRepository.GetFundCollectByPostID(ctx, post_id)
}
