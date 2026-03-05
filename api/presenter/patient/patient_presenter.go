package patient

import (
	"twichai/agnos-test/pkg/patient/entity"
)

type PatientPresenter struct {
	FirstNameTh  *string `json:"first_name_th"`
	MiddleNameTh *string `json:"middle_name_th,omitempty"`
	LastNameTh   *string `json:"last_name_th"`
	FirstNameEn  *string `json:"first_name_en,omitempty"`
	MiddleNameEn *string `json:"middle_name_en,omitempty"`
	LastNameEn   *string `json:"last_name_en,omitempty"`

	DateOfBirth *string `json:"date_of_birth"`
	PatientHN   string  `json:"patient_hn"`
	NationalID  *string `json:"national_id,omitempty"`
	PassportID  *string `json:"passport_id,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	Email       *string `json:"email,omitempty"`
	Gender      string  `json:"gender,omitempty"`
}

func ToPatientPresenter(patient *entity.Patient) *PatientPresenter {
	if patient == nil {
		return nil
	}

	respont := &PatientPresenter{
		PatientHN:    patient.PatientHospitals[0].PatientHN,
		FirstNameTh:  patient.FirstNameTH,
		MiddleNameTh: patient.MiddleNameTH,
		LastNameTh:   patient.LastNameTH,
		FirstNameEn:  patient.FirstNameEN,
		MiddleNameEn: patient.MiddleNameEN,
		LastNameEn:   patient.LastNameEN,
		DateOfBirth:  patient.DateOfBirth,
		NationalID:   patient.NationalID,
		PassportID:   patient.PassportID,
		PhoneNumber:  patient.PhoneNumber,
		Email:        patient.Email,
		Gender:       patient.Gender,
	}

	return respont
}

func ToPatientPresenters(patients []*entity.Patient) []*PatientPresenter {
	var presenters []*PatientPresenter
	for _, patient := range patients {
		presenters = append(presenters, ToPatientPresenter(patient))
	}
	return presenters
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
