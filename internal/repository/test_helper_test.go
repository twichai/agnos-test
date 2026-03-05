package repository

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
	hospitalEntity "twichai/agnos-test/pkg/hospital/entity"
	patientEntity "twichai/agnos-test/pkg/patient/entity"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=his port=5432 sslmode=disable TimeZone=Asia/Bangkok"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("Failed to open database: %v", err))
	}
	return db
}

func createHospital(t *testing.T, db *gorm.DB, name string, url string) hospitalEntity.Hospital {
	t.Helper()

	hospital := hospitalEntity.Hospital{
		ID:          generateUUID(),
		Name:        name,
		HospitalURL: url,
	}

	err := db.Create(&hospital).Error
	if err != nil {
		panic(fmt.Sprintf("Failed to create hospital: %v", err))
	}
	require.NoError(t, err, "failed to create hospital")

	return hospital
}

func createMalePatient(t *testing.T, db *gorm.DB, hospitalID string) patientEntity.Patient {
	t.Helper()

	nationalID := fmt.Sprintf("9%012d", time.Now().UnixNano()%1000000000000)
	firstNameTH := fmt.Sprintf("สมชาย%d", rand.Intn(1000))
	lastNameTH := "ใจดี"
	dateOfBirth := "1990-01-01"

	patient := patientEntity.Patient{
		ID:          generateUUID(),
		NationalID:  &nationalID,
		FirstNameTH: &firstNameTH,
		LastNameTH:  &lastNameTH,
		DateOfBirth: &dateOfBirth,
		Gender:      "M",
	}
	return createPatient(t, db, patient, hospitalID)
}

func createFemalePatient(t *testing.T, db *gorm.DB, hospitalID string) patientEntity.Patient {
	t.Helper()

	passportID := fmt.Sprintf("P%08d", time.Now().UnixNano()%100000000)
	firstNameEN := fmt.Sprintf("smith%d", rand.Intn(1000))
	lastNameEN := "jane"
	dateOfBirth := "1992-02-02"

	patient := patientEntity.Patient{
		ID:          generateUUID(),
		PassportID:  &passportID,
		FirstNameEN: &firstNameEN,
		LastNameEN:  &lastNameEN,
		DateOfBirth: &dateOfBirth,
		Gender:      "F",
	}
	return createPatient(t, db, patient, hospitalID)
}

func createPatient(t *testing.T, db *gorm.DB, patient patientEntity.Patient, hospitalID string) patientEntity.Patient {
	t.Helper()

	patientNew := patientEntity.Patient{
		ID:           generateUUID(),
		NationalID:   patient.NationalID,
		FirstNameTH:  patient.FirstNameTH,
		MiddleNameTH: patient.MiddleNameTH,
		LastNameTH:   patient.LastNameTH,
		FirstNameEN:  patient.FirstNameEN,
		MiddleNameEN: patient.MiddleNameEN,
		LastNameEN:   patient.LastNameEN,
		DateOfBirth:  patient.DateOfBirth,
		Gender: func() string {
			if patient.Gender != "" {
				return patient.Gender
			} else {
				return "M"
			}
		}(),
		PhoneNumber: patient.PhoneNumber,
		Email:       patient.Email,
	}

	err := db.Create(&patientNew).Error
	if err != nil {
		panic(fmt.Sprintf("Failed to create patient: %v", err))
	}
	err = db.Create(&patientEntity.PatientHospital{
		PatientID:  patientNew.ID,
		HospitalID: hospitalID,
		PatientHN:  fmt.Sprintf("HN-T-%d", time.Now().UnixNano()),
	}).Error
	if err != nil {
		panic(fmt.Sprintf("Failed to create patient_hospital mapping: %v", err))
	}

	require.NoError(t, err, "failed to create patient")

	return patientNew
}
