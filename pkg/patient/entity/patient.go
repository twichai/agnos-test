package entity

import (
	"time"
	hospitalEntity "twichai/agnos-test/pkg/hospital/entity"
)

type Patient struct {
	ID           string    `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	NationalID   *string   `json:"national_id" gorm:"column:national_id"`
	PassportID   *string   `json:"passport_id" gorm:"column:passport_id"`
	FirstNameTH  *string   `json:"first_name_th" gorm:"column:first_name_th"`
	MiddleNameTH *string   `json:"middle_name_th" gorm:"column:middle_name_th"`
	LastNameTH   *string   `json:"last_name_th" gorm:"column:last_name_th"`
	FirstNameEN  *string   `json:"first_name_en" gorm:"column:first_name_en"`
	MiddleNameEN *string   `json:"middle_name_en" gorm:"column:middle_name_en"`
	LastNameEN   *string   `json:"last_name_en" gorm:"column:last_name_en"`
	DateOfBirth  *string   `json:"date_of_birth" gorm:"column:date_of_birth"`
	Gender       string    `json:"gender" gorm:"column:gender"`
	PhoneNumber  *string   `json:"phone_number" gorm:"column:phone_number"`
	Email        *string   `json:"email" gorm:"column:email"`
	CreatedAt    time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"column:updated_at"`

	// Relationships
	PatientHospitals []PatientHospital         `json:"patient_hospitals" gorm:"foreignKey:PatientID;references:ID"`
	Hotels           []hospitalEntity.Hospital `json:"hospitals" gorm:"many2many:patient_hospitals;joinForeignKey:PatientID;joinReferences:HospitalID"`
}

type PatientHospital struct {
	PatientID  string `json:"patient_id" gorm:"column:patient_id;type:uuid;primaryKey"`
	HospitalID string `json:"hospital_id" gorm:"column:hospital_id;type:uuid;primaryKey"`
	PatientHN  string `json:"patient_hn" gorm:"column:patient_hn"`

	// Relationships
	Hospital hospitalEntity.Hospital `json:"hospital" gorm:"foreignKey:HospitalID;references:ID"`
}

type SeachPatientRequest struct {
	NationalID  *string `json:"national_id"`
	PassportID  *string `json:"passport_id"`
	FirstName   *string `json:"first_name"`
	MiddleName  *string `json:"middle_name"`
	LastName    *string `json:"last_name"`
	DateOfBirth *string `json:"date_of_birth"`
	PhoneNumber *string `json:"phone_number"`
	Email       *string `json:"email"`
	HospitalURL string  `json:"hospital_url" binding:"required"`
}

type CreatePatientRequest struct {
	NationalID  string `json:"national_id" binding:"required"`
	PassportID  string `json:"passport_id" binding:"required"`
	FirstNameTH string `json:"first_name_th" binding:"required"`
	LastNameTH  string `json:"last_name_th" binding:"required"`
	DateOfBirth string `json:"date_of_birth" binding:"required"`
}

func (PatientHospital) TableName() string {
	return "patient_hospitals"
}

func (Patient) TableName() string {
	return "patients"
}
