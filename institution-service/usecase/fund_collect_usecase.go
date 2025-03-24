package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/zuyatna/edu-connect/institution-service/model"
)

type IFundCollectUsecase interface {
	CreateFundCollect(ctx context.Context, fund_collect *model.FundCollect) (*model.FundCollect, error)
}

type FundCollectUsecase struct {
	fundCollectRepository IFundCollectUsecase
}

func NewFundCollectUsecase(fundCollectRepository IFundCollectUsecase) *FundCollectUsecase {
	return &FundCollectUsecase{
		fundCollectRepository: fundCollectRepository,
	}
}

func (u *FundCollectUsecase) CreateFundCollect(ctx context.Context, fund_collect *model.FundCollect) (*model.FundCollect, error) {
	var e []string

	if fund_collect.PostID.String() == "00000000-0000-0000-0000-000000000000" {
		e = append(e, "Post ID is required")
	}
	if fund_collect.UserID.String() == "00000000-0000-0000-0000-000000000000" {
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
