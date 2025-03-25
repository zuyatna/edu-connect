package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"transaction-service/middlewares"
	pbFuncCollect "transaction-service/pb/fund_collect"
	pbUser "transaction-service/pb/user"
	"transaction-service/usecase"

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

	objectID, err := primitive.ObjectIDFromHex(payload.TransactionID)
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
	transaction.PaymentID = payload.PaymentID

	if payload.Status == "PAID" {
		log.Printf("Updating fund collection for transaction %s", transaction.TransactionID)

		outCtx := r.Context()

		token := r.Header.Get("Authorization")
		if token == "" {
			token = "Bearer your-service-token-here"
		}

		md := metadata.Pairs("authorization", token)
		authCtx := metadata.NewOutgoingContext(outCtx, md)

		// This is commented out because the user service is not yet implemented
		// userResp, err := h.userClient.GetUserByID(authCtx, &pbUser.GetUserByIDRequest{
		// 	Id: transaction.UserID,
		// })
		// if err != nil {
		// 	log.Printf("Failed to get user: %v", err)
		// }
		// userName := userResp.Name

		var userName string
		if email := authCtx.Value(middlewares.EmailKey); email != nil {
			userName = email.(string)
		} else {
			userName = "Anonymous User"
			log.Printf("Warning: User email not found in context for transaction %s", transaction.TransactionID)
		}

		_, err = h.fundCollectClient.CreateFundCollect(authCtx, &pbFuncCollect.CreateFundCollectRequest{
			PostId:        transaction.PostID,
			UserId:        "00000000-0000-0000-0000-000000000000",
			UserName:      userName,
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
