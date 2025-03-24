package routes

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zuyatna/edu-connect/institution-service/httputil"
	pb "github.com/zuyatna/edu-connect/institution-service/pb/institution"
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
	e.POST("/institution/register", h.RegisterInstitution)
	e.POST("/institution/login", h.LoginInstitution)

	e.GET("/institution/:id", AuthMiddleware(h.GetInstitutionByID))
	e.PUT("/institution/:id", AuthMiddleware(h.UpdateInstitution))
	e.DELETE("/institution/:id", AuthMiddleware(h.DeleteInstitution))
}

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

	return c.JSON(http.StatusCreated, res)
}

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

	return c.JSON(http.StatusOK, res)
}

func (h *InstitutionHTTPHandler) GetInstitutionByID(c echo.Context) error {
	req := new(pb.GetInstitutionByIDRequest)
	req.InstitutionId = c.Param("id")

	res, err := h.institutionClient.GetInstitutionByID(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

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

	return c.JSON(http.StatusOK, res)
}

func (h *InstitutionHTTPHandler) DeleteInstitution(c echo.Context) error {
	req := new(pb.DeleteInstitutionRequest)
	req.InstitutionId = c.Param("id")

	res, err := h.institutionClient.DeleteInstitution(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Request().URL.Path
		if path == "/institution/register" || path == "/institution/login" {
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
