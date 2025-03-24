package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/zuyatna/edu-connect/institution-service/model"
)

type IPostUsecase interface {
	CreatePost(ctx context.Context, post *model.Post) (*model.Post, error)
	GetPostByID(ctx context.Context, post_id uuid.UUID) (*model.Post, error)
	GetAllPostByInstitutionID(ctx context.Context, institutionID uuid.UUID) ([]model.Post, error)
	UpdatePost(ctx context.Context, post *model.Post) (*model.Post, error)
	DeletePost(ctx context.Context, post_id uuid.UUID) error
	AddPostFundAchieved(ctx context.Context, post_id uuid.UUID, amount float64) (*model.Post, error)
}

type PostUsecase struct {
	postRepository IPostUsecase
}

func NewPostUsecase(postRepository IPostUsecase) *PostUsecase {
	return &PostUsecase{
		postRepository: postRepository,
	}
}

func (u *PostUsecase) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	var e []string

	if post.Title == "" {
		e = append(e, "Title is required")
	}
	if post.Body == "" {
		e = append(e, "Body is required")
	}
	if post.InstitutionID == uuid.Nil {
		e = append(e, "Institution ID is required")
	}
	if post.DateStart.IsZero() {
		e = append(e, "Date Start is required")
	}
	if post.DateEnd.IsZero() {
		e = append(e, "Date End is required")
	}
	if post.DateStart.After(post.DateEnd) {
		e = append(e, "Date Start must be before Date End")
	}
	if post.DateStart.Equal(post.DateEnd) {
		e = append(e, "Date Start must be different from Date End")
	}
	if post.FundTarget <= 0 {
		e = append(e, "Fund Target must be greater than 0")
	}

	if len(e) > 0 {
		return nil, errors.New(strings.Join(e, ", "))
	}

	return u.postRepository.CreatePost(ctx, post)
}

func (u *PostUsecase) GetPostByID(ctx context.Context, post_id uuid.UUID) (*model.Post, error) {
	return u.postRepository.GetPostByID(ctx, post_id)
}

func (u *PostUsecase) GetAllPostByInstitutionID(ctx context.Context, institutionID uuid.UUID) ([]model.Post, error) {
	return u.postRepository.GetAllPostByInstitutionID(ctx, institutionID)
}

func (u *PostUsecase) UpdatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	return u.postRepository.UpdatePost(ctx, post)
}

func (u *PostUsecase) DeletePost(ctx context.Context, post_id uuid.UUID) error {
	return u.postRepository.DeletePost(ctx, post_id)
}

func (u *PostUsecase) AddPostFundAchieved(ctx context.Context, post_id uuid.UUID, amount float64) (*model.Post, error) {
	return u.postRepository.AddPostFundAchieved(ctx, post_id, amount)
}
