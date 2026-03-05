package handler

import (
	"net/http"
	presenter "twichai/agnos-test/api/presenter/patient"
	"twichai/agnos-test/pkg/patient/entity"
	usecase "twichai/agnos-test/pkg/patient/usecase"

	"github.com/gin-gonic/gin"
)

type PatientHandler struct {
	patientUsecase usecase.PatientUsecase
}

func NewPatientHandler(patientUsecase usecase.PatientUsecase) *PatientHandler {
	return &PatientHandler{patientUsecase: patientUsecase}
}

func (h *PatientHandler) GetPatientByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "patient id is required",
		})
		return
	}

	patient, err := h.patientUsecase.SearchByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	respont := presenter.ToPatientPresenter(patient)

	c.JSON(http.StatusOK, respont)
}

func (h *PatientHandler) SearchPatient(c *gin.Context) {
	nationalID := c.Query("national_id")
	passportID := c.Query("passport_id")
	firstName := c.Query("first_name")
	middleName := c.Query("middle_name")
	lastName := c.Query("last_name")
	dateOfBirth := c.Query("date_of_birth")
	phoneNumber := c.Query("phone_number")
	email := c.Query("email")

	// get hospital id from context (set by auth middleware)
	hospitalID, exists := c.Get("hospital_id")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "hospital id not found in token",
		})
		return
	}

	patients, err := h.patientUsecase.Search(c.Request.Context(), &entity.SeachPatientRequest{
		NationalID:  &nationalID,
		PassportID:  &passportID,
		FirstName:   &firstName,
		MiddleName:  &middleName,
		LastName:    &lastName,
		DateOfBirth: &dateOfBirth,
		PhoneNumber: &phoneNumber,
		Email:       &email,
	}, hospitalID.(string))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, patients)
}
