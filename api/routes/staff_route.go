package routes

import (
	"twichai/agnos-test/api/handler"

	"github.com/gin-gonic/gin"
)

func RegisterStaffRoutes(router *gin.Engine, staffHandler *handler.StaffHandler) {
	router.POST("/staff/create", staffHandler.Create)
}
