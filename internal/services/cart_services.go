package services

import (
	"errors"
	"gorm.io/gorm"
	"marketplace-api/internal/models"
	"marketplace-api/internal/repository"
)

type CartService struct {
	cartRepository        *repository.CartRepository
	productRepository     *repository.ProductRepository
	distributorRepository *repository.DistributorRepository
}

func NewCartService(cartRepository *repository.CartRepository, productRepository *repository.ProductRepository, distributorRepository *repository.DistributorRepository) *CartService {
	return &CartService{cartRepository: cartRepository, productRepository: productRepository, distributorRepository: distributorRepository}
}

func (cs *CartService) AddCartItem(storeID int64, product *models.Product, quantity int64) error {
	cart, err := cs.cartRepository.GetCartByStoreID(storeID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			cart = &models.Cart{
				StoreID:    storeID,
				TotalPrice: 0,
			}
			if err := cs.cartRepository.CreateCart(cart); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	cartItem := &models.CartItem{
		CartID:    cart.ID,
		ProductID: product.ID,
		Quantity:  quantity,
	}

	oldCartItem, err := cs.cartRepository.GetCartItemByProductID(product.ID)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			cart.TotalPrice = cart.TotalPrice + product.Price*float64(quantity)
			err = cs.cartRepository.UpdateCart(cart)
			if err != nil {
				return err
			}
			return cs.cartRepository.AddCartItem(cartItem)
		} else {
			return err
		}
	}
	cart.TotalPrice = cart.TotalPrice - (product.Price * float64(oldCartItem.Quantity)) + product.Price*float64(quantity)
	err = cs.cartRepository.UpdateCart(cart)
	if err != nil {
		return err
	}
	return cs.cartRepository.UpdateCartItem(cartItem)

}

func (cs *CartService) UpdateCartItem(cartItem *models.CartItem) error {
	return cs.cartRepository.UpdateCartItem(cartItem)
}

func (cs *CartService) DeleteCartItem(storeID, productID int64) error {
	cart, err := cs.cartRepository.GetCartByStoreID(storeID)
	if err != nil {
		return err
	}

	cartItem, err := cs.cartRepository.GetCartItemByProductID(productID)
	if err != nil {
		return err
	}
	product, err := cs.productRepository.GetProductByID(productID)
	if err != nil {
		return err
	}
	err = cs.cartRepository.RemoveCartItem(cart.ID, productID)
	if err != nil {
		return err
	}

	cartItems, err := cs.cartRepository.GetCartItems(cart.ID)
	if err != nil {
		return err
	}
	if len(cartItems) == 0 {
		return cs.cartRepository.DeleteCart(storeID)
	}

	cart.TotalPrice = cart.TotalPrice - (product.Price * float64(cartItem.Quantity))
	err = cs.cartRepository.UpdateCart(cart)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CartService) GetCart(storeID int64) (*models.Cart, error) {
	cart, err := cs.cartRepository.GetCartByStoreID(storeID)
	if err != nil {
		return nil, err
	}
	cart.Items, err = cs.cartRepository.GetCartItems(cart.ID)
	if err != nil {
		return nil, err
	}
	for i, cartItem := range cart.Items {
		product, err := cs.productRepository.GetProductByID(cartItem.ProductID)
		if err != nil {
			return nil, err
		}
		distributor, err := cs.distributorRepository.GetDistributorByUserID(product.DistributorID)
		if err != nil {
			return nil, err
		}
		product.Distributor = *distributor
		cart.Items[i].Product = *product
	}
	return cart, nil
}

func (cs *CartService) DeleteCart(storeID int64) error {
	return cs.cartRepository.DeleteCart(storeID)
}
