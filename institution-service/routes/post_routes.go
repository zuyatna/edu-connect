package routes

import (
	"net/http"

	"institution-service/httputil"
	pb "institution-service/pb/post"

	"github.com/labstack/echo/v4"
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
	e.GET("/v1/posts", h.GetAllPost)

	groupPost := e.Group("/v1/post")
	groupPost.Use(AuthMiddleware)
	groupPost.POST("", h.CreatePost)
	groupPost.GET("/:id", h.GetPostByID)
	groupPost.GET("/institution/:id", h.GetAllPostByInstitutionID)
	groupPost.PUT("/:id", h.UpdatePost)
	groupPost.DELETE("/:id", h.DeletePost)
}

// GetAllPost godoc
// @Summary      Get all Post.
// @Description  Get all Post without authentication.
// @Tags         Post
// @Accept       json
// @Produce      json
// @Success      200  {object}  model.PostResponse "Success get post data"
// @Failure      401  {object}  httputil.HTTPError "Unauthorized"
// @Failure      404  {object}  httputil.HTTPError "Post not found"
// @Router       /v1/posts [get]
func (h *PostHTTPHandler) GetAllPost(c echo.Context) error {
	res, err := h.postClient.GetAllPost(c.Request().Context(), &pb.GetAllPostRequest{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success get post data",
		"data":    res,
	})
}

// CreatePost godoc
// @Summary      Create a new Post.
// @Description  Create post with title, body, etc.
// @Tags         Post
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header    string  true  "Bearer token"
// @Param        id            path      string    true  "Institution ID"
// @Param        request body model.PostRequest true "Post created details"
// @Success      200 {object} model.PostResponse "Post created successfully"
// @Failure      500 {object} httputil.HTTPError "Internal server error"
// @Router       /v1/post [post]
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

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "Post created successfully",
		"data":    res,
	})
}

// GetPostByID godoc
// @Summary      Get Post by ID.
// @Description  Get a specific post by ID with authorization.
// @Tags         Post
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header    string  true  "Bearer token"
// @Param        id            path      string    true  "Institution ID"
// @Success      200  {object}  model.PostResponse "Success get post data"
// @Failure      401  {object}  httputil.HTTPError "Unauthorized"
// @Failure      404  {object}  httputil.HTTPError "Post not found"
// @Router       /v1/post/{id} [get]
func (h *PostHTTPHandler) GetPostByID(c echo.Context) error {
	req := new(pb.GetPostByIDRequest)
	req.PostId = c.Param("id")

	res, err := h.postClient.GetPostByID(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success get post data",
		"data":    res,
	})
}

// GetAllPostByInstitutionID godoc
// @Summary      Get all Post by Institution ID.
// @Description  Get all post by Institution ID with authorization.
// @Tags         Post
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header    string  true  "Bearer token"
// @Param        id            path      string    true  "Institution ID"
// @Success      200  {object}  model.PostResponse "Success get post data"
// @Failure      401  {object}  httputil.HTTPError "Unauthorized"
// @Failure      404  {object}  httputil.HTTPError "Post not found"
// @Router       /v1/post/institution/{id} [get]
func (h *PostHTTPHandler) GetAllPostByInstitutionID(c echo.Context) error {
	req := new(pb.GetAllPostByInstitutionIDRequest)
	req.InstitutionId = c.Param("id")

	res, err := h.postClient.GetAllPostByInstitutionID(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success get post data",
		"data":    res,
	})
}

// UpdatePost godoc
// @Summary      Update Post.
// @Description  Update an existing post with authorization.
// @Tags         Post
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header    string  true  "Bearer token"
// @Param        id            path      string    true  "Institution ID"
// @Param        request       body      model.PostRequest  true  "Updated post data"
// @Success      200  {object}  model.PostResponse "Success update post data"
// @Failure      401  {object}  httputil.HTTPError "Unauthorized"
// @Failure      404  {object}  httputil.HTTPError "Post not found"
// @Router       /v1/post/{id} [put]
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

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success update post data",
		"data":    res,
	})
}

// DeletePost godoc
// @Summary      Delete Post.
// @Description  Delete an existing post with authorization.
// @Tags         Post
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        Authorization  header    string  true  "Bearer token"
// @Param        id            path      string     true  "Institution ID"
// @Success      200  {object}  model.PostDeleteResponse "Success delete post data"
// @Failure      500  {object}  httputil.HTTPError "Internal server error"
// @Router       /v1/post/{id} [delete]
func (h *PostHTTPHandler) DeletePost(c echo.Context) error {
	req := new(pb.DeletePostRequest)
	req.PostId = c.Param("id")

	_, err := h.postClient.DeletePost(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httputil.HTTPError{
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Success delete post data",
		"data":    map[string]interface{}{},
	})
}
