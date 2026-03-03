package entity

type CreateStaffRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	HospitalID string `json:"hospital_id" binding:"required"`
}

type Staff struct {
	ID         string  `json:"id"`
	Username   string  `json:"username"`
	Password   string  `json:"password"`
	HospitalID string  `json:"hospital_id"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  *string `json:"updated_at"`
}
