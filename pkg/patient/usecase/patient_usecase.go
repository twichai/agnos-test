package usecase

import (
	"context"
	"twichai/agnos-test/pkg/patient/entity"
)

type PatientUsecase interface {
	SearchByID(ctx context.Context, id string) (*entity.Patient, error)
	Search(ctx context.Context, req *entity.SeachPatientRequest, hospitalID string) ([]*entity.Patient, error)
}
