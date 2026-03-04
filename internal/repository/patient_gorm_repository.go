package repository

import (
	"context"
	"twichai/agnos-test/pkg/patient/entity"
	"twichai/agnos-test/pkg/patient/repository"

	"gorm.io/gorm"
)

type PatientGormRepository struct {
	db *gorm.DB
}

func NewPatientGormRepository(db *gorm.DB) repository.PatientRepository {
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

// Search implements [repository.PatientRepository].
func (r *PatientGormRepository) Search(ctx context.Context, req *entity.SeachPatientRequest, hospitalURL string) ([]*entity.Patient, error) {
	var patients []*entity.Patient
	query := r.db.WithContext(ctx).
		Preload("PatientHospitals").
		Joins("JOIN patient_hospitals ON patient_hospitals.patient_id = patients.id").
		Joins("JOIN hospitals ON patient_hospitals.hospital_id = hospitals.id").
		Where("hospitals.hospital_url = ?", hospitalURL)

	if *req.NationalID != "" {
		query = query.Where("patients.national_id = ?", req.NationalID)
	}
	if *req.PassportID != "" {
		query = query.Where("patients.passport_id = ?", req.PassportID)
	}
	if *req.FirstName != "" {
		query = query.Where("patients.first_name_th ILIKE ? OR patients.first_name_en ILIKE ?", "%"+*req.FirstName+"%", "%"+*req.FirstName+"%")
	}
	if *req.MiddleName != "" {
		query = query.Where("patients.middle_name_th ILIKE ? OR patients.middle_name_en ILIKE ?", "%"+*req.MiddleName+"%", "%"+*req.MiddleName+"%")
	}
	if *req.LastName != "" {
		query = query.Where("patients.last_name_th ILIKE ? OR patients.last_name_en ILIKE ?", "%"+*req.LastName+"%", "%"+*req.LastName+"%")
	}
	if *req.DateOfBirth != "" {
		query = query.Where("patients.date_of_birth = ?", req.DateOfBirth)
	}
	if *req.PhoneNumber != "" {
		query = query.Where("patients.phone_number = ?", req.PhoneNumber)
	}
	if *req.Email != "" {
		query = query.Where("patients.email = ?", req.Email)
	}

	if err := query.Find(&patients).Error; err != nil {
		return nil, err
	}

	return patients, nil
}
