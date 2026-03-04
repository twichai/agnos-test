package service

import (
	"context"
	"twichai/agnos-test/pkg/hospital/repository"
	"twichai/agnos-test/pkg/hospital/usecase"
)

type HospitalService struct {
	repo repository.HospitalRepository
}

func NewHospitalService(repo repository.HospitalRepository) usecase.HospitalUseCase {
	return &HospitalService{repo: repo}
}

func (s *HospitalService) GetHospitalByID(ctx context.Context, hospitalID string) (string, error) {
	return (s.repo).GetHospitalByID(ctx, hospitalID)
}
