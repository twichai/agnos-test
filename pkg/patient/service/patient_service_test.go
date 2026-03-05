package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"twichai/agnos-test/pkg/patient/entity"
	pUsecase "twichai/agnos-test/pkg/patient/usecase"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// --- mock repository ---

type mockPatientRepo struct {
	searchByIDFn func(ctx context.Context, id string) (*entity.Patient, error)
	searchFn     func(ctx context.Context, req *entity.SeachPatientRequest, hospitalID string) ([]*entity.Patient, error)
}

func (m mockPatientRepo) SearchByID(ctx context.Context, id string) (*entity.Patient, error) {
	return m.searchByIDFn(ctx, id)
}

func (m mockPatientRepo) Search(ctx context.Context, req *entity.SeachPatientRequest, hospitalID string) ([]*entity.Patient, error) {
	return m.searchFn(ctx, req, hospitalID)
}

// --- fixtures ---

func samplePatient() *entity.Patient {
	nationalID := "th123456789"
	passportID := "123456789"
	fnameEN := "John"
	lnameEN := "Doe"

	return &entity.Patient{
		ID:          "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		NationalID:  &nationalID,
		PassportID:  &passportID,
		FirstNameEN: &fnameEN,
		LastNameEN:  &lnameEN,
		Gender:      "M",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		PatientHospitals: []entity.PatientHospital{
			{
				PatientID:  "3fa85f64-5717-4562-b3fc-2c963f66afa6",
				HospitalID: "3fa85f64-5717-4562-b3fc-2c963f66afa8",
				PatientHN:  "HN123456",
			},
		},
	}
}

// --- SearchByID ---

func TestPatientService_SearchByID(t *testing.T) {
	tests := []struct {
		description string
		id          string
		repoResult  *entity.Patient
		repoErr     error
		expectErr   error
	}{
		{
			description: "found by id",
			id:          "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			repoResult:  samplePatient(),
			repoErr:     nil,
			expectErr:   nil,
		},
		{
			description: "empty id returns error",
			id:          "",
			repoResult:  nil,
			repoErr:     nil,
			expectErr:   errors.New("patient id is required"),
		},
		{
			description: "whitespace id returns error",
			id:          "   ",
			repoResult:  nil,
			repoErr:     nil,
			expectErr:   errors.New("patient id is required"),
		},
		{
			description: "record not found maps to ErrPatientNotFound",
			id:          "00000000-0000-0000-0000-000000000000",
			repoResult:  nil,
			repoErr:     gorm.ErrRecordNotFound,
			expectErr:   ErrPatientNotFound,
		},
		{
			description: "unexpected db error is propagated",
			id:          "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			repoResult:  nil,
			repoErr:     errors.New("db connection error"),
			expectErr:   errors.New("db connection error"),
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			repo := mockPatientRepo{
				searchByIDFn: func(ctx context.Context, id string) (*entity.Patient, error) {
					return test.repoResult, test.repoErr
				},
			}

			svc := NewPatientUsecase(repo)
			result, err := svc.SearchByID(context.Background(), test.id)

			if test.expectErr != nil {
				assert.Error(t, err)
				assert.Nil(t, result)
				if errors.Is(test.expectErr, ErrPatientNotFound) {
					assert.ErrorIs(t, err, ErrPatientNotFound)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, test.repoResult.ID, result.ID)
			}
		})
	}
}

// --- Search ---

func TestPatientService_Search(t *testing.T) {
	hospitalA := "3fa85f64-5717-4562-b3fc-2c963f66afa8"
	hospitalB := "3fa85f64-5717-4562-b3fc-2c963f66afa9"

	nationalID := "th123456789"
	emptyStr := ""

	tests := []struct {
		description string
		request     *entity.SeachPatientRequest
		hospitalID  string
		repoResult  []*entity.Patient
		repoErr     error
		expectCount int
		expectErr   bool
	}{
		{
			description: "returns all patients for hospital A",
			request:     &entity.SeachPatientRequest{NationalID: &emptyStr},
			hospitalID:  hospitalA,
			repoResult:  []*entity.Patient{samplePatient(), samplePatient()},
			repoErr:     nil,
			expectCount: 2,
			expectErr:   false,
		},
		{
			description: "filter by national id returns 1 patient",
			request:     &entity.SeachPatientRequest{NationalID: &nationalID},
			hospitalID:  hospitalA,
			repoResult:  []*entity.Patient{samplePatient()},
			repoErr:     nil,
			expectCount: 1,
			expectErr:   false,
		},
		{
			description: "hospital with no patients returns empty slice",
			request:     &entity.SeachPatientRequest{NationalID: &emptyStr},
			hospitalID:  hospitalB,
			repoResult:  []*entity.Patient{},
			repoErr:     nil,
			expectCount: 0,
			expectErr:   false,
		},
		{
			description: "repo error is propagated",
			request:     &entity.SeachPatientRequest{NationalID: &emptyStr},
			hospitalID:  hospitalA,
			repoResult:  nil,
			repoErr:     errors.New("db error"),
			expectCount: 0,
			expectErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			repo := mockPatientRepo{
				searchFn: func(ctx context.Context, req *entity.SeachPatientRequest, hospitalID string) ([]*entity.Patient, error) {
					return test.repoResult, test.repoErr
				},
			}

			svc := NewPatientUsecase(repo)
			var _ pUsecase.PatientUsecase = svc

			results, err := svc.Search(context.Background(), test.request, test.hospitalID)

			if test.expectErr {
				assert.Error(t, err)
				assert.Nil(t, results)
			} else {
				assert.NoError(t, err)
				assert.Len(t, results, test.expectCount)
			}
		})
	}
}
