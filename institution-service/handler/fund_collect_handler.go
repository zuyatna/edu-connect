package handler

import (
	"context"

	"github.com/google/uuid"
	"github.com/zuyatna/edu-connect/institution-service/model"
	pb "github.com/zuyatna/edu-connect/institution-service/pb/fund_collect"
	"github.com/zuyatna/edu-connect/institution-service/usecase"
)

type IFundCollectHandler interface {
	CreateFundCollect(ctx context.Context, fund_collect *model.FundCollect) (*model.FundCollect, error)
}

type FundCollectServer struct {
	pb.UnimplementedFundCollectServiceServer
	fundCollectUsecase usecase.IFundCollectUsecase
}

func NewFundCollectHandler(fundCollectUsecase usecase.IFundCollectUsecase) *FundCollectServer {
	return &FundCollectServer{
		fundCollectUsecase: fundCollectUsecase,
	}
}
func (s *FundCollectServer) CreateFundCollect(ctx context.Context, req *pb.CreateFundCollectRequest) (*pb.CreateFundCollectResponse, error) {
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

	return &pb.CreateFundCollectResponse{
		FundCollectId: fund_collect.FundCollectID.String(),
		PostId:        fund_collect.PostID.String(),
		UserId:        fund_collect.UserID.String(),
		UserName:      fund_collect.UserName,
		Amount:        float32(fund_collect.Amount),
		TransactionId: fund_collect.TransactionID,
	}, nil
}
