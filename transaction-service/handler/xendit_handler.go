package handler

import (
	"fmt"
	"log"
	"net/http"

	"transaction-service/model"
	pbFuncCollect "transaction-service/pb/fund_collect"
	pbUser "transaction-service/pb/user"
	"transaction-service/usecase"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/metadata"
)

type PaymentCallbackHandler struct {
	transactionUsecase usecase.ITransactionUsecase
	userClient         pbUser.UserServiceClient
	fundCollectClient  pbFuncCollect.FundCollectServiceClient
}

type XenditCallbackPayload struct {
	PaymentID     string  `json:"payment_id"`
	TransactionID string  `json:"transaction_id"`
	Status        string  `json:"status"`
	Amount        float64 `json:"amount"`
}

func NewPaymentCallbackHandler(
	transactionUsecase usecase.ITransactionUsecase,
	userClient pbUser.UserServiceClient,
	fundCollectClient pbFuncCollect.FundCollectServiceClient,
) *PaymentCallbackHandler {
	return &PaymentCallbackHandler{
		transactionUsecase: transactionUsecase,
		userClient:         userClient,
		fundCollectClient:  fundCollectClient,
	}
}

func (h *PaymentCallbackHandler) HandleSuccessRedirect(w http.ResponseWriter, r *http.Request) {
	transactionID := r.URL.Query().Get("external_id")
	if transactionID == "" {
		http.Error(w, "Missing transaction ID", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(transactionID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid transaction ID format: %v", err), http.StatusBadRequest)
		return
	}

	transaction, err := h.transactionUsecase.GetTransactionByID(r.Context(), objectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Transaction not found: %v", err), http.StatusNotFound)
		return
	}

	transaction.PaymentStatus = "PAID"

	log.Printf("Processing successful payment for transaction %s", transaction.TransactionID)

	outCtx := r.Context()

	token := r.Header.Get("Authorization")
	if token == "" {
		token = "Bearer your-service-token-here"
	}

	md := metadata.Pairs("authorization", token)
	authCtx := metadata.NewOutgoingContext(outCtx, md)

	var userName string
	if email := transaction.UserEmail; email != "" {
		userName = email
	} else {
		userName = "Anonymous User"
		log.Printf("Warning: User email not found for transaction %s", transaction.TransactionID)
	}

	postUUID, err := uuid.Parse(transaction.PostID)
	if err != nil {
		log.Printf("Failed to parse PostID as UUID: %v", err)
		http.Error(w, fmt.Sprintf("Invalid PostID format: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = h.transactionUsecase.CreateFundCollect(authCtx, &model.FundCollect{
		PostID:        postUUID,
		UserID:        transaction.UserID,
		UserName:      userName,
		Amount:        float64(transaction.Amount),
		TransactionID: transaction.TransactionID.Hex(),
	})
	if err != nil {
		log.Printf("Failed to create fund collect: %v", err)
	}

	_, err = h.transactionUsecase.AddPostFundAchieved(r.Context(), postUUID, transaction.Amount)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update post fund achieved: %v", err), http.StatusInternalServerError)
		return
	}

	_, err = h.transactionUsecase.UpdateTransaction(r.Context(), transaction)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "<html><body><h1>Payment Successful</h1><p>Thank you for your contribution!</p></body></html>")
}
