package repository

import (
	"context"
	"errors"
	"testing"
	"time"
	"twichai/agnos-test/pkg/patient/entity"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

var (
	patientID  = "aaaa1111-0000-0000-0000-000000000001"
	nationalID = "1100000000001"
)

func patientCols() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "national_id", "passport_id",
		"first_name_th", "middle_name_th", "last_name_th",
		"first_name_en", "middle_name_en", "last_name_en",
		"date_of_birth", "gender", "phone_number", "email",
		"created_at", "updated_at",
	})
}

func patientHospitalCols() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"patient_id", "hospital_id", "patient_hn"})
}

func TestSearchPatientByID(t *testing.T) {
	ctx := context.Background()

	t.Run("Search existing patient with national ID", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewPatientGormRepository(db)

		// Main query: SearchByID uses Joins + First
		mock.ExpectQuery(`SELECT .+ FROM "patients" JOIN patient_hospitals`).
			WithArgs(nationalID, nationalID, 1).
			WillReturnRows(patientCols().AddRow(
				patientID, nationalID, nil,
				"สมชาย", nil, "ใจดี",
				nil, nil, nil,
				"1990-01-01", "M", nil, nil,
				time.Now(), time.Now(),
			))
		// Preload PatientHospitals
		mock.ExpectQuery(`SELECT \* FROM "patient_hospitals"`).
			WithArgs(patientID).
			WillReturnRows(patientHospitalCols().AddRow(patientID, hospitalID, "HN-001"))

		patient, err := repo.SearchByID(ctx, nationalID)
		require.NoError(t, err)
		require.NotNil(t, patient)
		assert.Equal(t, patientID, patient.ID)
		assert.Equal(t, nationalID, *patient.NationalID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Search non-existing patient with id", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewPatientGormRepository(db)

		mock.ExpectQuery(`SELECT .+ FROM "patients" JOIN patient_hospitals`).
			WithArgs("9999999999999", "9999999999999", 1).
			WillReturnRows(patientCols()) // empty → gorm.ErrRecordNotFound

		patient, err := repo.SearchByID(ctx, "9999999999999")
		require.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound), "expected ErrRecordNotFound, got: %v", err)
		assert.Nil(t, patient)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestSearchPatient(t *testing.T) {
	ctx := context.Background()

	t.Run("Search existing patient with national ID", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewPatientGormRepository(db)

		nid := "1100000000001"

		// Main query: Search uses Joins + Find
		mock.ExpectQuery(`SELECT .+ FROM "patients" JOIN patient_hospitals`).
			WithArgs(hospitalID, nid).
			WillReturnRows(patientCols().AddRow(
				patientID, nid, nil,
				"สมชาย", nil, "ใจดี",
				nil, nil, nil,
				"1990-01-01", "M", nil, nil,
				time.Now(), time.Now(),
			))
		// Preload PatientHospitals
		mock.ExpectQuery(`SELECT \* FROM "patient_hospitals"`).
			WithArgs(patientID).
			WillReturnRows(patientHospitalCols().AddRow(patientID, hospitalID, "HN-001"))

		req := &entity.SeachPatientRequest{NationalID: &nid}
		patients, err := repo.Search(ctx, req, hospitalID)
		require.NoError(t, err)
		require.NotNil(t, patients)
		assert.Len(t, patients, 1)
		assert.Equal(t, nid, *patients[0].NationalID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Search non-existing patient", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewPatientGormRepository(db)

		nid := "9999999999999"

		mock.ExpectQuery(`SELECT .+ FROM "patients" JOIN patient_hospitals`).
			WithArgs(hospitalID, nid).
			WillReturnRows(patientCols()) // empty result

		req := &entity.SeachPatientRequest{NationalID: &nid}
		patients, err := repo.Search(ctx, req, hospitalID)
		require.NoError(t, err)
		require.NotNil(t, patients)
		assert.Len(t, patients, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
