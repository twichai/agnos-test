package entity

import (
	"time"
	"twichai/agnos-test/pkg/hospital/entity"
)

type CreateStaffRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	HospitalID string `json:"hospital_id" binding:"required,uuid"`
}

type Staff struct {
	ID         string    `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	Username   string    `json:"username" gorm:"column:username"`
	Password   string    `json:"password_hash" gorm:"column:password_hash"`
	HospitalID string    `json:"hospital_id" gorm:"column:hospital_id;type:uuid"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`

	// relationship
	Hospital entity.Hospital `json:"hospital" gorm:"foreignKey:HospitalID;references:ID"`
}

type StaffLoginRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	HospitalID string `json:"hospital_id" binding:"required,uuid"`
}

func (Staff) TableName() string {
	return "staffs"
}
