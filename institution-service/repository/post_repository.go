package repository

import (
	"context"
	"time"

	"institution-service/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPostRepository interface {
	CreatePost(ctx context.Context, post *model.Post) (*model.Post, error)
	GetAllPost(ctx context.Context) ([]model.Post, error)
	GetPostByID(ctx context.Context, post_id uuid.UUID) (*model.Post, error)
	GetAllPostByInstitutionID(ctx context.Context, institution_id uuid.UUID) ([]model.Post, error)
	UpdatePost(ctx context.Context, post *model.Post) (*model.Post, error)
	DeletePost(ctx context.Context, post_id uuid.UUID) error
	AddPostFundAchieved(ctx context.Context, post_id uuid.UUID, amount float64) (*model.Post, error)
}

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	if err := r.db.Create(post).Error; err != nil {
		return nil, err
	}

	return post, nil
}

func (r *PostRepository) GetAllPost(ctx context.Context) ([]model.Post, error) {
	var posts []model.Post

	err := r.db.Where("deleted_at IS NULL OR deleted_at = ?", "0001-01-01 00:00:00").Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetPostByID(ctx context.Context, post_id uuid.UUID) (*model.Post, error) {
	var post model.Post
	if err := r.db.Where("post_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		post_id, "0001-01-01 00:00:00").First(&post).Error; err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetAllPostByInstitutionID(ctx context.Context, institution_id uuid.UUID) ([]model.Post, error) {
	var posts []model.Post

	err := r.db.Where("institution_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		institution_id, "0001-01-01 00:00:00").Find(&posts).Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) UpdatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	updates := make(map[string]interface{})

	if post.Title != "" {
		updates["title"] = post.Title
	}
	if post.Body != "" {
		updates["body"] = post.Body
	}
	if !post.DateStart.IsZero() {
		updates["date_start"] = post.DateStart
	}
	if !post.DateEnd.IsZero() {
		updates["date_end"] = post.DateEnd
	}

	err := r.db.Model(&post).Where("post_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		post.PostID, "0001-01-01 00:00:00").Updates(updates).Error
	if err != nil {
		return nil, err
	}

	updatedPost, err := r.GetPostByID(ctx, post.PostID)
	if err != nil {
		return nil, err
	}

	return updatedPost, nil
}

func (r *PostRepository) DeletePost(ctx context.Context, post_id uuid.UUID) error {
	err := r.db.Model(&model.Post{}).Where("post_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		post_id, "0001-01-01 00:00:00").Update("deleted_at", time.Now()).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *PostRepository) AddPostFundAchieved(ctx context.Context, post_id uuid.UUID, amount float64) (*model.Post, error) {
	var post model.Post
	err := r.db.Where("post_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		post_id, "0001-01-01 00:00:00").First(&post).Error
	if err != nil {
		return nil, err
	}

	post.FundAchieved += amount
	err = r.db.Model(&post).Where("post_id = ? AND (deleted_at IS NULL OR deleted_at = ?)",
		post_id, "0001-01-01 00:00:00").Update("fund_achieved", post.FundAchieved).Error
	if err != nil {
		return nil, err
	}

	return &post, nil
}
