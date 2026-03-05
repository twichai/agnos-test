package repository

import (
	"context"
	"errors"
	"testing"
	"twichai/agnos-test/pkg/patient/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestSearchPatientByID(t *testing.T) {
	db := setupTestDB()
	repo := NewPatientGormRepository(db)
	ctx := context.Background()

	t.Run("Search existing patient with id", func(t *testing.T) {
		patient, err := repo.SearchByID(ctx, "1100000000001")
		require.NoError(t, err)
		require.NotNil(t, patient)
		require.NotNil(t, patient.NationalID)
		assert.Equal(t, "1100000000001", *patient.NationalID)
	})

	t.Run("Search non-existing patient with id", func(t *testing.T) {
		patient, err := repo.SearchByID(ctx, "9999999999999")
		require.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "expected ErrRecordNotFound, got: %v", err)
		assert.Nil(t, patient)
	})
}

func TestSearchPatient(t *testing.T) {
	db := setupTestDB()
	repo := NewPatientGormRepository(db)
	ctx := context.Background()
	hospital := createHospital(t, db, "hospital-a", "hospital-a")

	t.Run("Search existing patient with multiple fields", func(t *testing.T) {
		nationalID := "1100000000001"
		req := &entity.SeachPatientRequest{
			NationalID: &nationalID,
		}
		patients, err := repo.Search(ctx, req, hospital.ID)
		require.NoError(t, err)
		require.NotNil(t, patients)
		assert.Len(t, patients, 1)
		assert.Equal(t, nationalID, *patients[0].NationalID)
	})

	t.Run("Search non-existing patient with multiple fields", func(t *testing.T) {
		nationalID := "9999999999999"
		req := &entity.SeachPatientRequest{
			NationalID: &nationalID,
		}
		patients, err := repo.Search(ctx, req, hospital.ID)
		require.NoError(t, err)
		require.NotNil(t, patients)
		assert.Len(t, patients, 0)
	})
}
