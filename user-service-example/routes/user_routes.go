package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"user-service-example/httputil"
	pb "user-service-example/pb/user"
	"google.golang.org/grpc/metadata"
)

type UserHTTPHandler struct {
	userClient pb.UserServiceClient
}

func NewUserHTTPHandler(userClient pb.UserServiceClient) *UserHTTPHandler {
	return &UserHTTPHandler{
		userClient: userClient,
	}
}

func (h *UserHTTPHandler) Routes(e *echo.Echo) {
	e.POST("/user/register", h.RegisterUser)
	e.POST("/user/login", h.LoginUser)

	e.GET("/user/:id", h.authMiddleware(h.GetUserByID))
	e.PUT("/user/:id", h.authMiddleware(h.UpdateUser))
	e.PUT("/user/:id/order", h.authMiddleware(h.UpdateDonateCountUser))
	e.DELETE("/user/:id", h.authMiddleware(h.DeleteUser))
}

func (h *UserHTTPHandler) RegisterUser(c echo.Context) error {
	req := new(pb.RegisterUserRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, httputil.HTTPError{
			Message: "Invalid request body",
		})
	}

	res, err := h.userClient.RegisterUser(c.Request().Context(), &pb.RegisterUserRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *UserHTTPHandler) LoginUser(c echo.Context) error {
	req := new(pb.LoginUserRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, httputil.HTTPError{
			Message: "Invalid request body",
		})
	}

	res, err := h.userClient.LoginUser(c.Request().Context(), &pb.LoginUserRequest{
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

func (h *UserHTTPHandler) GetUserByID(c echo.Context) error {
	req := new(pb.GetUserByIDRequest)
	req.Id = c.Param("id")

	res, err := h.userClient.GetUserByID(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *UserHTTPHandler) UpdateUser(c echo.Context) error {
	req := new(pb.UpdateUserRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, httputil.HTTPError{
			Message: "Invalid request body",
		})
	}

	req.Id = c.Param("id")
	res, err := h.userClient.UpdateUser(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *UserHTTPHandler) UpdateDonateCountUser(c echo.Context) error {
	req := new(pb.UpdateDonateCountRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, httputil.HTTPError{
			Message: "Invalid request body",
		})
	}

	req.Id = c.Param("id")
	res, err := h.userClient.UpdateDonateCountUser(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *UserHTTPHandler) DeleteUser(c echo.Context) error {
	req := new(pb.DeleteUserRequest)
	req.Id = c.Param("id")

	resp, err := h.userClient.DeleteUser(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *UserHTTPHandler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
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
