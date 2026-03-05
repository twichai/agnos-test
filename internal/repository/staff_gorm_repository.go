package repository

import (
	"context"
	"twichai/agnos-test/pkg/staff/entity"
	"twichai/agnos-test/pkg/staff/repository"

	"gorm.io/gorm"
)

type StaffGormRepository struct {
	db *gorm.DB
}

func NewStaffGormRepository(db *gorm.DB) repository.StaffRepository {
	return &StaffGormRepository{db: db}
}

// Create implements
func (s *StaffGormRepository) Create(ctx context.Context, staff *entity.CreateStaffRequest) (*entity.Staff, error) {
	newStaff := &entity.Staff{
		ID:         generateUUID(),
		Username:   staff.Username,
		Password:   staff.Password,
		HospitalID: staff.HospitalID,
	}
	err := s.db.WithContext(ctx).Create(newStaff).Error
	if err != nil {
		return nil, err
	}
	return newStaff, nil
}

// GetByUsername implements [repository.StaffRepository].
func (s *StaffGormRepository) GetByUsername(ctx context.Context, username string, hospitalID string) (*entity.Staff, error) {
	var staff entity.Staff
	err := s.db.WithContext(ctx).Joins("Hospital").Where("username = ? AND hospital_id = ?", username, hospitalID).Preload("Hospital").First(&staff).Error
	if err != nil {
		return nil, err
	}
	return &staff, nil
}
