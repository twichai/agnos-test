package repository

import (
	"context"
	"twichai/agnos-test/pkg/hospital/entity"
	"twichai/agnos-test/pkg/hospital/repository"

	"gorm.io/gorm"
)

type HospitalGormRepository struct {
	db *gorm.DB
}

func NewHospitalGormRepository(db *gorm.DB) repository.HospitalRepository {
	return &HospitalGormRepository{db: db}
}

// GetHospitalByID implements [repository.HospitalRepository].
func (h *HospitalGormRepository) GetHospitalByID(ctx context.Context, hospitalID string) (string, error) {
	var hospital entity.Hospital
	err := h.db.WithContext(ctx).Where("id = ?", hospitalID).First(&hospital).Error
	if err != nil {
		return "", err
	}
	return hospital.Name, nil
}
