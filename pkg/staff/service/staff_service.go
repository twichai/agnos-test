package service

import (
	"context"
	"twichai/agnos-test/pkg/staff/entity"
	"twichai/agnos-test/pkg/staff/repository"
	"twichai/agnos-test/pkg/staff/usecase"
)

type StaffService struct {
	repo repository.StaffRepository
}

func NewStaffService(repo repository.StaffRepository) usecase.StaffUseCase {
	return &StaffService{repo: repo}
}

// Create implements [usecase.StaffUseCase].
func (s *StaffService) Create(ctx context.Context, staff *entity.CreateStaffRequest) (*entity.Staff, error) {
	staff.Password = "hashed_" + staff.Password
	return s.repo.Create(ctx, staff)
}
