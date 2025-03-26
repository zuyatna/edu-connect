package grpc

import (
	"context"
	"net"

	"userService/middleware"
	pb "userService/proto/user"
	"userService/repository"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	userRepo repository.IUserRepository
}

func NewGRPCServer(userRepo repository.IUserRepository) *Server {
	return &Server{userRepo: userRepo}
}

func (s *Server) GetUserByToken(ctx context.Context, _ *emptypb.Empty) (*pb.GetUserByTokenResponse, error) {

	logger := log.WithField("source", "grpc").WithField("method", "GetUserByToken")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Warn("Missing metadata in request")
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	token := md["authorization"]
	if len(token) == 0 {
		logger.Warn("Missing authorization token")
		return nil, status.Error(codes.Unauthenticated, "missing authorization token")
	}

	email, err := middleware.ValidateJWT(token[0])
	if err != nil {
		logger.WithError(err).Warn("Invalid JWT token")
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		logger.WithError(err).WithField("email", email).Warn("User not found")
		return nil, status.Error(codes.NotFound, "user not found")
	}

	logger.WithField("email", email).Info("User fetched by token")

	return &pb.GetUserByTokenResponse{
		Id:      int32(user.UserID),
		Name:    user.Name,
		Email:   user.Email,
		Balance: user.Balance,
	}, nil
}

func (s *Server) UpdateUserBalance(ctx context.Context, req *pb.UpdateUserBalanceRequest) (*pb.UpdateUserBalanceResponse, error) {

	logger := log.WithField("source", "grpc").WithField("method", "UpdateUserBalance")

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Warn("Missing metadata in request")
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	token := md["authorization"]
	if len(token) == 0 {
		logger.Warn("Missing authorization token")
		return nil, status.Error(codes.Unauthenticated, "missing authorization token")
	}

	email, err := middleware.ValidateJWT(token[0])
	if err != nil {
		logger.WithError(err).Warn("Invalid JWT token")
		return nil, status.Error(codes.Unauthenticated, "invalid token")
	}

	err = s.userRepo.UpdateBalanceByEmail(email, req.Balance)
	if err != nil {
		logger.WithError(err).WithField("email", email).Error("Failed to update balance")
		return nil, status.Error(codes.Internal, "failed to update balance")
	}

	logger.WithFields(log.Fields{
		"email":   email,
		"balance": req.Balance,
	}).Info("User balance updated successfully")

	return &pb.UpdateUserBalanceResponse{
		Message: "Balance updated successfully",
	}, nil
}

func (s *Server) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserByIdResponse, error) {

	logger := log.WithField("source", "grpc").WithField("method", "GetUserById")

	user, err := s.userRepo.GetByID(uint(req.Id))
	if err != nil {
		logger.WithError(err).WithField("id", req.Id).Warn("User not found by ID")
		return nil, status.Error(codes.NotFound, "user not found")
	}

	logger.WithField("id", req.Id).Info("User fetched by ID")

	return &pb.GetUserByIdResponse{
		Id:      int32(user.UserID),
		Name:    user.Name,
		Email:   user.Email,
		Balance: user.Balance,
	}, nil
}

func StartGRPCServer(userRepo repository.IUserRepository, grpcPort string) {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, NewGRPCServer(userRepo))

	reflection.Register(s)
	log.Println("gRPC Reflection enabled!")

	log.Printf("gRPC Server listening on port %s", grpcPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to start gRPC Server: %v", err)
	}
}
