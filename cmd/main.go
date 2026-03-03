package main

import (
	"log"
	"twichai/agnos-test/api/handler"
	"twichai/agnos-test/api/routes"
	"twichai/agnos-test/internal/database"
	"twichai/agnos-test/internal/repository"
	usecase "twichai/agnos-test/pkg/patient/service"

	"github.com/gin-gonic/gin"
)

func main() {
	db, err := database.NewPostgresGormFromEnv()
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	patientRepository := repository.NewPatientGormRepository(db)
	patientUsecase := usecase.NewPatientUsecase(patientRepository)
	patientHandler := handler.NewPatientHandler(patientUsecase)

	router := gin.Default()
	routes.RegisterPatientRoutes(router, patientHandler)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listens on 0.0.0.0:8080 by default
}
