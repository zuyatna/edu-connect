package routes

import (
	"net/http"

	pb "institution-service/pb/fund_collect"

	"github.com/labstack/echo/v4"
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
	groupFundCollect := e.Group("/v1/fund-collect")
	groupFundCollect.Use(AuthMiddleware)
	groupFundCollect.GET("/post/:id", h.GetFundCollectByPostID)
}

// GetFundCollectByPostID godoc
// @Summary      Get funding collection by Post ID.
// @Description  Get a specific funding collection by Post ID with authorization.
// @Tags         FundCollect
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header    string  true  "Bearer token"
// @Param        id            path      string    true  "Institution ID"
// @Success      200  {object}  model.FundCollectResponse "Success get funding collection data"
// @Failure      401  {object}  httputil.HTTPError "Unauthorized"
// @Failure      404  {object}  httputil.HTTPError "Funding collection not found"
// @Router       /v1/fund-collect/{id} [get]
func (h *FundCollectHTTPHandler) GetFundCollectByPostID(c echo.Context) error {
	req := new(pb.GetFundCollectByPostIDRequest)
	req.PostId = c.Param("id")

	res, err := h.fundCollectClient.GetFundCollectByPostID(c.Request().Context(), req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to retrieve fund collection data",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success get fund collection data",
		"data":    res,
	})
}
