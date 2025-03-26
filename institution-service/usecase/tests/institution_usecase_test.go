package tests

import (
	"context"
	"errors"
	"institution-service/mocks"
	"institution-service/model"
	"institution-service/usecase"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRegisterInstitution(t *testing.T) {
	t.Run("success - register institution", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockInstitutionRepo := mocks.NewMockIInstitutionRepository(ctrl)
		mockInstitutionUsecase := usecase.NewInstitutionUsecase(mockInstitutionRepo)

		institution := &model.Institution{
			InstitutionID: uuid.New(),
			Name:          "Institution Name",
			Email:         "test@email.com",
			Password:      "test_password",
			Phone:         "08123456789",
			Address:       "Institution Address",
			Website:       "institution.com",
		}

		mockInstitutionRepo.EXPECT().
			GetInstitutionByEmail(gomock.Any(), institution.Email).
			Return(nil, errors.New("institution not found"))

		mockInstitutionRepo.EXPECT().
			RegisterInstitution(gomock.Any(), gomock.Eq(institution)).
			Return(institution, nil)

		ctx := context.Background()
		result, err := mockInstitutionUsecase.RegisterInstitution(ctx, institution)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, institution.InstitutionID, result.InstitutionID)
		assert.Equal(t, institution.Name, result.Name)
		assert.Equal(t, institution.Email, result.Email)
		assert.Equal(t, institution.Phone, result.Phone)
		assert.Equal(t, institution.Address, result.Address)
		assert.Equal(t, institution.Website, result.Website)
	})

	t.Run("failed - duplicate email", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockInstitutionRepo := mocks.NewMockIInstitutionRepository(ctrl)
		mockInstitutionUsecase := usecase.NewInstitutionUsecase(mockInstitutionRepo)

		institution := &model.Institution{
			Name:     "Institution Name",
			Email:    "test@email.com",
			Password: "test_password",
		}

		existingInstitution := &model.Institution{
			InstitutionID: uuid.New(),
			Email:         "test@email.com",
		}

		mockInstitutionRepo.EXPECT().
			GetInstitutionByEmail(gomock.Any(), institution.Email).
			Return(existingInstitution, nil)

		ctx := context.Background()
		result, err := mockInstitutionUsecase.RegisterInstitution(ctx, institution)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("failed - registration error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockInstitutionRepo := mocks.NewMockIInstitutionRepository(ctrl)
		mockInstitutionUsecase := usecase.NewInstitutionUsecase(mockInstitutionRepo)

		institution := &model.Institution{
			Name:     "Institution Name",
			Email:    "test@email.com",
			Password: "test_password",
		}

		expectedErr := errors.New("failed to register institution")

		mockInstitutionRepo.EXPECT().
			GetInstitutionByEmail(gomock.Any(), institution.Email).
			Return(nil, errors.New("institution not found"))

		mockInstitutionRepo.EXPECT().
			RegisterInstitution(gomock.Any(), gomock.Eq(institution)).
			Return(nil, expectedErr)

		ctx := context.Background()
		result, err := mockInstitutionUsecase.RegisterInstitution(ctx, institution)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestGetInstitutionByID(t *testing.T) {
	t.Run("success - get institution by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockInstitutionRepo := mocks.NewMockIInstitutionRepository(ctrl)
		mockInstitutionUsecase := usecase.NewInstitutionUsecase(mockInstitutionRepo)

		institutionID := uuid.New()
		institution := &model.Institution{
			InstitutionID: institutionID,
			Name:          "Institution Name",
			Email:         "test@email.com",
			Phone:         "08123456789",
			Address:       "Institution Address",
			Website:       "institution.com",
		}

		mockInstitutionRepo.EXPECT().
			GetInstitutionByID(gomock.Any(), institutionID).
			Return(institution, nil)

		ctx := context.Background()
		result, err := mockInstitutionUsecase.GetInstitutionByID(ctx, institutionID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, institution.InstitutionID, result.InstitutionID)
		assert.Equal(t, institution.Name, result.Name)
		assert.Equal(t, institution.Email, result.Email)
		assert.Equal(t, institution.Phone, result.Phone)
		assert.Equal(t, institution.Address, result.Address)
		assert.Equal(t, institution.Website, result.Website)
	})

	t.Run("failed - get institution by id", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockInstitutionRepo := mocks.NewMockIInstitutionRepository(ctrl)
		mockInstitutionUsecase := usecase.NewInstitutionUsecase(mockInstitutionRepo)

		institutionID := uuid.New()

		mockInstitutionRepo.EXPECT().
			GetInstitutionByID(gomock.Any(), institutionID).
			Return(nil, errors.New("institution not found"))

		ctx := context.Background()
		result, err := mockInstitutionUsecase.GetInstitutionByID(ctx, institutionID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestUpdateInstitution(t *testing.T) {
	t.Run("success - update institution", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockInstitutionRepo := mocks.NewMockIInstitutionRepository(ctrl)
		mockInstitutionUsecase := usecase.NewInstitutionUsecase(mockInstitutionRepo)

		updatedInstitution := &model.Institution{
			InstitutionID: uuid.New(),
			Name:          "Updated Institution Name",
			Email:         "updated@email.com",
			Phone:         "08987654321",
			Address:       "Updated Institution Address",
			Website:       "updated-institution.com",
		}

		mockInstitutionRepo.EXPECT().
			GetInstitutionByEmail(gomock.Any(), updatedInstitution.Email).
			Return(nil, errors.New("institution not found"))

		mockInstitutionRepo.EXPECT().
			UpdateInstitution(gomock.Any(), gomock.Eq(updatedInstitution)).
			Return(updatedInstitution, nil)

		ctx := context.Background()
		result, err := mockInstitutionUsecase.UpdateInstitution(ctx, updatedInstitution)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, updatedInstitution.InstitutionID, result.InstitutionID)
		assert.Equal(t, updatedInstitution.Name, result.Name)
		assert.Equal(t, updatedInstitution.Email, result.Email)
		assert.Equal(t, updatedInstitution.Phone, result.Phone)
		assert.Equal(t, updatedInstitution.Address, result.Address)
		assert.Equal(t, updatedInstitution.Website, result.Website)
	})

	t.Run("failed - update error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockInstitutionRepo := mocks.NewMockIInstitutionRepository(ctrl)
		mockInstitutionUsecase := usecase.NewInstitutionUsecase(mockInstitutionRepo)

		updatedInstitution := &model.Institution{
			InstitutionID: uuid.New(),
			Name:          "Updated Institution Name",
			Email:         "updated@email.com",
		}

		expectedErr := errors.New("failed to update institution")

		mockInstitutionRepo.EXPECT().
			GetInstitutionByEmail(gomock.Any(), updatedInstitution.Email).
			Return(nil, errors.New("institution not found"))

		mockInstitutionRepo.EXPECT().
			UpdateInstitution(gomock.Any(), gomock.Eq(updatedInstitution)).
			Return(nil, expectedErr)

		ctx := context.Background()
		result, err := mockInstitutionUsecase.UpdateInstitution(ctx, updatedInstitution)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, expectedErr, err)
	})
}

func TestDeleteInstitution(t *testing.T) {
	t.Run("success - delete institution", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockInstitutionRepo := mocks.NewMockIInstitutionRepository(ctrl)
		mockInstitutionUsecase := usecase.NewInstitutionUsecase(mockInstitutionRepo)

		institutionID := uuid.New()

		mockInstitutionRepo.EXPECT().
			DeleteInstitution(gomock.Any(), institutionID).
			Return(nil)

		ctx := context.Background()
		err := mockInstitutionUsecase.DeleteInstitution(ctx, institutionID)

		assert.NoError(t, err)
	})

	t.Run("failed - delete institution", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockInstitutionRepo := mocks.NewMockIInstitutionRepository(ctrl)
		mockInstitutionUsecase := usecase.NewInstitutionUsecase(mockInstitutionRepo)

		institutionID := uuid.New()

		expectedErr := errors.New("failed to delete institution")

		mockInstitutionRepo.EXPECT().
			DeleteInstitution(gomock.Any(), institutionID).
			Return(expectedErr)

		ctx := context.Background()
		err := mockInstitutionUsecase.DeleteInstitution(ctx, institutionID)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
