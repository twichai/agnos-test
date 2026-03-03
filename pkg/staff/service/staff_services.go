package service

import (
	"context"
	"twichai/agnos-test/pkg/staff/entity"
	"twichai/agnos-test/pkg/staff/repository"
	"twichai/agnos-test/pkg/staff/usecase"
)

type staffUseCase struct {
	repo repository.StaffRepository
}

func NewStaffUsecase(repo repository.StaffRepository) usecase.StaffUseCase {
	return &staffUseCase{repo: repo}
}

func (s *staffUseCase) Create(ctx context.Context, staff *entity.CreateStaffRequest) (*entity.Staff, error) {
	return s.repo.Create(ctx, staff)
}

// Update implements [usecase.StaffUseCase].
func (s *staffUseCase) Update(ctx context.Context, staff *entity.CreateStaffRequest) (*entity.Staff, error) {
	panic("unimplemented")
}
