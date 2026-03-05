package service

import (
	"context"
	"errors"
	"testing"
	"time"
	hospitalEntity "twichai/agnos-test/pkg/hospital/entity"
	"twichai/agnos-test/pkg/staff/entity"
	"twichai/agnos-test/pkg/staff/usecase"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// --- mock repository ---

type mockStaffRepo struct {
	createFn        func(ctx context.Context, staff *entity.CreateStaffRequest) (*entity.Staff, error)
	getByUsernameFn func(ctx context.Context, username string, hospitalID string) (*entity.Staff, error)
}

func (m mockStaffRepo) Create(ctx context.Context, staff *entity.CreateStaffRequest) (*entity.Staff, error) {
	return m.createFn(ctx, staff)
}

func (m mockStaffRepo) GetByUsername(ctx context.Context, username string, hospitalID string) (*entity.Staff, error) {
	return m.getByUsernameFn(ctx, username, hospitalID)
}

// helper: pre-hash a password the same way the service does
func mustHash(password string) string {
	b, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b)
}

// --- Create ---

func TestStaffService_Create(t *testing.T) {
	tests := []struct {
		description string
		repoErr     error
		expectErr   bool
	}{
		{
			description: "success",
			repoErr:     nil,
			expectErr:   false,
		},
		{
			description: "repo returns error",
			repoErr:     errors.New("db error"),
			expectErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			repo := mockStaffRepo{
				createFn: func(ctx context.Context, req *entity.CreateStaffRequest) (*entity.Staff, error) {
					if test.repoErr != nil {
						return nil, test.repoErr
					}
					return &entity.Staff{
						ID:         "3fa85f64-5717-4562-b3fc-2c963f66afa6",
						Username:   req.Username,
						HospitalID: req.HospitalID,
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
					}, nil
				},
			}

			svc := NewStaffService(repo, "mysecretkey")
			result, err := svc.Create(context.Background(), &entity.CreateStaffRequest{
				Username:   "jane",
				Password:   "password123",
				HospitalID: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			})

			if test.expectErr {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "jane", result.Username)
			}
		})
	}
}

// --- Login ---

func TestStaffService_Login(t *testing.T) {
	validPasswordHash := mustHash("password123")
	hospitalID := "3fa85f64-5717-4562-b3fc-2c963f66afa6"

	existingStaff := &entity.Staff{
		ID:         "3fa85f64-5717-4562-b3fc-2c963f66afa7",
		Username:   "jane",
		Password:   validPasswordHash,
		HospitalID: hospitalID,
		Hospital: hospitalEntity.Hospital{
			ID:   hospitalID,
			Name: "Hospital A",
		},
	}

	tests := []struct {
		description string
		request     entity.StaffLoginRequest
		repoResult  *entity.Staff
		repoErr     error
		appSecret   string
		expectErr   error
		expectToken bool
	}{
		{
			description: "valid credentials",
			request:     entity.StaffLoginRequest{Username: "jane", Password: "password123", HospitalID: hospitalID},
			repoResult:  existingStaff,
			repoErr:     nil,
			appSecret:   "mysecretkey",
			expectErr:   nil,
			expectToken: true,
		},
		{
			description: "user not found",
			request:     entity.StaffLoginRequest{Username: "ghost", Password: "password123", HospitalID: hospitalID},
			repoResult:  nil,
			repoErr:     gorm.ErrRecordNotFound,
			appSecret:   "mysecretkey",
			expectErr:   usecase.ErrInvalidCredentials,
			expectToken: false,
		},
		{
			description: "wrong password",
			request:     entity.StaffLoginRequest{Username: "jane", Password: "wrongpassword", HospitalID: hospitalID},
			repoResult:  existingStaff,
			repoErr:     nil,
			appSecret:   "mysecretkey",
			expectErr:   usecase.ErrInvalidCredentials,
			expectToken: false,
		},
		{
			description: "repo error other than not found",
			request:     entity.StaffLoginRequest{Username: "jane", Password: "password123", HospitalID: hospitalID},
			repoResult:  nil,
			repoErr:     errors.New("db connection error"),
			appSecret:   "mysecretkey",
			expectErr:   errors.New("db connection error"),
			expectToken: false,
		},
		{
			description: "no jwt secret configured",
			request:     entity.StaffLoginRequest{Username: "jane", Password: "password123", HospitalID: hospitalID},
			repoResult:  existingStaff,
			repoErr:     nil,
			appSecret:   "",
			expectErr:   errors.New("jwt app secret is not configured"),
			expectToken: false,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			repo := mockStaffRepo{
				getByUsernameFn: func(ctx context.Context, username string, hospitalID string) (*entity.Staff, error) {
					return test.repoResult, test.repoErr
				},
			}

			svc := NewStaffService(repo, test.appSecret)
			result, err := svc.Login(context.Background(), &test.request)

			if test.expectErr != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
				if test.expectErr == usecase.ErrInvalidCredentials {
					assert.ErrorIs(t, err, usecase.ErrInvalidCredentials)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if test.expectToken {
					assert.NotEmpty(t, result.Token)
				}
			}
		})
	}
}
