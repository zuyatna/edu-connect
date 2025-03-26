package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"user-service-example/model"
)

type IUserUsecase interface {
	RegisterUser(ctx context.Context, user *model.User) (*model.User, error)
	LoginUser(ctx context.Context, email, password string) (*model.User, error)

	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateDonateCountUser(ctx context.Context, id uuid.UUID, amount float64) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type UserUsecase struct {
	userRepository IUserUsecase
}

func NewUserUsecase(userRepository IUserUsecase) *UserUsecase {
	return &UserUsecase{
		userRepository: userRepository,
	}
}

func (u *UserUsecase) RegisterUser(ctx context.Context, user *model.User) (*model.User, error) {
	var e []string

	if user.Name == "" {
		e = append(e, "Name is required")
	}
	if user.Email == "" {
		e = append(e, "Email is required")
	}

	if !strings.Contains(user.Email, "@") {
		e = append(e, "email must contain @")
	}

	commonTLDs := []string{".com", ".net", ".org", ".edu", ".co"}
	hasTLD := false
	for _, tld := range commonTLDs {
		if strings.Contains(user.Email, tld) {
			hasTLD = true
			break
		}
	}
	if !hasTLD {
		e = append(e, "email must contain a valid domain extension (.com, .net, etc)")
	}

	if _, err := u.userRepository.GetUserByEmail(ctx, user.Email); err == nil {
		e = append(e, "Email already exists")
	}

	if user.Password == "" {
		e = append(e, "Password is required")
	}
	if user.Password != "" && len(user.Password) < 6 {
		e = append(e, "Password must be at least 6 characters")
	}

	if len(e) > 0 {
		return nil, errors.New(strings.Join(e, ", "))
	}

	return u.userRepository.RegisterUser(ctx, user)
}

func (u *UserUsecase) LoginUser(ctx context.Context, email, password string) (*model.User, error) {
	var e []string

	if email == "" {
		e = append(e, "Email is required")
	}

	if !strings.Contains(email, "@") {
		e = append(e, "email must contain @")
	}

	commonTLDs := []string{".com", ".net", ".org", ".edu", ".co"}
	hasTLD := false
	for _, tld := range commonTLDs {
		if strings.Contains(email, tld) {
			hasTLD = true
			break
		}
	}
	if !hasTLD {
		e = append(e, "email must contain a valid domain extension (.com, .net, etc)")
	}

	if password == "" {
		e = append(e, "Password is required")
	}
	if password != "" && len(password) < 6 {
		e = append(e, "Password must be at least 6 characters")
	}

	if len(e) > 0 {
		return nil, errors.New(strings.Join(e, ", "))
	}

	return u.userRepository.LoginUser(ctx, email, password)
}

func (u *UserUsecase) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	return u.userRepository.GetUserByID(ctx, id)
}

func (u *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return u.userRepository.GetUserByEmail(ctx, email)
}

func (u *UserUsecase) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	var e []string

	if user.Email != "" {
		if !strings.Contains(user.Email, "@") {
			e = append(e, "email must contain @")
		}

		commonTLDs := []string{".com", ".net", ".org", ".edu", ".co"}
		hasTLD := false
		for _, tld := range commonTLDs {
			if strings.Contains(user.Email, tld) {
				hasTLD = true
				break
			}
		}
		if !hasTLD {
			e = append(e, "email must contain a valid domain extension (.com, .net, etc)")
		}

		if _, err := u.userRepository.GetUserByEmail(ctx, user.Email); err == nil {
			e = append(e, "Email already exists")
		}
	}

	if user.Password != "" && len(user.Password) < 6 {
		e = append(e, "Password must be at least 6 characters")
	}

	if len(e) > 0 {
		return nil, errors.New(strings.Join(e, ", "))
	}

	return u.userRepository.UpdateUser(ctx, user)
}

func (u *UserUsecase) UpdateDonateCountUser(ctx context.Context, id uuid.UUID, amount float64) error {
	return u.userRepository.UpdateDonateCountUser(ctx, id, amount)
}

func (u *UserUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return u.userRepository.DeleteUser(ctx, id)
}
