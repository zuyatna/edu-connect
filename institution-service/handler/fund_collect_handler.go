package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/zuyatna/edu-connect/institution-service/model"
	pbFundCollect "github.com/zuyatna/edu-connect/institution-service/pb/fund_collect"
	pbPost "github.com/zuyatna/edu-connect/institution-service/pb/post"
	"github.com/zuyatna/edu-connect/institution-service/usecase"
)

type IFundCollectHandler interface {
	CreateFundCollect(ctx context.Context, fund_collect *model.FundCollect) (*model.FundCollect, error)
}

type FundCollectServer struct {
	pbFundCollect.UnimplementedFundCollectServiceServer
	pbPost.UnimplementedPostServiceServer
	fundCollectUsecase usecase.IFundCollectUsecase
	postUsecase        usecase.IPostUsecase
}

func NewFundCollectHandler(fundCollectUsecase usecase.IFundCollectUsecase, postUsecase usecase.IPostUsecase) *FundCollectServer {
	return &FundCollectServer{
		fundCollectUsecase: fundCollectUsecase,
		postUsecase:        postUsecase,
	}
}
func (s *FundCollectServer) CreateFundCollect(ctx context.Context, req *pbFundCollect.CreateFundCollectRequest) (*pbFundCollect.CreateFundCollectResponse, error) {
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	fund_collect_model := &model.FundCollect{
		PostID:        postID,
		UserID:        userID,
		UserName:      req.UserName,
		Amount:        float64(req.Amount),
		TransactionID: req.TransactionId,
	}

	fund_collect, err := s.fundCollectUsecase.CreateFundCollect(ctx, fund_collect_model)
	if err != nil {
		return nil, err
	}

	_, err = s.postUsecase.AddPostFundAchieved(ctx, fund_collect.PostID, fund_collect.Amount)
	if err != nil {
		return nil, err
	}

	return &pbFundCollect.CreateFundCollectResponse{
		FundCollectId: fund_collect.FundCollectID.String(),
		PostId:        fund_collect.PostID.String(),
		UserId:        fund_collect.UserID.String(),
		UserName:      fund_collect.UserName,
		Amount:        float32(fund_collect.Amount),
		TransactionId: fund_collect.TransactionID,
	}, nil
}
