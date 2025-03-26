package handler

import (
	"context"
	"fmt"
	"log"

	"transaction-service/client"
	"transaction-service/middlewares"
	"transaction-service/model"
	pbFuncCollect "transaction-service/pb/fund_collect"
	pbTransaction "transaction-service/pb/transaction"
	pbUser "transaction-service/pb/user"
	"transaction-service/usecase"

	"github.com/google/uuid"
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
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get metadata from context")
	}

	fmt.Printf("Context keys: %v\n", ctx)
	if md["authorization"] != nil {
		fmt.Printf("Auth header present\n")
	}

	email, emailOk := ctx.Value(middlewares.EmailKey).(string)
	if !emailOk || email == "" {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated user email from context")
	}

	getUser, err := s.transactionUsecase.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("Failed to get user by email: %v", err)
	}

	authenticatedUserID := getUser.ID

	if userID, userIDOk := ctx.Value(middlewares.UserIDKey).(string); userIDOk && userID != "" {
		authenticatedUserID = "00000000-0000-0000-0000-000000000000" // userID
	}

	transaction_model := &model.Transaction{
		UserID:        authenticatedUserID,
		PostID:        req.PostId,
		UserEmail:     email,
		PaymentID:     "pending",
		Amount:        float64(req.Amount),
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
	}

	transaction, err := s.transactionUsecase.CreateTransaction(ctx, transaction_model)
	if err != nil {
		return nil, err
	}

	transactionIDStr := transaction.TransactionID.Hex()
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post ID format: %v", err)
	}

	post, err := s.transactionUsecase.GetPostByID(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get post: %v", err)
	}

	invoiceReq := client.CreateInvoiceRequest{
		ExternalID:         transactionIDStr,
		Amount:             transaction.Amount,
		PayerEmail:         email,
		Description:        fmt.Sprintf("Fund contribution for %s", post.Title),
		CustomerName:       "anonymous",
		InvoiceDuration:    86400, // 24 hours
		SuccessRedirectURL: fmt.Sprintf("https://transaction-service-1011483964797.asia-southeast2.run.app/payment/success?external_id=%s", transactionIDStr),
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
