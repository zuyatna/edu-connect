package tests

import (
	"context"
	"institution-service/model"
	"institution-service/repository"
	"log"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewInstitutionMockDB() (*gorm.DB, sqlmock.Sqlmock) {
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

func TestRegisterInstitution(t *testing.T) {
	t.Run("success - register institution", func(t *testing.T) {
		db, mock := NewInstitutionMockDB()
		repo := repository.NewInstitutionRepository(db)

		testInstitution := &model.Institution{
			Name:     "Test Institution",
			Email:    "test@email.com",
			Password: "test_password",
			Address:  "Test Address",
			Phone:    "1234567890",
			Website:  "testwebsite.com",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "institutions"`).
			WithArgs(
				testInstitution.Name,
				testInstitution.Email,
				testInstitution.Password,
				testInstitution.Address,
				testInstitution.Phone,
				testInstitution.Website,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnRows(
				sqlmock.NewRows([]string{"institution_id"}).AddRow(testInstitution.InstitutionID),
			)
		mock.ExpectCommit()

		ctx := context.Background()
		result, err := repo.RegisterInstitution(ctx, testInstitution)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, testInstitution.Name, result.Name)
		assert.Equal(t, testInstitution.Email, result.Email)
		assert.Equal(t, testInstitution.Address, result.Address)
		assert.Equal(t, testInstitution.Phone, result.Phone)
		assert.Equal(t, testInstitution.Website, result.Website)
	})

	t.Run("failed - register institution - invalid email", func(t *testing.T) {
		db, mock := NewInstitutionMockDB()
		repo := repository.NewInstitutionRepository(db)

		testInstitution := &model.Institution{
			Name:     "Test Institution",
			Email:    "testemail",
			Password: "test_password",
			Address:  "Test Address",
			Phone:    "1234567890",
			Website:  "testwebsite.com",
		}

		mock.ExpectBegin()

		ctx := context.Background()
		result, err := repo.RegisterInstitution(ctx, testInstitution)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("failed - register institution - invalid email domain", func(t *testing.T) {
		db, mock := NewInstitutionMockDB()
		repo := repository.NewInstitutionRepository(db)

		testInstitution := &model.Institution{
			Name:     "Test Institution",
			Email:    "testemail@.com",
			Password: "test_password",
			Address:  "Test Address",
			Phone:    "1234567890",
			Website:  "testwebsite.com",
		}

		mock.ExpectBegin()

		ctx := context.Background()
		result, err := repo.RegisterInstitution(ctx, testInstitution)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("failed - register institution - email already exists", func(t *testing.T) {
		db, mock := NewInstitutionMockDB()
		repo := repository.NewInstitutionRepository(db)

		testInstitution := &model.Institution{
			Name:     "Test Institution",
			Email:    "testemail.com",
			Password: "test_password",
			Address:  "Test Address",
			Phone:    "1234567890",
			Website:  "testwebsite.com",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "institutions"`).
			WithArgs(
				testInstitution.Name,
				testInstitution.Email,
				testInstitution.Password,
				testInstitution.Address,
				testInstitution.Phone,
				testInstitution.Website,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		ctx := context.Background()
		result, err := repo.RegisterInstitution(ctx, testInstitution)

		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("failed - register institution", func(t *testing.T) {
		db, mock := NewInstitutionMockDB()
		repo := repository.NewInstitutionRepository(db)

		testInstitution := &model.Institution{
			Name:     "Test Institution",
			Email:    "testemail.com",
			Password: "test_password",
			Address:  "Test Address",
			Phone:    "1234567890",
			Website:  "testwebsite.com",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`INSERT INTO "institutions"`).
			WithArgs(
				testInstitution.Name,
				testInstitution.Email,
				testInstitution.Password,
				testInstitution.Address,
				testInstitution.Phone,
				testInstitution.Website,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnError(&pq.Error{Code: "23505"})
		mock.ExpectRollback()

		ctx := context.Background()
		result, err := repo.RegisterInstitution(ctx, testInstitution)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetInstitutionByID(t *testing.T) {
	t.Run("success - get institution by id", func(t *testing.T) {
		db, mock := NewInstitutionMockDB()
		repo := repository.NewInstitutionRepository(db)
		testInstitution := &model.Institution{
			InstitutionID: uuid.New(),
			Name:          "Test Institution",
			Email:         "testemail.com",
			Address:       "Test Address",
			Phone:         "1234567890",
			Website:       "testwebsite.com",
		}

		rows := sqlmock.NewRows([]string{"institution_id", "name", "email", "address", "phone", "website"}).
			AddRow(testInstitution.InstitutionID, testInstitution.Name, testInstitution.Email, testInstitution.Address, testInstitution.Phone, testInstitution.Website)

		mock.ExpectQuery(`SELECT .* FROM "institutions" WHERE .*institution_id = \$1.*`).
			WithArgs(testInstitution.InstitutionID, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(rows)

		ctx := context.Background()
		result, err := repo.GetInstitutionByID(ctx, testInstitution.InstitutionID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, testInstitution.Name, result.Name)
		assert.Equal(t, testInstitution.Email, result.Email)
		assert.Equal(t, testInstitution.Address, result.Address)
		assert.Equal(t, testInstitution.Phone, result.Phone)
		assert.Equal(t, testInstitution.Website, result.Website)
	})

	t.Run("failed - get institution by id", func(t *testing.T) {
		db, mock := NewInstitutionMockDB()
		repo := repository.NewInstitutionRepository(db)

		testInstitution := &model.Institution{
			InstitutionID: uuid.New(),
		}

		mock.ExpectQuery(`SELECT .* FROM "institutions" WHERE .*institution_id = \$1.*`).
			WithArgs(testInstitution.InstitutionID, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		ctx := context.Background()
		result, err := repo.GetInstitutionByID(ctx, testInstitution.InstitutionID)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetInstitutionByEmail(t *testing.T) {
	t.Run("success - get institution by email", func(t *testing.T) {
		db, mock := NewInstitutionMockDB()
		repo := repository.NewInstitutionRepository(db)
		testInstitution := &model.Institution{
			InstitutionID: uuid.New(),
			Name:          "Test Institution",
			Email:         "testemail.com",
			Address:       "Test Address",
			Phone:         "1234567890",
			Website:       "testwebsite.com",
		}

		rows := sqlmock.NewRows([]string{"institution_id", "name", "email", "address", "phone", "website"}).
			AddRow(testInstitution.InstitutionID, testInstitution.Name, testInstitution.Email, testInstitution.Address, testInstitution.Phone, testInstitution.Website)

		mock.ExpectQuery(`SELECT .* FROM "institutions" WHERE .*email = \$1.*`).
			WithArgs(testInstitution.Email, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(rows)

		ctx := context.Background()
		result, err := repo.GetInstitutionByEmail(ctx, testInstitution.Email)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, testInstitution.Name, result.Name)
		assert.Equal(t, testInstitution.Email, result.Email)
		assert.Equal(t, testInstitution.Address, result.Address)
		assert.Equal(t, testInstitution.Phone, result.Phone)
		assert.Equal(t, testInstitution.Website, result.Website)
	})

	t.Run("failed - get institution by email", func(t *testing.T) {
		db, mock := NewInstitutionMockDB()
		repo := repository.NewInstitutionRepository(db)

		testInstitution := &model.Institution{
			Email: "testemail.com",
		}

		mock.ExpectQuery(`SELECT .* FROM "institutions" WHERE .*email = \$1.*`).
			WithArgs(testInstitution.Email, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(gorm.ErrRecordNotFound)

		ctx := context.Background()
		result, err := repo.GetInstitutionByEmail(ctx, testInstitution.Email)
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestUpdateInstitution(t *testing.T) {
	t.Run("success - update institution", func(t *testing.T) {
		db, mock := NewInstitutionMockDB()
		repo := repository.NewInstitutionRepository(db)
		testInstitution := &model.Institution{
			InstitutionID: uuid.New(),
			Name:          "Test Institution",
			Email:         "testemail.com",
			Address:       "Test Address",
			Phone:         "1234567890",
			Website:       "testwebsite.com",
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "institutions" SET`).
			WithArgs(
				testInstitution.Address,
				testInstitution.Email,
				testInstitution.Name,
				testInstitution.Phone,
				testInstitution.Website,
				sqlmock.AnyArg(),
				testInstitution.InstitutionID,
				"0001-01-01 00:00:00",
				testInstitution.InstitutionID,
			).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		rows := sqlmock.NewRows([]string{"institution_id", "name", "email", "address", "phone", "website"}).
			AddRow(testInstitution.InstitutionID, testInstitution.Name, testInstitution.Email, testInstitution.Address, testInstitution.Phone, testInstitution.Website)
		mock.ExpectQuery(`SELECT .* FROM "institutions" WHERE .*institution_id = \$1.*`).
			WithArgs(testInstitution.InstitutionID, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(rows)

		ctx := context.Background()
		result, err := repo.UpdateInstitution(ctx, testInstitution)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, testInstitution.Name, result.Name)
	})

	t.Run("failed - updated institution", func(t *testing.T) {
		db, mock := NewInstitutionMockDB()
		repo := repository.NewInstitutionRepository(db)
		testInstitution := &model.Institution{
			InstitutionID: uuid.New(),
			Name:          "Test Institution",
			Email:         "testemail.com",
			Address:       "Test Address",
			Phone:         "1234567890",
			Website:       "testwebsite.com",
		}

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "institutions" SET`).
			WithArgs(
				testInstitution.Address,
				testInstitution.Email,
				testInstitution.Name,
				testInstitution.Phone,
				testInstitution.Website,
				sqlmock.AnyArg(),
				testInstitution.InstitutionID,
				"0001-01-01 00:00:00",
				testInstitution.InstitutionID,
			).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		ctx := context.Background()
		result, err := repo.UpdateInstitution(ctx, testInstitution)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestDeleteInstitution(t *testing.T) {
	t.Run("success - delete institution", func(t *testing.T) {
		db, mock := NewInstitutionMockDB()
		repo := repository.NewInstitutionRepository(db)

		institutionID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "institutions" SET "deleted_at"=\$1,"updated_at"=\$2 WHERE .*`).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				institutionID,
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		ctx := context.Background()
		err := repo.DeleteInstitution(ctx, institutionID)

		assert.NoError(t, err)
	})

	t.Run("failed - delete institution", func(t *testing.T) {
		db, mock := NewInstitutionMockDB()
		repo := repository.NewInstitutionRepository(db)

		institutionID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectExec(`UPDATE "institutions" SET "deleted_at"=\$1,"updated_at"=\$2 WHERE .*`).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				institutionID,
				sqlmock.AnyArg(),
			).
			WillReturnError(gorm.ErrRecordNotFound)
		mock.ExpectRollback()

		ctx := context.Background()
		err := repo.DeleteInstitution(ctx, institutionID)

		assert.Error(t, err)
	})
}
