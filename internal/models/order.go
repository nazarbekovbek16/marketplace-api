package models

import "time"

const (
	OrderStatusActive  = "active"
	OrderStatusClosed  = "closed"
	StageNew           = "new"
	StageConfirmed     = "confirmed"
	StageProcessing    = "processing"
	StageShipped       = "shipped"
	StageSuccess       = "success"
	StageStatusSuccess = "success"
	StageStatusWarning = "warning"
	StageStatusError   = "error"
)

// Order model info
type Order struct {
	ID               int64       `json:"id" gorm:"primaryKey"`
	StoreID          int64       `json:"store_id" gorm:"not null"`
	Store            Store       `gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ProductID        int64       `json:"product_id" gorm:"not null"`
	Product          Product     `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"product"`
	Quantity         int64       `json:"quantity"`
	TotalPrice       float64     `json:"total_price"`
	Timestamp        time.Time   `json:"timestamp"`
	DistributorID    int64       `json:"distributor_id"`
	Distributor      Distributor `gorm:"foreignKey:DistributorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Status           string      `json:"status"`
	StageID          int64       `json:"order_id"`
	Stage            Stage       `gorm:"foreignKey:StageID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"stage"`
	City             string      `json:"city"`
	Address          string      `json:"address"`
	StoreEmail       string      `json:"store_email"`
	DistributorEmail string      `json:"distributor_email"`
}

// Stage model info
type Stage struct {
	ID     int64  `json:"id" gorm:"primaryKey"`
	Stage  string `json:"stage"`
	Status string `json:"status"`
}

func (s *Stage) GetNextStage() string {
	switch s.Stage {
	case StageNew:
		return StageConfirmed
	case StageConfirmed:
		return StageProcessing
	case StageProcessing:
		return StageShipped
	case StageShipped:
		return StageSuccess
	case StageSuccess:
		return ""
	default:
		return ""
	}
}

func (s *Stage) GetPrevStage(stage string) string {
	switch stage {
	case StageNew:
		return ""
	case StageConfirmed:
		return StageNew
	case StageProcessing:
		return StageConfirmed
	case StageShipped:
		return StageProcessing
	case StageSuccess:
		return StageShipped
	default:
		return ""
	}
}
