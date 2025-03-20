package handler

import (
	"context"

	"github.com/zuyatna/edu-connect/user-service/model"
	pb "github.com/zuyatna/edu-connect/user-service/pb/user"
	"github.com/zuyatna/edu-connect/user-service/usecase"
	"github.com/zuyatna/edu-connect/user-service/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IUserHandler interface {
	RegisterUser(ctx context.Context, req *pb.RegisterUserRequest) (*pb.UserResponse, error)
	LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.UserResponse, error)
	GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.UserResponse, error)
	UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error)
	UpdateOrderCountUser(ctx context.Context, req *pb.UpdateOrderCountRequest) (*pb.UpdateOrderCountResponse, error)
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
		Name:       req.Name,
		Email:      req.Email,
		Password:   req.Password,
		OrderCount: 0,
	}

	user, err = s.userUsecase.RegisterUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &pb.UserResponse{
		Id:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := s.userUsecase.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login user: %v", err)
	}

	token, err := utils.GenerateToken(user.ID.Hex())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login user: %v", err)
	}

	return &pb.LoginUserResponse{
		Token: token,
	}, nil
}

func (s *Server) GetUserByID(ctx context.Context, req *pb.GetUserByIDRequest) (*pb.UserResponse, error) {
	_, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	authenticatedUserID, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated user ID from context")
	}

	if authenticatedUserID != req.Id {
		return nil, status.Errorf(codes.PermissionDenied, "you can only update your own user data")
	}

	user, err := s.userUsecase.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "register user error: %v", err)
	}

	return &pb.UserResponse{
		Id:    user.ID.Hex(),
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
		Id:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	authenticatedUserID, ok := ctx.Value("user_id").(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated user ID from context")
	}

	if authenticatedUserID != req.Id {
		return nil, status.Errorf(codes.PermissionDenied, "you can only update your own user data")
	}

	getUser, err := s.userUsecase.GetUserByID(ctx, req.Id)
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
		ID:       objectID,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	user, err = s.userUsecase.UpdateUser(ctx, user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &pb.UserResponse{
		Id:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
	}, nil
}

func (s *Server) UpdateOrderCountUser(ctx context.Context, req *pb.UpdateOrderCountRequest) (*pb.UpdateOrderCountResponse, error) {
	err := s.userUsecase.UpdateOrderCountUser(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update order count user: %v", err)
	}

	return &pb.UpdateOrderCountResponse{
		Message: "Order count updated successfully",
	}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.userUsecase.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	return &pb.DeleteUserResponse{
		Message: "User deleted successfully",
	}, nil
}
