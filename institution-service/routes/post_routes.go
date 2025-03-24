package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/zuyatna/edu-connect/institution-service/httputil"
	pb "github.com/zuyatna/edu-connect/institution-service/pb/post"
)

type PostHTTPHandler struct {
	postClient pb.PostServiceClient
}

func NewPostHTTPHandler(postClient pb.PostServiceClient) *PostHTTPHandler {
	return &PostHTTPHandler{
		postClient: postClient,
	}
}

func (h *PostHTTPHandler) Routes(e *echo.Echo) {
	groupPost := e.Group("/post")
	e.Use(AuthMiddleware)
	groupPost.POST("", h.CreatePost)
	groupPost.GET("/:id", h.GetPostByID)
	groupPost.GET("/institution/:id", h.GetAllPostByInstitutionID)
	groupPost.PUT("/:id", h.UpdatePost)
	groupPost.DELETE("/:id", h.DeletePost)
}

func (h *PostHTTPHandler) CreatePost(c echo.Context) error {
	req := new(pb.CreatePostRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, httputil.HTTPError{
			Message: "Invalid request body",
		})
	}

	res, err := h.postClient.CreatePost(c.Request().Context(), &pb.CreatePostRequest{
		Title:      req.Title,
		Body:       req.Body,
		DateStart:  req.DateStart,
		DateEnd:    req.DateEnd,
		FundTarget: req.FundTarget,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, res)
}

func (h *PostHTTPHandler) GetPostByID(c echo.Context) error {
	req := new(pb.GetPostByIDRequest)
	req.PostId = c.Param("id")

	res, err := h.postClient.GetPostByID(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *PostHTTPHandler) GetAllPostByInstitutionID(c echo.Context) error {
	req := new(pb.GetAllPostByInstitutionIDRequest)
	req.InstitutionId = c.Param("id")

	res, err := h.postClient.GetAllPostByInstitutionID(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *PostHTTPHandler) UpdatePost(c echo.Context) error {
	req := new(pb.UpdatePostRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, httputil.HTTPError{
			Message: "Invalid request body",
		})
	}

	res, err := h.postClient.UpdatePost(c.Request().Context(), &pb.UpdatePostRequest{
		PostId:     c.Param("id"),
		Title:      req.Title,
		Body:       req.Body,
		DateStart:  req.DateStart,
		DateEnd:    req.DateEnd,
		FundTarget: req.FundTarget,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, res)
}

func (h *PostHTTPHandler) DeletePost(c echo.Context) error {
	req := new(pb.DeletePostRequest)
	req.PostId = c.Param("id")

	_, err := h.postClient.DeletePost(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
