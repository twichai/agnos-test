package usecase

import "context"

type HospitalUseCase interface {
	GetHospitalByID(ctx context.Context, hospitalID string) (string, error)
}
