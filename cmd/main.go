package main

import (
	"log"
	"twichai/agnos-test/api/handler"
	"twichai/agnos-test/api/routes"
	"twichai/agnos-test/internal/config"
	"twichai/agnos-test/internal/database"
	patientRepo "twichai/agnos-test/internal/repository"
	staffRepo "twichai/agnos-test/internal/repository"
	patientUsecase "twichai/agnos-test/pkg/patient/service"
	staffUsecase "twichai/agnos-test/pkg/staff/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgresGorm(cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	patientRepository := patientRepo.NewPatientGormRepository(db)
	patientUsecase := patientUsecase.NewPatientUsecase(patientRepository)
	patientHandler := handler.NewPatientHandler(patientUsecase)

	staffRepository := staffRepo.NewStaffGormRepository(db)
	staffUsecase := staffUsecase.NewStaffService(staffRepository)
	staffHandler := handler.NewStaffHandler(staffUsecase)

	router := gin.Default()
	routes.RegisterPatientRoutes(router, patientHandler)
	routes.RegisterStaffRoutes(router, staffHandler)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run(":" + cfg.AppPort)
}
