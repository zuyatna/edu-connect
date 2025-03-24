package repository

import (
	"context"

	"github.com/zuyatna/edu-connect/institution-service/model"
	"gorm.io/gorm"
)

type IFundCollectRepository interface {
	CreateFundCollect(ctx context.Context, fund_collect *model.FundCollect) (*model.FundCollect, error)
}

type FundCollectRepository struct {
	db *gorm.DB
}

func NewFundCollectRepository(db *gorm.DB) *FundCollectRepository {
	return &FundCollectRepository{
		db: db,
	}
}

func (r *FundCollectRepository) CreateFundCollect(ctx context.Context, fund_collect *model.FundCollect) (*model.FundCollect, error) {
	if err := r.db.Create(fund_collect).Error; err != nil {
		return nil, err
	}

	return fund_collect, nil
}
