package repository

import (
	"gorm.io/gorm"
	"marketplace-api/internal/models"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (cr *CartRepository) CreateCart(cart *models.Cart) error {
	return cr.db.Create(cart).Error
}

func (cr *CartRepository) GetCartByStoreID(storeID int64) (*models.Cart, error) {
	var cart models.Cart
	if err := cr.db.Where("store_id = ?", storeID).First(&cart).Error; err != nil {
		return nil, err
	}
	return &cart, nil
}

func (cr *CartRepository) GetCartItemByProductID(productId int64) (*models.CartItem, error) {
	var cartItem models.CartItem
	if err := cr.db.Where("product_id = ?", productId).First(&cartItem).Error; err != nil {
		return nil, err
	}
	return &cartItem, nil
}

func (cr *CartRepository) AddCartItem(cartItem *models.CartItem) error {
	return cr.db.Create(cartItem).Error
}

func (cr *CartRepository) UpdateCartItem(cartItem *models.CartItem) error {
	return cr.db.Where("product_id = ?", cartItem.ProductID).Updates(cartItem).Error
}

func (cr *CartRepository) UpdateCart(cart *models.Cart) error {
	return cr.db.Where("id = ?", cart.ID).Updates(&cart).Error
}

func (cr *CartRepository) RemoveCartItem(cartID int64, productID int64) error {
	return cr.db.Where("product_id = ? AND cart_id = ?", productID, cartID).Delete(&models.CartItem{}).Error
}

func (cr *CartRepository) DeleteCart(storeID int64) error {
	cart, err := cr.GetCartByStoreID(storeID)
	if err != nil {
		return err
	}

	err = cr.db.Where("cart_id = ? ", cart.ID).Delete(&models.CartItem{}).Error
	if err != nil {
		return err
	}
	err = cr.db.Where("id = ? ", cart.ID).Delete(cart).Error
	if err != nil {
		return err
	}
	return nil
}

func (cr *CartRepository) GetCartItems(cartID int64) ([]models.CartItem, error) {
	var cartItems []models.CartItem
	if err := cr.db.Where("cart_id = ?", cartID).Find(&cartItems).Error; err != nil {
		return nil, err
	}
	return cartItems, nil
}
