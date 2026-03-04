package handler

import (
	"errors"
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

func (h *StaffHandler) Login(ginContext *gin.Context) {
	var request entity.StaffLoginRequest
	if err := ginContext.ShouldBindJSON(&request); err != nil {
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body: username, password, and hospital_id are required",
		})
		return
	}

	request.Username = strings.TrimSpace(request.Username)
	request.Password = strings.TrimSpace(request.Password)
	request.HospitalID = strings.TrimSpace(request.HospitalID)

	if request.HospitalID == "" {
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "hospital_id is required",
		})
		return
	}

	if request.Username == "" || request.Password == "" {
		ginContext.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "username and password are required",
		})
		return
	}

	staff, err := h.usecase.Login(ginContext.Request.Context(), &entity.StaffLoginRequest{
		Username:   request.Username,
		Password:   request.Password,
		HospitalID: request.HospitalID,
	})
	if err != nil {
		if errors.Is(err, usecase.ErrInvalidCredentials) {
			ginContext.JSON(http.StatusUnauthorized, gin.H{
				"error": "username or password is incorrect",
			})
			return
		}

		ginContext.JSON(http.StatusInternalServerError, gin.H{
			"error": "internal server error",
		})
		return
	}

	ginContext.JSON(http.StatusOK, staff)
}
