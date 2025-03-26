package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"institution-service/model"
)

type IInstitutionUsecase interface {
	RegisterInstitution(ctx context.Context, institution *model.Institution) (*model.Institution, error)
	LoginInstitution(ctx context.Context, email, password string) (*model.Institution, error)

	GetInstitutionByID(ctx context.Context, id uuid.UUID) (*model.Institution, error)
	GetInstitutionByEmail(ctx context.Context, email string) (*model.Institution, error)
	UpdateInstitution(ctx context.Context, institution *model.Institution) (*model.Institution, error)
	DeleteInstitution(ctx context.Context, id uuid.UUID) error
}

type InstitutionUsecase struct {
	institutionRepository IInstitutionUsecase
}

func NewInstitutionUsecase(institutionRepository IInstitutionUsecase) *InstitutionUsecase {
	return &InstitutionUsecase{
		institutionRepository: institutionRepository,
	}
}

func (u *InstitutionUsecase) RegisterInstitution(ctx context.Context, institution *model.Institution) (*model.Institution, error) {
	var e []string

	if institution.Name == "" {
		e = append(e, "Name is required")
	}
	if institution.Email == "" {
		e = append(e, "Email is required")
	}

	if !strings.Contains(institution.Email, "@") {
		e = append(e, "email must contain @")
	}

	commonTLDs := []string{".com", ".net", ".org", ".edu", ".co"}
	hasTLD := false
	for _, tld := range commonTLDs {
		if strings.Contains(institution.Email, tld) {
			hasTLD = true
			break
		}
	}
	if !hasTLD {
		e = append(e, "email must contain a valid domain extension (.com, .net, etc)")
	}

	if _, err := u.institutionRepository.GetInstitutionByEmail(ctx, institution.Email); err == nil {
		e = append(e, "Email already exists")
	}
	if len(e) > 0 {
		return nil, errors.New(strings.Join(e, ", "))
	}

	institution, err := u.institutionRepository.RegisterInstitution(ctx, institution)
	if err != nil {
		return nil, err
	}

	return institution, nil
}

func (u *InstitutionUsecase) LoginInstitution(ctx context.Context, email, password string) (*model.Institution, error) {
	var e []string

	if email == "" {
		e = append(e, "Email is required")
	}
	if password == "" {
		e = append(e, "Password is required")
	}

	if len(e) > 0 {
		return nil, errors.New(strings.Join(e, ", "))
	}

	institution, err := u.institutionRepository.LoginInstitution(ctx, email, password)
	if err != nil {
		return nil, err
	}

	return institution, nil
}

func (u *InstitutionUsecase) GetInstitutionByID(ctx context.Context, id uuid.UUID) (*model.Institution, error) {
	return u.institutionRepository.GetInstitutionByID(ctx, id)
}

func (u *InstitutionUsecase) GetInstitutionByEmail(ctx context.Context, email string) (*model.Institution, error) {
	return u.institutionRepository.GetInstitutionByEmail(ctx, email)
}

func (u *InstitutionUsecase) UpdateInstitution(ctx context.Context, institution *model.Institution) (*model.Institution, error) {
	var e []string

	if institution.Email != "" {
		if !strings.Contains(institution.Email, "@") {
			e = append(e, "email must contain @")
		}

		commonTLDs := []string{".com", ".net", ".org", ".edu", ".co"}
		hasTLD := false
		for _, tld := range commonTLDs {
			if strings.Contains(institution.Email, tld) {
				hasTLD = true
				break
			}
		}
		if !hasTLD {
			e = append(e, "email must contain a valid domain extension (.com, .net, etc)")
		}

		if _, err := u.institutionRepository.GetInstitutionByEmail(ctx, institution.Email); err == nil {
			e = append(e, "Email already exists")
		}

	}

	if institution.Password != "" && len(institution.Password) < 6 {
		e = append(e, "Password must be at least 6 characters")
	}

	if len(e) > 0 {
		return nil, errors.New("no updates provided")
	}

	return u.institutionRepository.UpdateInstitution(ctx, institution)
}

func (u *InstitutionUsecase) DeleteInstitution(ctx context.Context, id uuid.UUID) error {
	return u.institutionRepository.DeleteInstitution(ctx, id)
}
