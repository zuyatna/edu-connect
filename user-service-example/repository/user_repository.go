package repository

import (
	"context"

	"github.com/zuyatna/edu-connect/user-service/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	RegisterUser(ctx context.Context, user *model.User) (*model.User, error)
	LoginUser(ctx context.Context, email, password string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) (*model.User, error)
	UpdateOrderCountUser(ctx context.Context, id string) error
	DeleteUser(ctx context.Context, id string) error
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
	user := new(model.User)
	err := r.db.Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user := new(model.User)
	err := r.db.Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := new(model.User)
	err := r.db.Where("email = ?", email).First(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *model.User) (*model.User, error) {
	err := r.db.Save(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) UpdateOrderCountUser(ctx context.Context, id string) error {
	user := new(model.User)
	err := r.db.Where("id = ?", id).First(user).Error
	if err != nil {
		return err
	}

	user.OrderCount++
	err = r.db.Save(user).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	user := new(model.User)
	err := r.db.Where("id = ?", id).First(user).Error
	if err != nil {
		return err
	}

	err = r.db.Delete(user).Error
	if err != nil {
		return err
	}

	return nil
}
