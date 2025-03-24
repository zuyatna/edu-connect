package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	pbFuncCollect "github.com/zuyatna/edu-connect/transaction-service/pb/fund_collect"
	pbUser "github.com/zuyatna/edu-connect/transaction-service/pb/user"
	"github.com/zuyatna/edu-connect/transaction-service/usecase"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PaymentCallbackHandler struct {
	transactionUsecase usecase.ITransactionUsecase
	userClient         pbUser.UserServiceClient
	fundCollectClient  pbFuncCollect.FundCollectServiceClient
}

type XenditCallbackPayload struct {
	ID         string  `json:"id"`
	ExternalID string  `json:"external_id"`
	Status     string  `json:"status"`
	Amount     float64 `json:"amount"`
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

func (h *PaymentCallbackHandler) HandleCallback(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var payload XenditCallbackPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	objectID, err := primitive.ObjectIDFromHex(payload.ExternalID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid transaction ID format: %v", err), http.StatusBadRequest)
		return
	}

	transaction, err := h.transactionUsecase.GetTransactionByID(r.Context(), objectID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Transaction not found: %v", err), http.StatusNotFound)
		return
	}

	transaction.PaymentStatus = payload.Status

	if payload.Status == "PAID" {
		log.Printf("Updating fund collection for transaction %s", transaction.TransactionID)

		outCtx := r.Context()

		// userResp, err := h.userClient.GetUserByID(outCtx, &pbUser.GetUserByIDRequest{
		// 	Id: transaction.UserID,
		// })
		// if err != nil {
		// 	log.Printf("Failed to get user: %v", err)
		// }

		// userName := userResp.Name

		_, err = h.fundCollectClient.CreateFundCollect(outCtx, &pbFuncCollect.CreateFundCollectRequest{
			PostId:        transaction.PostID,
			UserId:        transaction.UserID,
			UserName:      "anonymous",
			Amount:        float32(transaction.Amount),
			TransactionId: transaction.TransactionID.Hex(),
		})
		if err != nil {
			log.Printf("Failed to create fund collect: %v", err)
		}
	} else if payload.Status == "EXPIRED" {
		log.Printf("Payment for transaction %s has expired", transaction.TransactionID)
	} else if payload.Status == "FAILED" {
		log.Printf("Payment for transaction %s has failed", transaction.TransactionID)
	}

	_, err = h.transactionUsecase.UpdateTransaction(r.Context(), transaction)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update transaction: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
