package handler

import (
	"context"
	"fmt"

	"github.com/zuyatna/edu-connect/transaction-service/client"
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
	xenditClient       *client.XenditClient
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
		xenditClient:       client.NewXenditClient(),
	}
}

func (s *TransactionServer) CreateTransaction(ctx context.Context, req *pbTransaction.CreateTransactionRequest) (*pbTransaction.CreateTransactionResponse, error) {
	authenticatedUserID, ok := ctx.Value(middlewares.UserIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated user ID from context")
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

	transaction_model := &model.Transaction{
		UserID:        authenticatedUserID,
		PostID:        req.PostId,
		PaymentID:     "pending",
		Amount:        float64(req.Amount),
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
	}

	transaction, err := s.transactionUsecase.CreateTransaction(ctx, transaction_model)
	if err != nil {
		return nil, err
	}

	invoiceReq := client.CreateInvoiceRequest{
		ExternalID:         transaction.TransactionID.String(),
		Amount:             transaction.Amount,
		PayerEmail:         userResp.Email,
		Description:        fmt.Sprintf("Fund contribution for post %s", req.PostId),
		CustomerName:       userResp.Name,
		InvoiceDuration:    86400, // 24 hours
		SuccessRedirectURL: "https://edu-connect.example.com/payment/success",
		FailureRedirectURL: "https://edu-connect.example.com/payment/failed",
		CallbackURL:        "https://edu-connect.example.com/api/payment/callback",
	}

	invoice, err := s.xenditClient.CreateInvoice(invoiceReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create payment invoice: %v", err)
	}

	transaction.PaymentID = invoice.ID
	transaction.PaymentURL = invoice.InvoiceURL
	transaction.PaymentStatus = "PENDING"

	_, err = s.transactionUsecase.UpdateTransaction(ctx, transaction)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update transaction with payment details: %v", err)
	}

	return &pbTransaction.CreateTransactionResponse{
		TransactionId: transaction.TransactionID.String(),
		PaymentId:     transaction.PaymentID,
		Amount:        float32(transaction.Amount),
		AccountNumber: transaction.AccountNumber,
		AccountName:   transaction.AccountName,
		PaymentUrl:    invoice.InvoiceURL,
		Status:        "PENDING",
	}, nil
}
