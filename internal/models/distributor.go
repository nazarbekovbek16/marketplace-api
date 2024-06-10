package models

// Distributor model info
type Distributor struct {
	ID          int64  `json:"id" gorm:"primaryKey"`
	Name        string `json:"name"`
	CompanyName string `json:"company_name"`
	Details     string `json:"details"`
	PhoneNumber string `json:"phone_number"`
	City        string `json:"city"`
	BIN         string `json:"bin"`
	ImgUrl      string `json:"img_url"`
	UserID      int64  `gorm:"not null;" json:"user_id"`
	User        User   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}
