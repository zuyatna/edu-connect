package handler

import (
	"context"
	"time"

	"institution-service/middlewares"
	"institution-service/model"
	pb "institution-service/pb/post"
	"institution-service/usecase"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IPostHandler interface {
	CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.PostResponse, error)
	GetAllPost(ctx context.Context, req *pb.GetAllPostRequest) (*pb.GetAllPostResponse, error)
	GetPostByID(ctx context.Context, req *pb.GetPostByIDRequest) (*pb.PostResponse, error)
	GetAllPostByInstitutionID(ctx context.Context, req *pb.GetAllPostByInstitutionIDRequest) (*pb.GetAllPostByInstitutionIDResponse, error)
	UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.PostResponse, error)
	DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error)
	AddPostFundAchieved(ctx context.Context, req *pb.AddPostFundAchievedRequest) (*pb.AddPostFundAchievedResponse, error)
}

type PostServer struct {
	pb.UnimplementedPostServiceServer
	postUsecase usecase.IPostUsecase
}

func NewPostHandler(postUsecase usecase.IPostUsecase) *PostServer {
	return &PostServer{
		postUsecase: postUsecase,
	}
}

func (s *PostServer) CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.PostResponse, error) {
	authenticatedInstitutionID, ok := ctx.Value(middlewares.InstitutionIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated institution ID from context")
	}

	institutionID, err := uuid.Parse(authenticatedInstitutionID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse authenticated institution ID: %v", err)
	}

	var dateStart, dateEnd time.Time

	dateStart, err = time.Parse(time.RFC3339, req.DateStart)
	if err != nil {
		dateStart, err = time.Parse("2006-01-02", req.DateStart)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid date_start format, expected YYYY-MM-DD or RFC3339: %v", err)
		}
	}

	dateEnd, err = time.Parse(time.RFC3339, req.DateEnd)
	if err != nil {
		dateEnd, err = time.Parse("2006-01-02", req.DateEnd)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid date_end format, expected YYYY-MM-DD or RFC3339: %v", err)
		}
	}

	post := &model.Post{
		Title:         req.Title,
		Body:          req.Body,
		InstitutionID: institutionID,
		DateStart:     dateStart,
		DateEnd:       dateEnd,
		FundTarget:    float64(req.FundTarget),
		FundAchieved:  0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	createdPost, err := s.postUsecase.CreatePost(ctx, post)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create post error: %v", err)
	}

	return &pb.PostResponse{
		PostId:       createdPost.PostID.String(),
		Title:        createdPost.Title,
		Body:         createdPost.Body,
		DateStart:    createdPost.DateStart.Format("2006-01-02"),
		DateEnd:      createdPost.DateEnd.Format("2006-01-02"),
		FundTarget:   float32(createdPost.FundTarget),
		FuncAchieved: float32(createdPost.FundAchieved),
	}, nil
}

func (s *PostServer) GetAllPost(ctx context.Context, req *pb.GetAllPostRequest) (*pb.GetAllPostResponse, error) {
	posts, err := s.postUsecase.GetAllPost(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get all post error: %v", err)
	}

	var postResponses []*pb.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, &pb.PostResponse{
			PostId:       post.PostID.String(),
			Title:        post.Title,
			Body:         post.Body,
			DateStart:    post.DateStart.Format("2006-01-02"),
			DateEnd:      post.DateEnd.Format("2006-01-02"),
			FundTarget:   float32(post.FundTarget),
			FuncAchieved: float32(post.FundAchieved),
		})
	}

	return &pb.GetAllPostResponse{
		Posts: postResponses,
	}, nil
}

func (s *PostServer) GetPostByID(ctx context.Context, req *pb.GetPostByIDRequest) (*pb.PostResponse, error) {
	authenticatedInstitutionID, ok := ctx.Value(middlewares.InstitutionIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated institution ID from context")
	}

	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post ID format: %v", err)
	}

	post, err := s.postUsecase.GetPostByID(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get post by ID error: %v", err)
	}

	if post.InstitutionID.String() != authenticatedInstitutionID {
		return nil, status.Errorf(codes.PermissionDenied, "unauthorized access")
	}

	return &pb.PostResponse{
		PostId:       post.PostID.String(),
		Title:        post.Title,
		Body:         post.Body,
		DateStart:    post.DateStart.Format("2006-01-02"),
		DateEnd:      post.DateEnd.Format("2006-01-02"),
		FundTarget:   float32(post.FundTarget),
		FuncAchieved: float32(post.FundAchieved),
	}, nil
}

