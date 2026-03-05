package routes

import (
	"twichai/agnos-test/api/handler"
	"twichai/agnos-test/api/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterPatientRoutes(router *gin.Engine, patientHandler *handler.PatientHandler, appSecret string) {
	router.GET("/patient/search/:id", patientHandler.GetPatientByID)
	authorized := router.Group("")
	authorized.Use(middleware.AuthMiddleware(appSecret))
	authorized.GET("/patient/search", patientHandler.SearchPatient)

}
