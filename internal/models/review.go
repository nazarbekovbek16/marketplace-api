package models

// Review model info
type Review struct {
	ID            int64  `json:"id" gorm:"primaryKey"`
	DistributorId int64  `json:"distributor_id"`
	ProductId     int64  `json:"product_id"`
	StoreId       int64  `json:"store_id"`
	Rating        int    `json:"rating"`
	Text          string `json:"text"`
}

type ReviewInput struct {
	DistributorId int64  `json:"distributor_id"`
	ProductId     int64  `json:"product_id"`
	Rating        int    `json:"rating"`
	Text          string `json:"text"`
}
