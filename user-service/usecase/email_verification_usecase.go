package usecase

import (
	"time"
	"userService/queue"
	"userService/repository"

	customErr "userService/error"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IVerificationUseCase interface {
	GenerateVerification(email string) error
	VerifyToken(token string) error
	ResendVerification(email string) error
}

type verificationUseCase struct {
	userRepo         repository.IUserRepository
	verificationRepo repository.IVerificationRepository
	emailPublisher   queue.IEmailPublisher
}

func NewVerificationUseCase(
	userRepo repository.IUserRepository,
	verificationRepo repository.IVerificationRepository,
	emailPublisher queue.IEmailPublisher,
) IVerificationUseCase {
	return &verificationUseCase{
		userRepo:         userRepo,
		verificationRepo: verificationRepo,
		emailPublisher:   emailPublisher,
	}
}

func (v *verificationUseCase) GenerateVerification(email string) error {
	logger := logrus.WithField("email", email)

	token := uuid.NewString()
	expiresAt := time.Now().Add(15 * time.Minute)

	err := v.verificationRepo.CreateToken(email, token, expiresAt)
	if err != nil {
		logger.WithError(err).Error("Failed to create email verification token")
		return customErr.ErrInternalServer
	}

	err = v.emailPublisher.PublishVerificationToken(email, token)
	if err != nil {
		logger.WithError(err).Error("Failed to publish verification email")
	}

	logger.Info("Verification token generated and published")
	return nil
}

func (v *verificationUseCase) VerifyToken(token string) error {
	logger := logrus.WithField("token", token)

	data, err := v.verificationRepo.ValidateToken(token)
	if err != nil {
		logger.WithError(err).Warn("Invalid or expired token")
		return customErr.ErrVerificationTokenInvalid
	}

	err = v.userRepo.UpdateIsVerified(data.Email, true)
	if err != nil {
		logger.WithError(err).Error("Failed to update user verification status")
		return customErr.ErrInternalServer
	}

	err = v.verificationRepo.MarkTokenUsed(token)
	if err != nil {
		logger.WithError(err).Warn("Token used but failed to mark as used")
	}

	logger.WithField("email", data.Email).Info("User verified successfully")
	return nil
}

func (v *verificationUseCase) ResendVerification(email string) error {
	logger := logrus.WithField("email", email)

	user, err := v.userRepo.GetByEmail(email)
	if err != nil {
		logger.Warn("Resend verification failed: Email not found")
		return customErr.ErrLoginEmailNotFound
	}

	activeToken, err := v.verificationRepo.GetActiveVerificationByEmail(email)
	if err == nil && activeToken != nil {
		logger.Warn("Resend verification denied: Active token exists")
		return customErr.ErrVerificationTokenStillValid
	}

	token := uuid.NewString()
	expiresAt := time.Now().Add(15 * time.Minute)

	err = v.verificationRepo.CreateToken(user.Email, token, expiresAt)
	if err != nil {
		logger.WithError(err).Error("Failed to create verification token")
		return customErr.ErrInternalServer
	}

	err = v.emailPublisher.PublishVerificationToken(user.Email, token)
	if err != nil {
		logger.WithError(err).Error("Failed to publish verification email")
	}

	logger.Info("Verification token resent successfully")
	return nil
}
