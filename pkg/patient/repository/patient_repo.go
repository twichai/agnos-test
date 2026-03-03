package repository

import (
	"context"
	"twichai/agnos-test/pkg/patient/entity"
)

type PatientRepository interface {
	SearchByID(ctx context.Context, id string, hospitalURL string) (*entity.Patient, error)
}
