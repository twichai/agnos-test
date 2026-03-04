package usecase

import (
	"context"
	"errors"
	"strings"
	"twichai/agnos-test/pkg/patient/entity"
	"twichai/agnos-test/pkg/patient/repository"
	"twichai/agnos-test/pkg/patient/usecase"

	"gorm.io/gorm"
)

var ErrPatientNotFound = errors.New("patient not found")

type patientUsecase struct {
	repo repository.PatientRepository
}

func NewPatientUsecase(repo repository.PatientRepository) usecase.PatientUsecase {
	return &patientUsecase{repo: repo}
}

func (u *patientUsecase) SearchByID(ctx context.Context, id string, hospitalURL string) (*entity.Patient, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("patient id is required")
	}

	patient, err := u.repo.SearchByID(ctx, id, hospitalURL)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPatientNotFound
		}
		return nil, err
	}

	return patient, nil
}

// Search implements [usecase.PatientUsecase].
func (u *patientUsecase) Search(ctx context.Context, req *entity.SeachPatientRequest, hospitalURL string) ([]*entity.Patient, error) {
	return u.repo.Search(ctx, req, hospitalURL)
}
