package handler

import (
	"context"

	"institution-service/model"
	pbFundCollect "institution-service/pb/fund_collect"
	pbPost "institution-service/pb/post"
	"institution-service/usecase"

	"github.com/google/uuid"
)

type IFundCollectHandler interface {
	CreateFundCollect(ctx context.Context, fundCollect *model.FundCollect) (*model.FundCollect, error)
	GetFundCollectByPostID(ctx context.Context, postID string) ([]model.FundCollect, error)
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

	fund_collect_model := &model.FundCollect{
		PostID:        postID,
		UserID:        req.UserId,
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
		UserId:        fund_collect.UserID,
		UserName:      fund_collect.UserName,
		Amount:        float32(fund_collect.Amount),
		TransactionId: fund_collect.TransactionID,
	}, nil
}

func (s *FundCollectServer) GetFundCollectByPostID(ctx context.Context, req *pbFundCollect.GetFundCollectByPostIDRequest) (*pbFundCollect.GetFundCollectByPostIDResponse, error) {
	fund_collects, err := s.fundCollectUsecase.GetFundCollectByPostID(ctx, req.PostId)
	if err != nil {
		return nil, err
	}

	var fund_collect_responses []*pbFundCollect.FundCollectResponse
	for _, fund_collect := range fund_collects {
		fund_collect_responses = append(fund_collect_responses, &pbFundCollect.FundCollectResponse{
			FundCollectId: fund_collect.FundCollectID.String(),
			PostId:        fund_collect.PostID.String(),
			UserId:        fund_collect.UserID,
			UserName:      fund_collect.UserName,
			Amount:        float32(fund_collect.Amount),
			TransactionId: fund_collect.TransactionID,
		})
	}

	return &pbFundCollect.GetFundCollectByPostIDResponse{
		Funds: fund_collect_responses,
	}, nil
}
