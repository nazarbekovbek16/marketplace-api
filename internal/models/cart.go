package models

// Cart model info
type Cart struct {
	ID         int64      `json:"id" gorm:"primaryKey"`
	StoreID    int64      `json:"store_id" gorm:"not null;unique"`
	Store      Store      `gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Items      []CartItem `json:"items" gorm:"foreignKey:CartID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TotalPrice float64    `json:"total_price"`
}

// CartItem model info
type CartItem struct {
	ID        int64   `json:"id" gorm:"primaryKey"`
	CartID    int64   `json:"cart_id"`
	ProductID int64   `json:"product_id"`
	Quantity  int64   `json:"quantity"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
