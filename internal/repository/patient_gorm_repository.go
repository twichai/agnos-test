package repository

import (
	"context"
	"twichai/agnos-test/pkg/patient/entity"

	"gorm.io/gorm"
)

type PatientGormRepository struct {
	db *gorm.DB
}

func NewPatientGormRepository(db *gorm.DB) *PatientGormRepository {
	return &PatientGormRepository{db: db}
}

func (r *PatientGormRepository) SearchByID(ctx context.Context, id string, hospitalURL string) (*entity.Patient, error) {
	var patient entity.Patient
	if err := r.db.WithContext(ctx).
		Preload("PatientHospitals").
		Where("patients.national_id = ? OR patients.passport_id = ?", id, id).
		Joins("JOIN patient_hospitals ON patient_hospitals.patient_id = patients.id").
		Joins("JOIN hospitals ON patient_hospitals.hospital_id = hospitals.id").
		Where("hospitals.hospital_url = ?", hospitalURL).
		First(&patient).Error; err != nil {
		return nil, err
	}

	return &patient, nil
}
