package repository

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"
	"twichai/agnos-test/pkg/staff/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=his port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to open database: %v", err))
	}
	return db
}

func mustGetHospitalIDByURL(t *testing.T, db *gorm.DB, hospitalURL string) string {
	t.Helper()

	var hospital struct {
		ID string
	}

	err := db.Table("hospitals").
		Select("id").
		Where("hospital_url = ?", hospitalURL).
		Take(&hospital).Error
	require.NoError(t, err, "failed to find hospital by url %q", hospitalURL)

	return hospital.ID
}

func TestCreateStaff(t *testing.T) {
	db := setupTestDB()
	repo := NewStaffGormRepository(db)

	ctx := context.Background()
	hospitalID := mustGetHospitalIDByURL(t, db, "hospital-a")
	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())

	t.Run("Add staff success.", func(t *testing.T) {
		staff := &entity.CreateStaffRequest{
			Username:   username,
			Password:   "testpassword",
			HospitalID: hospitalID,
		}

		created, err := repo.Create(ctx, staff)
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.Equal(t, staff.Username, created.Username)
		assert.Equal(t, staff.HospitalID, created.HospitalID)

		t.Cleanup(func() {
			err := db.WithContext(ctx).Delete(&entity.Staff{}, "id = ?", created.ID).Error
			assert.NoError(t, err)
		})
	})

	t.Run("Add duplicate staff's username and hospital ID", func(t *testing.T) {
		staff1 := &entity.CreateStaffRequest{
			Username:   username,
			Password:   "testpassword",
			HospitalID: hospitalID,
		}
		staff2 := &entity.CreateStaffRequest{
			Username:   username,
			Password:   "testpassword2",
			HospitalID: hospitalID,
		}

		created1, err := repo.Create(ctx, staff1)
		require.NoError(t, err)
		require.NotNil(t, created1)

		t.Cleanup(func() {
			err := db.WithContext(ctx).Delete(&entity.Staff{}, "id = ?", created1.ID).Error
			assert.NoError(t, err)
		})

		created2, err := repo.Create(ctx, staff2)
		require.Error(t, err)
		assert.Nil(t, created2)
		assert.True(t, strings.Contains(strings.ToLower(err.Error()), "duplicate"), "expected duplicate constraint error, got: %v", err)
	})

	t.Run("Add staff duplicate username but different hospital ID", func(t *testing.T) {
		otherHospitalID := mustGetHospitalIDByURL(t, db, "hospital-b")

		staff1 := &entity.CreateStaffRequest{
			Username:   username,
			Password:   "testpassword",
			HospitalID: hospitalID,
		}
		staff2 := &entity.CreateStaffRequest{
			Username:   username,
			Password:   "testpassword2",
			HospitalID: otherHospitalID,
		}

		created1, err := repo.Create(ctx, staff1)
		require.NoError(t, err)
		require.NotNil(t, created1)

		created2, err := repo.Create(ctx, staff2)
		require.NoError(t, err)
		require.NotNil(t, created2)
	})
}

func TestGetByUsername(t *testing.T) {
	db := setupTestDB()
	repo := NewStaffGormRepository(db)

	ctx := context.Background()
	hospitalID := mustGetHospitalIDByURL(t, db, "hospital-a")
	username := fmt.Sprintf("testuser_%d", time.Now().UnixNano())

	staff := &entity.CreateStaffRequest{
		Username:   username,
		Password:   "testpassword",
		HospitalID: hospitalID,
	}

	created, err := repo.Create(ctx, staff)
	require.NoError(t, err)
	require.NotNil(t, created)

	t.Cleanup(func() {
		err := db.WithContext(ctx).Delete(&entity.Staff{}, "id = ?", created.ID).Error
		assert.NoError(t, err)
	})

	t.Run("Get existing staff by correct username and hospital ID", func(t *testing.T) {
		fetched, err := repo.GetByUsername(ctx, username, hospitalID)
		require.NoError(t, err)
		require.NotNil(t, fetched)
		assert.Equal(t, created.ID, fetched.ID)
		assert.Equal(t, created.Username, fetched.Username)
		assert.Equal(t, created.Password, fetched.Password)
		assert.Equal(t, created.HospitalID, fetched.HospitalID)
	})

	t.Run("Get non-existing staff by correct username but wrong hospital ID", func(t *testing.T) {
		otherHospitalID := mustGetHospitalIDByURL(t, db, "hospital-b")
		fetched, err := repo.GetByUsername(ctx, username, otherHospitalID)
		require.Error(t, err)
		assert.Nil(t, fetched)
	})

	t.Run("Get non-existing staff by wrong username but correct hospital ID", func(t *testing.T) {
		fetched, err := repo.GetByUsername(ctx, "nonexistinguser", hospitalID)
		require.Error(t, err)
		assert.Nil(t, fetched)
	})
}
