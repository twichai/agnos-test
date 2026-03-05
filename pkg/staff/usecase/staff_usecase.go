package usecase

import (
	"context"
	"errors"
	"twichai/agnos-test/api/presenter/patient"
	"twichai/agnos-test/pkg/staff/entity"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type StaffUseCase interface {
	Create(ctx context.Context, staff *entity.CreateStaffRequest) (*entity.Staff, error)
	Login(ctx context.Context, staff *entity.StaffLoginRequest) (*patient.LoginStaffPresenter, error)
}
