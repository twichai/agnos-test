package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"twichai/agnos-test/api/middleware"
	"twichai/agnos-test/pkg/patient/entity"
	patientEntity "twichai/agnos-test/pkg/patient/entity"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func setupPatientRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := NewPatientHandler(mockPatientUsecase{})

	appSecret := "mysecretkey"

	router.GET("/patients/:id", handler.GetPatientByID)
	authorized := router.Group("")
	authorized.Use(middleware.AuthMiddleware(appSecret))
	authorized.GET("/patient/search", handler.SearchPatient)
	return router
}

func makeTestAuthToken(secret string, hospitalID string) string {
	claims := middleware.Claims{
		StaffName:  "tester",
		HospitalID: hospitalID,
		StaffID:    "3fa85f64-5717-4562-b3fc-2c963f66afa7",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

type mockPatientUsecase struct{}

func (m mockPatientUsecase) SearchByID(ctx context.Context, id string) (*patientEntity.Patient, error) {
	if id == "123456789" || id == "th123456789" {
		nativeId := "th123456789"
		passportId := "123456789"
		FnameEN := "John"
		LnameEN := "Doe"

		return &patientEntity.Patient{
			ID:           "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			NationalID:   &nativeId,
			PassportID:   &passportId,
			FirstNameTH:  nil,
			MiddleNameTH: nil,
			LastNameTH:   nil,
			FirstNameEN:  &FnameEN,
			MiddleNameEN: nil,
			LastNameEN:   &LnameEN,
			DateOfBirth:  nil,
			Gender:       "M",
			PhoneNumber:  nil,
			Email:        nil,
			CreatedAt:    time.Time{},
			UpdatedAt:    time.Time{},

			PatientHospitals: []patientEntity.PatientHospital{
				{
					PatientID:  "3fa85f64-5717-4562-b3fc-2c963f66afa6",
					HospitalID: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
					PatientHN:  "HN123456",
				},
			},
		}, nil
	}

	return nil, errors.New("patient not found")
}

func (m mockPatientUsecase) Search(ctx context.Context, req *entity.SeachPatientRequest, hospitalID string) ([]*entity.Patient, error) {
	nativeId := "th123456789"
	passportId := "123456789"
	FnameEN := "John"
	LnameEN := "Doe"

	nativeId2 := "th123456790"
	passportId2 := "123456790"
	FnameEN2 := "John"
	LnameEN2 := "Doe"

	patients := []*entity.Patient{
		{
			ID:           "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			NationalID:   &nativeId,
			PassportID:   &passportId,
			FirstNameTH:  nil,
			MiddleNameTH: nil,
			LastNameTH:   nil,
			FirstNameEN:  &FnameEN,
			MiddleNameEN: nil,
			LastNameEN:   &LnameEN,
			DateOfBirth:  nil,
			Gender:       "M",
			PhoneNumber:  nil,
			Email:        nil,
			CreatedAt:    time.Time{},
			UpdatedAt:    time.Time{},

			PatientHospitals: []patientEntity.PatientHospital{
				{
					PatientID:  "3fa85f64-5717-4562-b3fc-2c963f66afa6",
					HospitalID: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
					PatientHN:  "HN123456",
				},
			},
		},
		{
			ID:           "3fa85f64-5717-4562-b3fc-2c963f66afa7",
			NationalID:   &nativeId2,
			PassportID:   &passportId2,
			FirstNameTH:  nil,
			MiddleNameTH: nil,
			LastNameTH:   nil,
			FirstNameEN:  &FnameEN2,
			MiddleNameEN: nil,
			LastNameEN:   &LnameEN2,
			DateOfBirth:  nil,
			Gender:       "M",
			PhoneNumber:  nil,
			Email:        nil,
			CreatedAt:    time.Time{},
			UpdatedAt:    time.Time{},

			PatientHospitals: []patientEntity.PatientHospital{
				{
					PatientID:  "3fa85f64-5717-4562-b3fc-2c963f66afa7",
					HospitalID: "3fa85f64-5717-4562-b3fc-2c963f66afa6",
					PatientHN:  "HN123457",
				},
			},
		},
	}

	if hospitalID == "3fa85f64-5717-4562-b3fc-2c963f66afa6" && req.NationalID != nil && *req.NationalID == "th123456789" {
		return patients[:1], nil
	}

	if hospitalID == "3fa85f64-5717-4562-b3fc-2c963f66afa6" && (req.NationalID == nil || *req.NationalID == "") {
		return patients, nil
	}

	return nil, nil
}

func TestPatientGetPatientByID(t *testing.T) {
	router := setupPatientRouter()
	tests := []struct {
		description  string
		id           string
		expectStatus int
	}{
		{
			description:  "invalid id wiht national id",
			id:           "th123456789",
			expectStatus: http.StatusOK,
		},
		{
			description:  "invalid id wiht passport id",
			id:           "123456789",
			expectStatus: http.StatusOK,
		},
		{
			description:  "patient not found",
			id:           "00000000-0000-0000-0000-000000000000",
			expectStatus: http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/patients/"+test.id, nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectStatus, w.Code)
			if w.Code == http.StatusOK {
				assert.Contains(t, w.Body.String(), "John")
				assert.Contains(t, w.Body.String(), "Doe")
				assert.Contains(t, w.Body.String(), "th123456789")
				assert.Contains(t, w.Body.String(), "123456789")
				assert.Contains(t, w.Body.String(), "M")
				assert.Contains(t, w.Body.String(), "HN123456")
			}
		})
	}
}

func TestPatientSearchPatient(t *testing.T) {
	router := setupPatientRouter()
	type PatientRequest struct {
		nationalID  string
		passportID  string
		firstName   string
		middleName  string
		lastName    string
		dateOfBirth string
		phoneNumber string
		email       string
	}

	tests := []struct {
		description    string
		patientRequest PatientRequest
		expectStatus   int
		hospitalID     string
		countData      int
	}{
		{
			description:    "search patient with empty request hospital id A",
			patientRequest: PatientRequest{},
			expectStatus:   http.StatusOK,
			hospitalID:     "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			countData:      2,
		},
		{
			description:    "search patient with empty request hospital id B",
			patientRequest: PatientRequest{},
			expectStatus:   http.StatusOK,
			hospitalID:     "3fa85f64-5717-4562-b3fc-2c963f66afa7",
			countData:      0,
		},
		{
			description:    "search patient with national id",
			patientRequest: PatientRequest{nationalID: "th123456789"},
			expectStatus:   http.StatusOK,
			hospitalID:     "3fa85f64-5717-4562-b3fc-2c963f66afa6",
			countData:      1,
		},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			body := map[string]string{}
			if test.patientRequest.nationalID != "" {
				body["national_id"] = test.patientRequest.nationalID
			}
			if test.patientRequest.passportID != "" {
				body["passport_id"] = test.patientRequest.passportID
			}
			if test.patientRequest.firstName != "" {
				body["first_name"] = test.patientRequest.firstName
			}
			if test.patientRequest.middleName != "" {
				body["middle_name"] = test.patientRequest.middleName
			}
			if test.patientRequest.lastName != "" {
				body["last_name"] = test.patientRequest.lastName
			}
			if test.patientRequest.dateOfBirth != "" {
				body["date_of_birth"] = test.patientRequest.dateOfBirth
			}
			if test.patientRequest.phoneNumber != "" {
				body["phone_number"] = test.patientRequest.phoneNumber
			}
			if test.patientRequest.email != "" {
				body["email"] = test.patientRequest.email
			}

			jsonBody, _ := json.Marshal(body)
			req, _ := http.NewRequest(http.MethodGet, "/patient/search", bytes.NewReader(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+makeTestAuthToken("mysecretkey", test.hospitalID))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectStatus, w.Code)
			if w.Code == http.StatusOK {
				var patients []map[string]any
				err := json.Unmarshal(w.Body.Bytes(), &patients)
				assert.NoError(t, err)
				assert.Len(t, patients, test.countData)
			}
		})
	}
}
