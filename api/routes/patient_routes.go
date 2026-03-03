package routes

import (
	"twichai/agnos-test/api/handler"

	"github.com/gin-gonic/gin"
)

func RegisterPatientRoutes(router *gin.Engine, patientHandler *handler.PatientHandler) {
	router.GET("/patient/search/:id", patientHandler.GetPatientByID)
}
