package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"institution-service/model"
	"gorm.io/gorm"
)

type IInstitutionRepository interface {
	RegisterInstitution(ctx context.Context, institution *model.Institution) (*model.Institution, error)
	LoginInstitution(ctx context.Context, email, password string) (*model.Institution, error)

	GetInstitutionByID(ctx context.Context, id uuid.UUID) (*model.Institution, error)
	GetInstitutionByEmail(ctx context.Context, email string) (*model.Institution, error)
	UpdateInstitution(ctx context.Context, institution *model.Institution) (*model.Institution, error)
	DeleteInstitution(ctx context.Context, id uuid.UUID) error
}

type InstitutionRepository struct {
	db *gorm.DB
}

func NewInstitutionRepository(db *gorm.DB) *InstitutionRepository {
	return &InstitutionRepository{
		db: db,
	}
}

func (r *InstitutionRepository) RegisterInstitution(ctx context.Context, institution *model.Institution) (*model.Institution, error) {
	if err := r.db.Create(institution).Error; err != nil {
		return nil, err
	}

	return institution, nil
}

func (r *InstitutionRepository) LoginInstitution(ctx context.Context, email, password string) (*model.Institution, error) {
	var institution model.Institution
	if err := r.db.Where("email = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		email, "0001-01-01 00:00:00").First(&institution).Error; err != nil {
		return nil, err
	}

	if err := institution.CompareHashAndPassword(password); err != nil {
		return nil, err
	}

	return &institution, nil
}

func (r *InstitutionRepository) GetInstitutionByID(ctx context.Context, id uuid.UUID) (*model.Institution, error) {
	var institution model.Institution
	if err := r.db.Where("institution_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		id, "0001-01-01 00:00:00").First(&institution).Error; err != nil {
		return nil, err
	}

	return &institution, nil
}

func (r *InstitutionRepository) GetInstitutionByEmail(ctx context.Context, email string) (*model.Institution, error) {
	var institution model.Institution
	if err := r.db.Where("email = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		email, "0001-01-01 00:00:00").First(&institution).Error; err != nil {
		return nil, err
	}

	return &institution, nil
}

func (r *InstitutionRepository) UpdateInstitution(ctx context.Context, institution *model.Institution) (*model.Institution, error) {
	updates := make(map[string]interface{})

	if institution.Name != "" {
		updates["name"] = institution.Name
	}

	if institution.Email != "" {
		updates["email"] = institution.Email
	}

	if institution.Password != "" {
		updates["password"] = institution.Password
	}

	if institution.Address != "" {
		updates["address"] = institution.Address
	}

	if institution.Phone != "" {
		updates["phone"] = institution.Phone
	}

	if institution.Website != "" {
		updates["website"] = institution.Website
	}

	err := r.db.Model(&institution).Where("institution_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		institution.InstitutionID, "0001-01-01 00:00:00").Updates(updates).Error
	if err != nil {
		return nil, err
	}

	updatedInstitution, err := r.GetInstitutionByID(ctx, institution.InstitutionID)
	if err != nil {
		return nil, err
	}

	return updatedInstitution, nil
}

func (r *InstitutionRepository) DeleteInstitution(ctx context.Context, id uuid.UUID) error {
    err := r.db.Model(&model.Institution{}).Where("institution_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
        id, "0001-01-01 00:00:00").Update("deleted_at", time.Now()).Error
    if err != nil {
        return err
    }

    return nil
}
