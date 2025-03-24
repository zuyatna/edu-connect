package handler

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/zuyatna/edu-connect/institution-service/middlewares"
	"github.com/zuyatna/edu-connect/institution-service/model"
	pb "github.com/zuyatna/edu-connect/institution-service/pb/post"
	"github.com/zuyatna/edu-connect/institution-service/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IPostHandler interface {
	CreatePost(ctx context.Context, req *pb.CreatePostRequest) (*pb.PostResponse, error)
	GetPostByID(ctx context.Context, req *pb.GetPostByIDRequest) (*pb.PostResponse, error)
	GetAllPostByInstitutionID(ctx context.Context, req *pb.GetAllPostByInstitutionIDRequest) (*pb.GetAllPostByInstitutionIDResponse, error)
	AddPostFundAchieved(ctx context.Context, req *pb.AddPostFundAchievedRequest) (*pb.AddPostFundAchievedResponse, error)
	UpdatePost(ctx context.Context, req *pb.UpdatePostRequest) (*pb.PostResponse, error)
	DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error)
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

	// More flexible date parsing that handles both date-only and full RFC3339 formats
	var dateStart, dateEnd time.Time

	// Try RFC3339 first, then fallback to date-only format
	dateStart, err = time.Parse(time.RFC3339, req.DateStart)
	if err != nil {
		// Try simple date format
		dateStart, err = time.Parse("2006-01-02", req.DateStart)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid date_start format, expected YYYY-MM-DD or RFC3339: %v", err)
		}
	}

	dateEnd, err = time.Parse(time.RFC3339, req.DateEnd)
	if err != nil {
		// Try simple date format
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
		DateStart:    createdPost.DateStart.String(),
		DateEnd:      createdPost.DateEnd.String(),
		FundTarget:   float32(createdPost.FundTarget),
		FuncAchieved: float32(createdPost.FundAchieved),
	}, nil
}

func (s *PostServer) GetPostByID(ctx context.Context, req *pb.GetPostByIDRequest) (*pb.PostResponse, error) {
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post ID format: %v", err)
	}

	post, err := s.postUsecase.GetPostByID(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get post by ID error: %v", err)
	}

	return &pb.PostResponse{
		PostId:       post.PostID.String(),
		Title:        post.Title,
		Body:         post.Body,
		DateStart:    post.DateStart.String(),
		DateEnd:      post.DateEnd.String(),
		FundTarget:   float32(post.FundTarget),
		FuncAchieved: float32(post.FundAchieved),
	}, nil
}

func (s *PostServer) GetAllPostByInstitutionID(ctx context.Context, req *pb.GetAllPostByInstitutionIDRequest) (*pb.GetAllPostByInstitutionIDResponse, error) {
	institutionID, err := uuid.Parse(req.InstitutionId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid institution ID format: %v", err)
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
			DateStart:    post.DateStart.String(),
			DateEnd:      post.DateEnd.String(),
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

	// More flexible date parsing that handles both date-only and full RFC3339 formats
	var dateStart, dateEnd time.Time

	// Try RFC3339 first, then fallback to date-only format
	dateStart, err = time.Parse(time.RFC3339, req.DateStart)
	if err != nil {
		// Try simple date format
		dateStart, err = time.Parse("2006-01-02", req.DateStart)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid date_start format, expected YYYY-MM-DD or RFC3339: %v", err)
		}
	}

	dateEnd, err = time.Parse(time.RFC3339, req.DateEnd)
	if err != nil {
		// Try simple date format
		dateEnd, err = time.Parse("2006-01-02", req.DateEnd)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid date_end format, expected YYYY-MM-DD or RFC3339: %v", err)
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
		DateStart:    updatedPost.DateStart.String(),
		DateEnd:      updatedPost.DateEnd.String(),
		FundTarget:   float32(updatedPost.FundTarget),
		FuncAchieved: float32(updatedPost.FundAchieved),
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

func (s *PostServer) DeletePost(ctx context.Context, req *pb.DeletePostRequest) (*pb.DeletePostResponse, error) {
	postID, err := uuid.Parse(req.PostId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid post ID format: %v", err)
	}

	err = s.postUsecase.DeletePost(ctx, postID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete post error: %v", err)
	}

	return &pb.DeletePostResponse{
		Message: "Post deleted successfully",
	}, nil
}