func (s *PostServer) GetAllPostByInstitutionID(ctx context.Context, req *pb.GetAllPostByInstitutionIDRequest) (*pb.GetAllPostByInstitutionIDResponse, error) {
	authenticatedInstitutionID, ok := ctx.Value(middlewares.InstitutionIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated institution ID from context")
	}

	institutionID, err := uuid.Parse(req.InstitutionId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid institution ID format: %v", err)
	}

	if institutionID.String() != authenticatedInstitutionID {
		return nil, status.Errorf(codes.PermissionDenied, "unauthorized access")
	}

	posts, err := s.postUsecase.GetAllPostByInstitutionID(ctx, institutionID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get all post by institution ID error: %v", err)
	}

	var postResponses []*pb.PostResponse
	for _, post := range posts {
		postResponses = append(postResponses, &pb.PostResponse{
			PostId:       post.PostID.String(),
			Title:        post.Title,
			Body:         post.Body,
			DateStart:    post.DateStart.Format("2006-01-02"),
			DateEnd:      post.DateEnd.Format("2006-01-02"),
			FundTarget:   float32(post.FundTarget),
			FuncAchieved: float32(post.FundAchieved),
		})
	}

	return &pb.GetAllPostByInstitutionIDResponse{
		Posts: postResponses,
	}, nil
}

func (s *PostServer) UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.PostResponse, error) {
	authenticatedInstitutionID, ok := ctx.Value(middlewares.InstitutionIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated institution ID from context")
	}

	institutionID, err := uuid.Parse(authenticatedInstitutionID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to parse authenticated institution ID: %v", err)
	}

	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post ID format: %v", err)
	}

	getPost, err := s.postUsecase.GetPostByID(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get post by ID error: %v", err)
	}

	if getPost.InstitutionID.String() != authenticatedInstitutionID {
		return nil, status.Errorf(codes.PermissionDenied, "unauthorized access")
	}

	var dateStart, dateEnd time.Time

	if req.DateStart == "" {
		dateStart = getPost.DateStart
	} else {
		dateStart, err = time.Parse(time.RFC3339, req.DateStart)
		if err != nil {
			dateStart, err = time.Parse("2006-01-02", req.DateStart)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid date_start format, expected YYYY-MM-DD or RFC3339: %v", err)
			}
		}
	}

	if req.DateEnd == "" {
		dateEnd = getPost.DateEnd
	} else {
		dateEnd, err = time.Parse(time.RFC3339, req.DateEnd)
		if err != nil {
			dateEnd, err = time.Parse("2006-01-02", req.DateEnd)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid date_end format, expected YYYY-MM-DD or RFC3339: %v", err)
			}
		}
	}

	post := &model.Post{
		PostID:        postID,
		Title:         req.Title,
		Body:          req.Body,
		InstitutionID: institutionID,
		DateStart:     dateStart,
		DateEnd:       dateEnd,
		FundTarget:    float64(req.FundTarget),
		FundAchieved:  0,
		UpdatedAt:     time.Now(),
	}

	updatedPost, err := s.postUsecase.UpdatePost(ctx, post)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update post error: %v", err)
	}

	return &pb.PostResponse{
		PostId:       updatedPost.PostID.String(),
		Title:        updatedPost.Title,
		Body:         updatedPost.Body,
		DateStart:    updatedPost.DateStart.Format("2006-01-02"),
		DateEnd:      updatedPost.DateEnd.Format("2006-01-02"),
		FundTarget:   float32(updatedPost.FundTarget),
		FuncAchieved: float32(updatedPost.FundAchieved),
	}, nil
}

func (s *PostServer) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	authenticatedInstitutionID, ok := ctx.Value(middlewares.InstitutionIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated institution ID from context")
	}

	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post ID format: %v", err)
	}

	getPost, err := s.postUsecase.GetPostByID(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get post by ID error: %v", err)
	}

	if getPost.InstitutionID.String() != authenticatedInstitutionID {
		return nil, status.Errorf(codes.PermissionDenied, "unauthorized access")
	}

	err = s.postUsecase.DeletePost(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete post error: %v", err)
	}

	return &pb.DeletePostResponse{
		Message: "Post deleted successfully",
	}, nil
}

func (s *PostServer) AddPostFundAchieved(ctx context.Context, req *pb.AddPostFundAchievedRequest) (*pb.AddPostFundAchievedResponse, error) {
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post ID format: %v", err)
	}

	amount := float64(req.Amount)

	post, err := s.postUsecase.AddPostFundAchieved(ctx, postID, amount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "add post fund achieved error: %v", err)
	}

	return &pb.AddPostFundAchievedResponse{
		PostId:       post.PostID.String(),
		FuncAchieved: float32(post.FundAchieved),
	}, nil
}
