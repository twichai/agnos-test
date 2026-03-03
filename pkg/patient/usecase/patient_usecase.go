package usecase

import (
	"context"
	"twichai/agnos-test/pkg/patient/entity"
)

type PatientUsecase interface {
	SearchByID(ctx context.Context, id string, hospitalURL string) (*entity.Patient, error)
}
