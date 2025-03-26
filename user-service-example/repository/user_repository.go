package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"user-service-example/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	RegisterUser(ctx context.Context, user *model.User) (*model.User, error)
	LoginUser(ctx context.Context, email, password string) (*model.User, error)

	GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateDonateCountUser(ctx context.Context, id uuid.UUID, amount float64) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) RegisterUser(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.db.Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) LoginUser(ctx context.Context, email, password string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		email, "0001-01-01 00:00:00").First(&user).Error; err != nil {
		return nil, err
	}

	if err := user.CompareHashAndPassword(password); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.Where("user_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		id, "0001-01-01 00:00:00").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		email, "0001-01-01 00:00:00").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	updates := make(map[string]interface{})

	if user.Name != "" {
		updates["name"] = user.Name
	}

	if user.Email != "" {
		updates["email"] = user.Email
	}

	if user.Password != "" {
		updates["password"] = user.Password
	}

	if len(updates) == 0 {
		return user, nil
	}

	err := r.db.Model(&model.User{}).Where("user_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		user.UserID, "0001-01-01 00:00:00").Updates(updates).Error
	if err != nil {
		return nil, err
	}

	updatedUser, err := r.GetUserByID(ctx, user.UserID)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

func (r *UserRepository) UpdateDonateCountUser(ctx context.Context, id uuid.UUID, amount float64) error {
	err := r.db.Model(&model.User{}).Where("user_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		id, "0001-01-01 00:00:00").Update("donate_count", amount).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	err := r.db.Model(&model.User{}).Where("user_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		id, "0001-01-01 00:00:00").Update("deleted_at", time.Now()).Error
	if err != nil {
		return err
	}

	return nil
}
