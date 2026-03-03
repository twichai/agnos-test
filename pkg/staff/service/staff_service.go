package service

import (
	"context"
	"twichai/agnos-test/pkg/staff/entity"
	"twichai/agnos-test/pkg/staff/repository"
	"twichai/agnos-test/pkg/staff/usecase"

	"golang.org/x/crypto/bcrypt"
)

type StaffService struct {
	repo repository.StaffRepository
}

func NewStaffService(repo repository.StaffRepository) usecase.StaffUseCase {
	return &StaffService{repo: repo}
}

// Create implements [usecase.StaffUseCase].
func (s *StaffService) Create(ctx context.Context, staff *entity.CreateStaffRequest) (*entity.Staff, error) {
	hashedPassword, err := hashPassword(staff.Password)
	if err != nil {
		return nil, err
	}

	staff.Password = hashedPassword
	return s.repo.Create(ctx, staff)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
