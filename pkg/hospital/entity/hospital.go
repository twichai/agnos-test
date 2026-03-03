package entity

type Hospital struct {
	ID          string `json:"id" gorm:"column:id;type:uuid;primaryKey"`
	Name        string `json:"name" gorm:"column:name"`
	HospitalURL string `json:"hospital_url" gorm:"column:hospital_url;unique"`
}

func (Hospital) TableName() string {
	return "hospitals"
}
