package handler

import (
	"net/http"
	"strings"
	"twichai/agnos-test/pkg/staff/entity"
	"twichai/agnos-test/pkg/staff/usecase"

	"github.com/gin-gonic/gin"
)

type StaffHandler struct {
	usecase usecase.StaffUseCase
}

func NewStaffHandler(usecase usecase.StaffUseCase) *StaffHandler {
	return &StaffHandler{usecase: usecase}
}

func (h *StaffHandler) Create(ginContext *gin.Context) {
	var request entity.CreateStaffRequest
	if err := ginContext.ShouldBindJSON(&request); err != nil {
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body: username, password, hospital_id(uuid) are required",
		})
		return
	}

	request.Username = strings.TrimSpace(request.Username)
	request.Password = strings.TrimSpace(request.Password)
	request.HospitalID = strings.TrimSpace(request.HospitalID)

	if request.Username == "" || request.Password == "" || request.HospitalID == "" {
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "username, password, and hospital_id are required",
		})
		return
	}

	staff, err := h.usecase.Create(ginContext.Request.Context(), &entity.CreateStaffRequest{
		Username:   request.Username,
		Password:   request.Password,
		HospitalID: request.HospitalID,
	})

	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ginContext.JSON(http.StatusOK, staff)
}
