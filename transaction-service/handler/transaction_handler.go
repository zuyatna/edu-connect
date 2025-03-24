package handler

import (
	"context"

	"github.com/zuyatna/edu-connect/transaction-service/middlewares"
	"github.com/zuyatna/edu-connect/transaction-service/model"
	pbFuncCollect "github.com/zuyatna/edu-connect/transaction-service/pb/fund_collect"
	pbTransaction "github.com/zuyatna/edu-connect/transaction-service/pb/transaction"
	pbUser "github.com/zuyatna/edu-connect/transaction-service/pb/user"
	"github.com/zuyatna/edu-connect/transaction-service/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type ITransactionHandler interface {
	CreateTransaction(ctx context.Context, req *pbTransaction.CreateTransactionRequest) (*pbTransaction.CreateTransactionResponse, error)
}

type TransactionServer struct {
	pbTransaction.UnimplementedTransactionServiceServer
	transactionUsecase usecase.ITransactionUsecase
	userClient         pbUser.UserServiceClient
	fundCollectClient  pbFuncCollect.FundCollectServiceClient
}

func NewTransactionHandler(
	transactionUsecase usecase.ITransactionUsecase,
	userClient pbUser.UserServiceClient,
	fundCollectClient pbFuncCollect.FundCollectServiceClient,
) *TransactionServer {
	return &TransactionServer{
		transactionUsecase: transactionUsecase,
		userClient:         userClient,
		fundCollectClient:  fundCollectClient,
	}
}

func (s *TransactionServer) CreateTransaction(ctx context.Context, req *pbTransaction.CreateTransactionRequest) (*pbTransaction.CreateTransactionResponse, error) {
	authenticatedUserID, ok := ctx.Value(middlewares.UserIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated user ID from context")
	}

	transaction_model := &model.Transaction{
		UserID:        authenticatedUserID,
		PostID:        req.PostId,
		PaymentID:     "00000000-0000-0000-0000-000000000000",
		Amount:        float64(req.Amount),
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
	}

	transaction, err := s.transactionUsecase.CreateTransaction(ctx, transaction_model)
	if err != nil {
		return nil, err
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get metadata from context")
	}

	outCtx := ctx
	if authHeaders, exists := md["authorization"]; exists && len(authHeaders) > 0 {
		outCtx = metadata.AppendToOutgoingContext(ctx, "authorization", authHeaders[0])
	}

	userResp, err := s.userClient.GetUserByID(outCtx, &pbUser.GetUserByIDRequest{
		Id: authenticatedUserID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	userName := userResp.Name
	if req.HideName {
		userName = "Anonymous"
	}

	_, err = s.fundCollectClient.CreateFundCollect(outCtx, &pbFuncCollect.CreateFundCollectRequest{
		PostId:        req.PostId,
		UserId:        authenticatedUserID,
		UserName:      userName,
		Amount:        req.Amount,
		TransactionId: transaction.TransactionID.String(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create fund collect: %v", err)
	}

	return &pbTransaction.CreateTransactionResponse{
		TransactionId: transaction.TransactionID.String(),
		PaymentId:     transaction.PaymentID,
		Amount:        float32(transaction.Amount),
		AccountNumber: transaction.AccountNumber,
		AccountName:   transaction.AccountName,
	}, nil
}
