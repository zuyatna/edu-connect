package routes

import (
	"net/http"
	"strings"

	"transaction-service/httputil"
	pb "transaction-service/pb/transaction"
	"transaction-service/utils"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/metadata"
)

type TransactionHTTPHandler struct {
	transactionClient pb.TransactionServiceClient
}

func NewTransactionHTTPHandler(transactionClient pb.TransactionServiceClient) *TransactionHTTPHandler {
	return &TransactionHTTPHandler{
		transactionClient: transactionClient,
	}
}

func (h *TransactionHTTPHandler) Routes(e *echo.Echo) {
	e.POST("/transaction", h.authMiddleware2(h.CreateTransaction))
}

func (h *TransactionHTTPHandler) CreateTransaction(c echo.Context) error {
	req := new(pb.CreateTransactionRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, httputil.HTTPError{
			Message: "Invalid request body",
		})
	}

	res, err := h.transactionClient.CreateTransaction(c.Request().Context(), &pb.CreateTransactionRequest{
		PostId:        req.PostId,
		Amount:        req.Amount,
		AccountNumber: req.AccountNumber,
		AccountName:   req.AccountName,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *TransactionHTTPHandler) authMiddleware2(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, httputil.HTTPError{
				Message: "Unauthorized",
			})
		}

		tokenParts := strings.Split(token, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			return echo.NewHTTPError(http.StatusUnauthorized, httputil.HTTPError{
				Message: "Invalid token format",
			})
		}

		userID, email, err := utils.ValidateJWT(tokenParts[1])
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, httputil.HTTPError{
				Message: "Invalid token: " + err.Error(),
			})
		}

		metadataMap := map[string]string{
			"authorization": token,
		}

		if userID != "" {
			metadataMap["id"] = userID
		}

		if email != "" {
			metadataMap["email"] = email
		}

		md := metadata.New(metadataMap)
		ctx := metadata.NewOutgoingContext(c.Request().Context(), md)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

// func (h *TransactionHTTPHandler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		token := c.Request().Header.Get("Authorization")
// 		if token == "" {
// 			return echo.NewHTTPError(http.StatusUnauthorized, httputil.HTTPError{
// 				Message: "Unauthorized",
// 			})
// 		}

// 		tokenParts := strings.Split(token, " ")
// 		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
// 			return echo.NewHTTPError(http.StatusUnauthorized, httputil.HTTPError{
// 				Message: "Invalid token format",
// 			})
// 		}

// 		claims, err := utils.ValidateToken(tokenParts[1])
// 		if err != nil {
// 			return echo.NewHTTPError(http.StatusUnauthorized, httputil.HTTPError{
// 				Message: "Invalid token: " + err.Error(),
// 			})
// 		}

// 		userID := (*claims)["user_id"].(string)

// 		md := metadata.New(map[string]string{
// 			"authorization": token,
// 			"user_id":       userID,
// 		})
// 		ctx := metadata.NewOutgoingContext(c.Request().Context(), md)
// 		c.SetRequest(c.Request().WithContext(ctx))

// 		return next(c)
// 	}
// }
