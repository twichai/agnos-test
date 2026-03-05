package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"twichai/agnos-test/api/presenter/patient"
	"twichai/agnos-test/pkg/staff/entity"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type createStaffRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	HospitalID string `json:"hospital_id"`
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	handler := NewStaffHandler(mockStaffUseCase{})
	router.POST("/staff/create", handler.Create)
	router.POST("/staff/login", handler.Login)
	return router
}

type mockStaffUseCase struct{}

func (m mockStaffUseCase) Create(ctx context.Context, staff *entity.CreateStaffRequest) (*entity.Staff, error) {
	// case error duplicate username and hospital ID can be added here if needed
	if staff.Username == "duplicate" && staff.HospitalID == "3fa85f64-5717-4562-b3fc-2c963f66afa6" {
		return nil, assert.AnError
	}
	return &entity.Staff{
		ID:         "3fa85f64-5717-4562-b3fc-2c963f66afa6",
		Username:   staff.Username,
		HospitalID: staff.HospitalID,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}, nil
}

func (m mockStaffUseCase) Login(ctx context.Context, staff *entity.StaffLoginRequest) (*patient.LoginStaffPresenter, error) {
	if staff.Username == "Jane Doe" && staff.Password == "password123" && staff.HospitalID == "3fa85f64-5717-4562-b3fc-2c963f66afa6" {
		return &patient.LoginStaffPresenter{
			Token: "mocked-jwt-token",
		}, nil
	}
	return nil, assert.AnError
}

func TestStaffHandler_Create(t *testing.T) {
	router := setupRouter()

	// Define test cases
	tests := []struct {
		description  string
		requestBody  createStaffRequest
		expectStatus int
	}{
		{
			description:  "Valid input",
			requestBody:  createStaffRequest{"Jane Doe", "password123", "3fa85f64-5717-4562-b3fc-2c963f66afa6"},
			expectStatus: http.StatusOK,
		},
		{
			description:  "Invalid username",
			requestBody:  createStaffRequest{"", "password123", "3fa85f64-5717-4562-b3fc-2c963f66afa6"},
			expectStatus: http.StatusBadRequest,
		},
		{
			description:  "Invalid password",
			requestBody:  createStaffRequest{"Jane Doe", "", "3fa85f64-5717-4562-b3fc-2c963f66afa6"},
			expectStatus: http.StatusBadRequest,
		},
		{
			description:  "Invalid hospital ID",
			requestBody:  createStaffRequest{"Jane Doe", "password123", ""},
			expectStatus: http.StatusBadRequest,
		},
		{
			description:  "Invalid all fields",
			requestBody:  createStaffRequest{"", "", ""},
			expectStatus: http.StatusBadRequest,
		},
		{
			description:  "Duplicate username and hospital ID",
			requestBody:  createStaffRequest{"duplicate", "password123", "3fa85f64-5717-4562-b3fc-2c963f66afa6"},
			expectStatus: http.StatusInternalServerError,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(test.requestBody)
			req, _ := http.NewRequest("POST", "/staff/create", bytes.NewReader(reqBody))
			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectStatus, w.Code)

		})
	}
}

type loginStaffRequest struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	HospitalID string `json:"hospital_id"`
}

func TestStaffHandler_Login(t *testing.T) {
	router := setupRouter()

	// Define test cases
	tests := []struct {
		description  string
		requestBody  loginStaffRequest
		expectStatus int
	}{
		{
			description:  "login success",
			requestBody:  loginStaffRequest{"Jane Doe", "password123", "3fa85f64-5717-4562-b3fc-2c963f66afa6"},
			expectStatus: http.StatusOK,
		},
		{
			description:  "Invalid username",
			requestBody:  loginStaffRequest{"", "password123", "3fa85f64-5717-4562-b3fc-2c963f66afa6"},
			expectStatus: http.StatusBadRequest,
		},
		{
			description:  "Invalid password",
			requestBody:  loginStaffRequest{"Jane Doe", "", "3fa85f64-5717-4562-b3fc-2c963f66afa6"},
			expectStatus: http.StatusBadRequest,
		},
		{
			description:  "Invalid hospital ID",
			requestBody:  loginStaffRequest{"Jane Doe", "password123", ""},
			expectStatus: http.StatusBadRequest,
		},
		{
			description:  "Invalid all fields",
			requestBody:  loginStaffRequest{"", "", ""},
			expectStatus: http.StatusBadRequest,
		},
		{
			description:  "Login failed",
			requestBody:  loginStaffRequest{"Jane Doe", "wrongpassword", "3fa85f64-5717-4562-b3fc-2c963f66afa6"},
			expectStatus: http.StatusUnauthorized,
		},
	}
	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(test.requestBody)
			req, _ := http.NewRequest("POST", "/staff/login", bytes.NewReader(reqBody))
			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectStatus, w.Code)

			if w.Code == http.StatusOK {
				// should return token
				assert.Contains(t, w.Body.String(), "token")
			}
		})
	}
}
