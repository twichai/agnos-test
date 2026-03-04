package repository

import "context"

type HospitalRepository interface {
	GetHospitalByID(ctx context.Context, hospitalID string) (string, error)
}
