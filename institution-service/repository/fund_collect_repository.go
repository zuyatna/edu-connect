package repository

import (
	"context"

	"institution-service/model"
	"gorm.io/gorm"
)

type IFundCollectRepository interface {
	CreateFundCollect(ctx context.Context, fundCollect *model.FundCollect) (*model.FundCollect, error)
	GetFundCollectByPostID(ctx context.Context, postID string) ([]model.FundCollect, error)
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

func (r *FundCollectRepository) GetFundCollectByPostID(ctx context.Context, post_id string) ([]model.FundCollect, error) {
	var fund_collects []model.FundCollect

	if err := r.db.Where("post_id = ?", post_id).Find(&fund_collects).Error; err != nil {
		return nil, err
	}

	return fund_collects, nil
}
