package handler

import (
	"context"

	"user-service-example/middlewares"
	"user-service-example/model"
	pb "user-service-example/pb/user"
	"user-service-example/usecase"
	"user-service-example/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IUserHandler interface {
	RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.UserResponse, error)
	LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.UserResponse, error)

	GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.UserResponse, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error)
	UpdateDonateCountUser(ctx context.Context, req *pb.UpdateDonateCountRequest) (*pb.UpdateDonateCountResponse, error)
	DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error)
}

type Server struct {
	pb.UnimplementedUserServiceServer
	userUsecase usecase.IUserUsecase
}

func NewUserHandler(userUsecase usecase.IUserUsecase) *Server {
	return &Server{
		userUsecase: userUsecase,
	}
}

func (s *Server) RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.UserResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}
	req.Password = string(hashedPassword)

	user := &model.User{
		Name:        req.Name,
		Email:       req.Email,
		Password:    req.Password,
		DonateCount: 0,
	}

	user, err = s.userUsecase.RegisterUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &pb.UserResponse{
		Id:    user.UserID.String(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := s.userUsecase.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login user: %v", err)
	}

	token, err := utils.GenerateToken(user.UserID.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login user: %v", err)
	}

	return &pb.LoginUserResponse{
		Token: token,
	}, nil
}
func (s *Server) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.UserResponse, error) {
	authenticatedUserID, ok := ctx.Value(middlewares.UserIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated user ID from context")
	}

	if authenticatedUserID != req.Id {
		return nil, status.Errorf(codes.PermissionDenied, "you can only update your own user data")
	}

	userID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	user, err := s.userUsecase.GetUserByID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "register user error: %v", err)
	}

	return &pb.UserResponse{
		Id:    user.UserID.String(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *Server) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.UserResponse, error) {
	user, err := s.userUsecase.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to get user by email: %v", err)
	}

	return &pb.UserResponse{
		Id:    user.UserID.String(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	authenticatedUserID, ok := ctx.Value(middlewares.UserIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated user ID from context")
	}

	if authenticatedUserID != req.Id {
		return nil, status.Errorf(codes.PermissionDenied, "you can only update your own user data")
	}

	userID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	getUser, err := s.userUsecase.GetUserByID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "register user error: %v", err)
	}

	if req.Email != getUser.Email {
		_, err := s.userUsecase.GetUserByEmail(ctx, req.Email)
		if err == nil {
			return nil, status.Errorf(codes.AlreadyExists, "email already exists")
		}
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash password: %v", err)
		}
		req.Password = string(hashedPassword)
	}

	user := &model.User{
		UserID:   userID,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err = s.userUsecase.UpdateUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &pb.UserResponse{
		Id:    user.UserID.String(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *Server) UpdateDonateCountUser(ctx context.Context, req *pb.UpdateDonateCountRequest) (*pb.UpdateDonateCountResponse, error) {
	authenticatedUserID, ok := ctx.Value(middlewares.UserIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated user ID from context")
	}

	if authenticatedUserID != req.Id {
		return nil, status.Errorf(codes.PermissionDenied, "you can only update your own user data")
	}

	userID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	err = s.userUsecase.UpdateDonateCountUser(ctx, userID, float64(req.DonateCount))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update donate count: %v", err)
	}

	return &pb.UpdateDonateCountResponse{
		Message: "Donate count updated successfully",
	}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	authenticatedUserID, ok := ctx.Value(middlewares.UserIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated user ID from context")
	}

	if authenticatedUserID != req.Id {
		return nil, status.Errorf(codes.PermissionDenied, "you can only update your own user data")
	}

	userID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format: %v", err)
	}

	err = s.userUsecase.DeleteUser(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &pb.DeleteUserResponse{
		Message: "User deleted successfully",
	}, nil
}
