package tests

import (
	"context"
	"fmt"
	"institution-service/model"
	"institution-service/repository"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostMockDB() (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})

	if err != nil {
		log.Fatalf("An error '%s' was not expected when opening gorm database", err)
	}

	return gormDB, mock
}

func TestCreatePost(t *testing.T) {
	t.Run("success - create post", func(t *testing.T) {
		db, mock := NewPostMockDB()
		repo := repository.NewPostRepository(db)

		post := &model.Post{
			PostID:        uuid.New(),
			InstitutionID: uuid.New(),
			Title:         "Title",
			Body:          "Body",
			DateStart:     time.Now(),
			DateEnd:       time.Now(),
			FundTarget:    1000000,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "posts"`).
			WithArgs(
				post.InstitutionID,
				post.Title,
				post.Body,
				post.DateStart,
				post.DateEnd,
				post.FundTarget,
				float64(0),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				nil,
				post.PostID,
			).
			WillReturnRows(sqlmock.NewRows([]string{"post_id"}).AddRow(post.PostID))
		mock.ExpectCommit()

		ctx := context.Background()
		result, err := repo.CreatePost(ctx, post)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, post.PostID, result.PostID)
		assert.Equal(t, post.InstitutionID, result.InstitutionID)
		assert.Equal(t, post.Title, result.Title)
		assert.Equal(t, post.Body, result.Body)
		assert.Equal(t, post.DateStart, result.DateStart)
		assert.Equal(t, post.DateEnd, result.DateEnd)
		assert.Equal(t, post.FundTarget, result.FundTarget)
	})

	t.Run("failed - create post", func(t *testing.T) {
		db, mock := NewPostMockDB()
		repo := repository.NewPostRepository(db)

		post := &model.Post{
			PostID:        uuid.New(),
			InstitutionID: uuid.New(),
			Title:         "Title",
			Body:          "Body",
			DateStart:     time.Now(),
			DateEnd:       time.Now(),
			FundTarget:    1000000,
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "posts"`).
			WithArgs(
				post.InstitutionID,
				post.Title,
				post.Body,
				post.DateStart,
				post.DateEnd,
				post.FundTarget,
				float64(0),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				nil,
				post.PostID,
			).
			WillReturnError(fmt.Errorf("unexpected error"))
		mock.ExpectRollback()

		ctx := context.Background()
		result, err := repo.CreatePost(ctx, post)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetPostByID(t *testing.T) {
	t.Run("success - get post by id", func(t *testing.T) {
		db, mock := NewPostMockDB()
		repo := repository.NewPostRepository(db)

		postID := uuid.New()
		post := &model.Post{
			PostID:        postID,
			InstitutionID: uuid.New(),
			Title:         "Title",
			Body:          "Body",
			DateStart:     time.Now(),
			DateEnd:       time.Now(),
			FundTarget:    1000000,
		}

		rows := sqlmock.NewRows([]string{"post_id", "institution_id", "title", "body", "date_start", "date_end", "fund_target", "fund_achieved"}).
			AddRow(post.PostID, post.InstitutionID, post.Title, post.Body, post.DateStart, post.DateEnd, post.FundTarget, float64(0))

		mock.ExpectQuery(`SELECT \* FROM "posts" WHERE .+post_id = \$1.+`).
			WithArgs(
				postID,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnRows(rows)

		ctx := context.Background()
		result, err := repo.GetPostByID(ctx, postID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, post.PostID, result.PostID)
		assert.Equal(t, post.InstitutionID, result.InstitutionID)
		assert.Equal(t, post.Title, result.Title)
		assert.Equal(t, post.Body, result.Body)
		assert.Equal(t, post.DateStart, result.DateStart)
		assert.Equal(t, post.DateEnd, result.DateEnd)
		assert.Equal(t, post.FundTarget, result.FundTarget)
	})

	t.Run("failed - get post by id", func(t *testing.T) {
		db, mock := NewPostMockDB()
		repo := repository.NewPostRepository(db)

		postID := uuid.New()

		mock.ExpectQuery(`SELECT \* FROM "posts" WHERE .+post_id = \$1.+`).
			WithArgs(
				postID,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnError(fmt.Errorf("unexpected error"))

		ctx := context.Background()
		result, err := repo.GetPostByID(ctx, postID)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestUpdatePost(t *testing.T) {
	t.Run("success - update post", func(t *testing.T) {
		db, mock := NewPostMockDB()
		repo := repository.NewPostRepository(db)

		postID := uuid.New()
		post := &model.Post{
			PostID:        postID,
			InstitutionID: uuid.New(),
			Title:         "Title",
			Body:          "Body",
			DateStart:     time.Now(),
			DateEnd:       time.Now(),
			FundTarget:    1000000,
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "posts" SET .+ WHERE .+post_id = \$\d+.+`).
			WithArgs(
				post.Body,
				post.DateEnd,
				post.DateStart,
				post.Title,
				sqlmock.AnyArg(),
				post.PostID,
				sqlmock.AnyArg(),
				post.PostID,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		rows := sqlmock.NewRows([]string{"post_id", "institution_id", "title", "body", "date_start", "date_end", "fund_target", "fund_achieved"}).
			AddRow(post.PostID, post.InstitutionID, post.Title, post.Body, post.DateStart, post.DateEnd, post.FundTarget, float64(0))

		mock.ExpectQuery(`SELECT \* FROM "posts" WHERE .+post_id = \$1.+`).
			WithArgs(
				post.PostID,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnRows(rows)

		ctx := context.Background()
		result, err := repo.UpdatePost(ctx, post)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, post.PostID, result.PostID)
		assert.Equal(t, post.InstitutionID, result.InstitutionID)
		assert.Equal(t, post.Title, result.Title)
		assert.Equal(t, post.Body, result.Body)
		assert.Equal(t, post.DateStart, result.DateStart)
		assert.Equal(t, post.DateEnd, result.DateEnd)
		assert.Equal(t, post.FundTarget, result.FundTarget)
	})

	t.Run("failed - update post", func(t *testing.T) {
		db, mock := NewPostMockDB()
		repo := repository.NewPostRepository(db)

		postID := uuid.New()
		post := &model.Post{
			PostID:        postID,
			InstitutionID: uuid.New(),
			Title:         "Title",
			Body:          "Body",
			DateStart:     time.Now(),
			DateEnd:       time.Now(),
			FundTarget:    1000000,
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "posts" SET .+ WHERE .+post_id = \$\d+.+`).
			WithArgs(
				post.Body,
				post.DateEnd,
				post.DateStart,
				post.Title,
				sqlmock.AnyArg(),
				post.PostID,
				sqlmock.AnyArg(),
				post.PostID,
			).
			WillReturnError(fmt.Errorf("unexpected error"))
		mock.ExpectRollback()

		ctx := context.Background()
		result, err := repo.UpdatePost(ctx, post)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestDeletePost(t *testing.T) {
	t.Run("success - delete post", func(t *testing.T) {
		db, mock := NewPostMockDB()
		repo := repository.NewPostRepository(db)

		postID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "posts" SET "deleted_at"=\$1,"updated_at"=\$2 WHERE .+`).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				postID,
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		ctx := context.Background()
		err := repo.DeletePost(ctx, postID)

		assert.NoError(t, err)
	})

	t.Run("failed - delete post", func(t *testing.T) {
		db, mock := NewPostMockDB()
		repo := repository.NewPostRepository(db)

		postID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "posts" SET "deleted_at"=\$1,"updated_at"=\$2 WHERE .+`).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				postID,
				sqlmock.AnyArg(),
			).
			WillReturnError(fmt.Errorf("unexpected error"))
		mock.ExpectRollback()

		ctx := context.Background()
		err := repo.DeletePost(ctx, postID)

		assert.Error(t, err)
	})
}
