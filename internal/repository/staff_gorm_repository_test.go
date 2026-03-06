package repository

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
	"twichai/agnos-test/pkg/staff/entity"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func newMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { sqlDB.Close() })

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

var (
	hospitalID      = "3fa85f64-5717-4562-b3fc-2c963f66afa6"
	otherHospitalID = "3fa85f64-5717-4562-b3fc-2c963f66afa7"
	staffID         = "3fa85f64-5717-4562-b3fc-2c963f66afa8"
)

func staffCols() *sqlmock.Rows {
	return sqlmock.NewRows([]string{
		"id", "username", "password_hash", "hospital_id", "created_at", "updated_at",
		"Hospital__id", "Hospital__name", "Hospital__hospital_url",
	})
}

func TestCreateStaff(t *testing.T) {
	ctx := context.Background()
	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())

	t.Run("Add staff success", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewStaffGormRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "staffs"`).
			WithArgs(sqlmock.AnyArg(), username, "testpassword", hospitalID, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		created, err := repo.Create(ctx, &entity.CreateStaffRequest{
			Username:   username,
			Password:   "testpassword",
			HospitalID: hospitalID,
		})

		require.NoError(t, err)
		require.NotNil(t, created)
		assert.Equal(t, username, created.Username)
		assert.Equal(t, hospitalID, created.HospitalID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Add duplicate staff's username and hospital ID", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewStaffGormRepository(db)

		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "staffs"`).
			WithArgs(sqlmock.AnyArg(), username, "testpassword", hospitalID, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnError(fmt.Errorf(`ERROR: duplicate key value violates unique constraint "staffs_username_hospital_id_key" (SQLSTATE 23505)`))
		mock.ExpectRollback()

		created, err := repo.Create(ctx, &entity.CreateStaffRequest{
			Username:   username,
			Password:   "testpassword",
			HospitalID: hospitalID,
		})

		require.Error(t, err)
		assert.Nil(t, created)
		assert.True(t, strings.Contains(strings.ToLower(err.Error()), "duplicate"), "expected duplicate constraint error, got: %v", err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Add staff duplicate username but different hospital ID", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewStaffGormRepository(db)

		// First insert - hospital A
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "staffs"`).
			WithArgs(sqlmock.AnyArg(), username, "testpassword", hospitalID, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Second insert - hospital B (different hospital, same username is allowed)
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO "staffs"`).
			WithArgs(sqlmock.AnyArg(), username, "testpassword2", otherHospitalID, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		created1, err := repo.Create(ctx, &entity.CreateStaffRequest{
			Username:   username,
			Password:   "testpassword",
			HospitalID: hospitalID,
		})
		require.NoError(t, err)
		require.NotNil(t, created1)

		created2, err := repo.Create(ctx, &entity.CreateStaffRequest{
			Username:   username,
			Password:   "testpassword2",
			HospitalID: otherHospitalID,
		})
		require.NoError(t, err)
		require.NotNil(t, created2)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestGetByUsername(t *testing.T) {
	ctx := context.Background()
	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())

	t.Run("Get existing staff by correct username and hospital ID", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewStaffGormRepository(db)

		// Main query with LEFT JOIN (Hospital data is loaded via join columns)
		mock.ExpectQuery(`SELECT .+ FROM "staffs" LEFT JOIN "hospitals"`).
			WithArgs(username, hospitalID, 1).
			WillReturnRows(staffCols().AddRow(
				staffID, username, "testpassword", hospitalID, time.Now(), time.Now(),
				hospitalID, "hospital-a", "http://hospital-a.com",
			))

		fetched, err := repo.GetByUsername(ctx, username, hospitalID)
		require.NoError(t, err)
		require.NotNil(t, fetched)
		assert.Equal(t, staffID, fetched.ID)
		assert.Equal(t, username, fetched.Username)
		assert.Equal(t, hospitalID, fetched.HospitalID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Get non-existing staff by correct username but wrong hospital ID", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewStaffGormRepository(db)

		mock.ExpectQuery(`SELECT .+ FROM "staffs" LEFT JOIN "hospitals"`).
			WithArgs(username, otherHospitalID, 1).
			WillReturnRows(staffCols()) // empty rows → gorm.ErrRecordNotFound

		fetched, err := repo.GetByUsername(ctx, username, otherHospitalID)
		require.Error(t, err)
		assert.Nil(t, fetched)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Get non-existing staff by wrong username but correct hospital ID", func(t *testing.T) {
		db, mock := newMockDB(t)
		repo := NewStaffGormRepository(db)

		mock.ExpectQuery(`SELECT .+ FROM "staffs" LEFT JOIN "hospitals"`).
			WithArgs("nonexistinguser", hospitalID, 1).
			WillReturnRows(staffCols()) // empty rows → gorm.ErrRecordNotFound

		fetched, err := repo.GetByUsername(ctx, "nonexistinguser", hospitalID)
		require.Error(t, err)
		assert.Nil(t, fetched)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
