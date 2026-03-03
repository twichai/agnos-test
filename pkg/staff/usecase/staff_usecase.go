package usecase

import (
	"context"
	"twichai/agnos-test/pkg/staff/entity"
)

type StaffUseCase interface {
	Create(ctx context.Context, staff *entity.CreateStaffRequest) (*entity.Staff, error)
}
