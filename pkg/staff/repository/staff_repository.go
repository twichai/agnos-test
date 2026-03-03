package repository

import (
	"context"
	"twichai/agnos-test/pkg/staff/entity"
)

type StaffRepository interface {
	Create(ctx context.Context, staff *entity.CreateStaffRequest) (*entity.Staff, error)
}
