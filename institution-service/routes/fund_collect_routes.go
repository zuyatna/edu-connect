package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	pb "github.com/zuyatna/edu-connect/institution-service/pb/fund_collect"
)

type FundCollectHTTPHandler struct {
	fundCollectClient pb.FundCollectServiceClient
}

func NewFundCollectHTTPHandler(fundCollectClient pb.FundCollectServiceClient) *FundCollectHTTPHandler {
	return &FundCollectHTTPHandler{
		fundCollectClient: fundCollectClient,
	}
}

func (h *FundCollectHTTPHandler) Routes(e *echo.Echo) {
	groupFundCollect := e.Group("/fund-collect")
	groupFundCollect.Use(AuthMiddleware)
	groupFundCollect.GET("/post/:id", h.GetFundCollectByPostID)
}

func (h *FundCollectHTTPHandler) GetFundCollectByPostID(c echo.Context) error {
	req := new(pb.GetFundCollectByPostIDRequest)
	req.PostId = c.Param("id")

	res, err := h.fundCollectClient.GetFundCollectByPostID(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve fund collection data",
		})
	}

	return c.JSON(http.StatusOK, res)
}
