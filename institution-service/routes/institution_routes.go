package routes

import (
	"context"
	"net/http"

	"institution-service/httputil"
	pb "institution-service/pb/institution"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/metadata"
)

type InstitutionHTTPHandler struct {
	institutionClient pb.InstitutionServiceClient
}

func NewInstitutionHTTPHandler(institutionClient pb.InstitutionServiceClient) *InstitutionHTTPHandler {
	return &InstitutionHTTPHandler{
		institutionClient: institutionClient,
	}
}

func (h *InstitutionHTTPHandler) Routes(e *echo.Echo) {
	e.POST("/v1/institution/register", h.RegisterInstitution)
	e.POST("/v1/institution/login", h.LoginInstitution)

	e.GET("/v1/institution", AuthMiddleware(h.GetInstitutionByID))
	e.PUT("/v1/institution/:id", AuthMiddleware(h.UpdateInstitution))
	e.DELETE("/v1/institution/:id", AuthMiddleware(h.DeleteInstitution))
}

// RegisterInstitution godoc
// @Summary      Register a new Institution.
// @Description  Register institution with name, email, etc. Email must be unique and password will be hashed before saved to database.
// @Tags         Institution
// @Accept       json
// @Produce      json
// @Param        request body model.InstitutionRequest true "Institution created details"
// @Success      200 {object} model.InstitutionResponse "Institution created successfully"
// @Failure      500 {object} httputil.HTTPError "Internal server error"
// @Router       /v1/institution/register [post]
func (h *InstitutionHTTPHandler) RegisterInstitution(c echo.Context) error {
	req := new(pb.RegisterInstitutionRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, httputil.HTTPError{
			Message: "Invalid request body",
		})
	}

	cleanCtx := context.Background()

	res, err := h.institutionClient.RegisterInstitution(cleanCtx, &pb.RegisterInstitutionRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Address:  req.Address,
		Phone:    req.Phone,
		Website:  req.Website,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Institution created successfully",
		"data":    res,
	})
}

// LoginInstitution godoc
// @Summary      Login Institution.
// @Description  Login Institution with email and password.
// @Tags         Institution
// @Accept       json
// @Produce      json
// @Param        request body model.InstitutionLoginRequest true "Institution login"
// @Success      200 {object} model.InstitutionToken "Institution login successfully"
// @Failure      500 {object} httputil.HTTPError "Internal server error"
// @Router       /v1/institution/login [post]
func (h *InstitutionHTTPHandler) LoginInstitution(c echo.Context) error {
	req := new(pb.LoginInstitutionRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, httputil.HTTPError{
			Message: "Invalid request body",
		})
	}

	cleanCtx := context.Background()

	res, err := h.institutionClient.LoginInstitution(cleanCtx, &pb.LoginInstitutionRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Institution login successfully",
		"data":    res,
	})
}

// GetInstitutionByID godoc
// @Summary      Get Institution.
// @Description  Get data institution with authorization.
// @Tags         Institution
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header    string  true  "Bearer token"
// @Success      200  {object}  model.InstitutionResponse "Success get institution data"
// @Failure      401  {object}  httputil.HTTPError "Unauthorized"
// @Failure      404  {object}  httputil.HTTPError "Data not found"
// @Router       /v1/institution [get]
func (h *InstitutionHTTPHandler) GetInstitutionByID(c echo.Context) error {
	req := new(pb.GetInstitutionByIDRequest)
	req.InstitutionId = c.Param("id")

	res, err := h.institutionClient.GetInstitutionByID(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success get institution data",
		"data":    res,
	})
}

// UpdateInstitution godoc
// @Summary      Update Institution.
// @Description  Update an existing institution with authorization.
// @Tags         Institution
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header    string  true  "Bearer token"
// @Param        id            path      string    true  "Institution ID"
// @Param        user       body      model.InstitutionRequest  true  "Updated institution data"
// @Success      200  {object}  model.InstitutionResponse "Success update institution data"
// @Failure      401  {object}  httputil.HTTPError "Unauthorized"
// @Failure      404  {object}  httputil.HTTPError "User not found"
// @Router       /v1/institution/{id} [put]
func (h *InstitutionHTTPHandler) UpdateInstitution(c echo.Context) error {
	req := new(pb.UpdateInstitutionRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, httputil.HTTPError{
			Message: "Invalid request body",
		})
	}

	req.InstitutionId = c.Param("id")
	res, err := h.institutionClient.UpdateInstitution(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success update institution data",
		"data":    res,
	})
}

// DeleteInstitution godoc
// @Summary      Delete Institution.
// @Description  Delete an existing institution with authorization.
// @Tags         Institution
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header    string  true  "Bearer token"
// @Param        id            path      string     true  "Institution ID"
// @Success      200  {object}  model.InstitutionDeleteResponse "Success delete institution data"
// @Failure      500  {object}  httputil.HTTPError "Internal server error"
// @Router       /v1/institution/{id} [delete]
func (h *InstitutionHTTPHandler) DeleteInstitution(c echo.Context) error {
	req := new(pb.DeleteInstitutionRequest)
	req.InstitutionId = c.Param("id")

	_, err := h.institutionClient.DeleteInstitution(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success delete institution data",
		"data":    map[string]interface{}{},
	})
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path
		if path == "/v1/institution/register" ||
			path == "/v1/institution/login" ||
			path == "/v1/posts" {
			return next(c)
		}

		token := c.Request().Header.Get("Authorization")
		if token == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, httputil.HTTPError{
				Message: "Unauthorized",
			})
		}

		md := metadata.New(map[string]string{"authorization": token})
		ctx := metadata.NewOutgoingContext(c.Request().Context(), md)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}
