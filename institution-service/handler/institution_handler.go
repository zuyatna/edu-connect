package handler

import (
	"context"

	"institution-service/middlewares"
	"institution-service/model"
	pb "institution-service/pb/institution"
	"institution-service/usecase"
	"institution-service/utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IInstitutionHandler interface {
	RegisterInstitution(ctx context.Context, req *pb.RegisterInstitutionRequest) (*pb.InstitutionResponse, error)
	LoginInstitution(ctx context.Context, req *pb.LoginInstitutionRequest) (*pb.LoginInstitutionResponse, error)

	GetInstitutionByID(ctx context.Context, req *pb.GetInstitutionByIDRequest) (*pb.InstitutionResponse, error)
	GetInstitutionByEmail(ctx context.Context, req *pb.GetInstitutionByEmailRequest) (*pb.InstitutionResponse, error)
	UpdateInstitution(ctx context.Context, req *pb.UpdateInstitutionRequest) (*pb.InstitutionResponse, error)
	DeleteInstitution(ctx context.Context, req *pb.DeleteInstitutionRequest) (*pb.DeleteInstitutionResponse, error)
}

type InstitutionServer struct {
	pb.UnimplementedInstitutionServiceServer
	userUsecase usecase.IInstitutionUsecase
}

func NewInstitutionHandler(userUsecase usecase.IInstitutionUsecase) *InstitutionServer {
	return &InstitutionServer{
		userUsecase: userUsecase,
	}
}

func (s *InstitutionServer) RegisterInstitution(ctx context.Context, req *pb.RegisterInstitutionRequest) (*pb.InstitutionResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register institution: %v", err)
	}
	req.Password = string(hashedPassword)

	institution := &model.Institution{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Address:  req.Address,
		Phone:    req.Phone,
		Website:  req.Website,
	}

	institution, err = s.userUsecase.RegisterInstitution(ctx, institution)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register institution: %v", err)
	}

	return &pb.InstitutionResponse{
		InstitutionId: institution.InstitutionID.String(),
		Name:          institution.Name,
		Email:         institution.Email,
	}, nil
}

func (s *InstitutionServer) LoginInstitution(ctx context.Context, req *pb.LoginInstitutionRequest) (*pb.LoginInstitutionResponse, error) {
	institution, err := s.userUsecase.LoginInstitution(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login institution: %v", err)
	}

	token, err := utils.GenerateInstitutionToken(institution.InstitutionID.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login institution: %v", err)
	}

	return &pb.LoginInstitutionResponse{
		Token: token,
	}, nil
}

func (s *InstitutionServer) GetInstitutionByID(ctx context.Context, req *pb.GetInstitutionByIDRequest) (*pb.InstitutionResponse, error) {
	authenticatedInstitutionID, ok := ctx.Value(middlewares.InstitutionIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated institution ID from context")
	}

	institutionID, err := uuid.Parse(authenticatedInstitutionID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid institution ID format: %v", err)
	}

	institution, err := s.userUsecase.GetInstitutionByID(ctx, institutionID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get institution by ID error: %v", err)
	}

	return &pb.InstitutionResponse{
		InstitutionId: institution.InstitutionID.String(),
		Name:          institution.Name,
		Email:         institution.Email,
		Address:       institution.Address,
		Phone:         institution.Phone,
		Website:       institution.Website,
	}, nil
}

func (s *InstitutionServer) GetInstitutionByEmail(ctx context.Context, req *pb.GetInstitutionByEmailRequest) (*pb.InstitutionResponse, error) {
	institution, err := s.userUsecase.GetInstitutionByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get institution by email error: %v", err)
	}

	return &pb.InstitutionResponse{
		InstitutionId: institution.InstitutionID.String(),
		Name:          institution.Name,
		Email:         institution.Email,
	}, nil
}

func (s *InstitutionServer) UpdateInstitution(ctx context.Context, req *pb.UpdateInstitutionRequest) (*pb.InstitutionResponse, error) {
	authenticatedInstitutionID, ok := ctx.Value(middlewares.InstitutionIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated institution ID from context")
	}

	if authenticatedInstitutionID != req.InstitutionId {
		return nil, status.Errorf(codes.PermissionDenied, "unauthorized access")
	}

	institutionID, err := uuid.Parse(req.InstitutionId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid institution ID format: %v", err)
	}

	getInstitution, err := s.userUsecase.GetInstitutionByID(ctx, institutionID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get institution by ID error: %v", err)
	}

	if req.Email != getInstitution.Email {
		_, err := s.userUsecase.GetInstitutionByEmail(ctx, req.Email)
		if err == nil {
			return nil, status.Errorf(codes.AlreadyExists, "email already exists")
		}
	}

	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "update institution error: %v", err)
		}
		req.Password = string(hashedPassword)
	}

	institution := &model.Institution{
		InstitutionID: institutionID,
		Name:          req.Name,
		Email:         req.Email,
		Address:       req.Address,
		Phone:         req.Phone,
		Website:       req.Website,
	}

	institution, err = s.userUsecase.UpdateInstitution(ctx, institution)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update institution error: %v", err)
	}

	return &pb.InstitutionResponse{
		InstitutionId: institution.InstitutionID.String(),
		Name:          institution.Name,
		Email:         institution.Email,
	}, nil
}

func (s *InstitutionServer) DeleteInstitution(ctx context.Context, req *pb.DeleteInstitutionRequest) (*pb.DeleteInstitutionResponse, error) {
	authenticatedInstitutionID, ok := ctx.Value(middlewares.InstitutionIDKey).(string)
	if !ok {
		return nil, status.Errorf(codes.Internal, "failed to get authenticated institution ID from context")
	}

	if authenticatedInstitutionID != req.InstitutionId {
		return nil, status.Errorf(codes.PermissionDenied, "unauthorized access")
	}

	institutionID, err := uuid.Parse(req.InstitutionId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid institution ID format: %v", err)
	}

	err = s.userUsecase.DeleteInstitution(ctx, institutionID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "delete institution error: %v", err)
	}

	return &pb.DeleteInstitutionResponse{
		Message: "Institution deleted successfully",
	}, nil
}
