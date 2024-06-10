package models

import (
	"github.com/lib/pq"
)

// Product model info
type Product struct {
	ID                 int64          `json:"id" gorm:"primaryKey"`
	ProductName        string         `json:"product_name"`
	ProductDescription string         `json:"product_description"`
	Price              float64        `json:"price"`
	ImgURLs            pq.StringArray `json:"ImgURLs" gorm:"type:text[]"`
	MinimumQuantity    int64          `json:"minimum_quantity"`
	DistributorID      int64          `gorm:"not null;" json:"distributor_id"`
	Distributor        Distributor    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"distributor"`
	Stock              int64          `json:"stock"`
	City               string         `json:"city"`
	Category           string         `json:"category"`
}

/*
8) некоторые отмененные заказы имеют статус актив
*/
