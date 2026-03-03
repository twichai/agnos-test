package handler

import (
	"net/http"
	presenter "twichai/agnos-test/api/presenter/patient"
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

	patient, err := h.patientUsecase.SearchByID(c.Request.Context(), id, "hospital-a")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	respont := presenter.ToPatientPresenter(patient)

	c.JSON(http.StatusOK, respont)
}
