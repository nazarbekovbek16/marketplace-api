package services

import (
	"errors"
	"marketplace-api/internal/models"
	"marketplace-api/internal/repository"
	"math"
	"time"
)

type OrderService struct {
	orderRepository   *repository.OrderRepository
	productRepository *repository.ProductRepository
}

func NewOrderService(orderRepository *repository.OrderRepository, productRepository *repository.ProductRepository) *OrderService {
	return &OrderService{orderRepository: orderRepository, productRepository: productRepository}
}

func (os *OrderService) CreatOrder(cart *models.Cart, city, address string) error {
	for _, cartItem := range cart.Items {
		product, err := os.productRepository.GetProductByID(cartItem.ProductID)
		if err != nil {
			return err
		}
		if product.Stock < cartItem.Quantity {
			return errors.New("not enough quantity in stock for product " + product.ProductName)
		}
	}
	for _, cartItem := range cart.Items {
		product, err := os.productRepository.GetProductByID(cartItem.ProductID)
		if err != nil {
			return err
		}
		product.Stock = product.Stock - cartItem.Quantity
		err = os.productRepository.UpdateProduct(product)
		if err != nil {
			return err
		}
	}

	storeEmail, err := os.productRepository.GetEmail(cart.StoreID)
	if err != nil {
		return err
	}
	for _, cartItem := range cart.Items {
		distributorEmail, err := os.productRepository.GetEmail(cartItem.Product.DistributorID)
		if err != nil {
			return err
		}
		order := &models.Order{
			StoreID:          cart.StoreID,
			ProductID:        cartItem.ProductID,
			Product:          cartItem.Product,
			Quantity:         cartItem.Quantity,
			TotalPrice:       float64(cartItem.Quantity) * cartItem.Product.Price,
			Timestamp:        time.Now(),
			Distributor:      cartItem.Product.Distributor,
			DistributorID:    cartItem.Product.DistributorID,
			Status:           models.OrderStatusActive,
			City:             city,
			Address:          address,
			StoreEmail:       storeEmail,
			DistributorEmail: distributorEmail,
		}
		stage := &models.Stage{
			Stage:  models.StageNew,
			Status: models.StageSuccess,
		}
		err = os.orderRepository.CreateOrderStage(stage)
		order.Stage = *stage
		order.StageID = stage.ID
		err = os.orderRepository.CreateOrder(order)
		if err != nil {
			return err
		}
	}
	return nil
}

func (os *OrderService) ChangeOrderStatus(order models.Order, stageStatus string) error {
	stage, err := os.orderRepository.GetStageByID(order.StageID)
	if err != nil {
		return err
	}
	switch stage.Status {
	case models.StageStatusSuccess:
		if stage.Stage == models.StageProcessing || stage.Stage == models.StageShipped {
			if stageStatus == models.StageStatusWarning {
				stage.Stage = models.StageShipped
				stage.Status = stageStatus
				err = os.orderRepository.UpdateOrderStage(stage)
				return nil
			} else if stageStatus == models.StageStatusError {
				stage.Stage = stage.GetNextStage()
				stage.Status = stageStatus
				err = os.orderRepository.UpdateOrderStage(stage)
				if err != nil {
					return err
				}
				if stageStatus == models.StageStatusError {
					order.Status = models.OrderStatusClosed
					return os.orderRepository.UpdateOrderStatus(&order)
				}
			}
			stage.Stage = models.StageSuccess
			stage.Status = stageStatus
			err = os.orderRepository.UpdateOrderStage(stage)
			order.Status = models.OrderStatusClosed
			return os.orderRepository.UpdateOrderStatus(&order)
		}
		stage.Stage = stage.GetNextStage()
		stage.Status = stageStatus
		err = os.orderRepository.UpdateOrderStage(stage)
		if err != nil {
			return err
		}
		if stageStatus == models.StageStatusError {

			order.Status = models.OrderStatusClosed
			return os.orderRepository.UpdateOrderStatus(&order)
		}
		return nil
	case models.StageStatusWarning:
		if stageStatus == models.StageStatusError {
			stage.Stage = stage.GetNextStage()
			stage.Status = stageStatus
			err = os.orderRepository.UpdateOrderStage(stage)
			if err != nil {
				return err
			}
			order.Status = models.OrderStatusClosed
			return os.orderRepository.UpdateOrderStatus(&order)
		}
		if stage.Stage == models.StageShipped {
			stage.Stage = models.StageSuccess
			stage.Status = stageStatus
			err = os.orderRepository.UpdateOrderStage(stage)
			order.Status = models.OrderStatusClosed
			return os.orderRepository.UpdateOrderStatus(&order)
		}
		stage.Status = stageStatus
		err = os.orderRepository.UpdateOrderStage(stage)
		if err != nil {
			return err
		}
	}
	return nil
}

func (os *OrderService) GetOrderByID(userID, orderID int64, role string) (*models.Order, error) {
	order, err := os.orderRepository.GetOrderByID(userID, orderID, role)
	if err != nil {
		return nil, err
	}
	stage, err := os.orderRepository.GetStageByID(order.StageID)
	if err != nil {
		return nil, err
	}
	order.Stage = *stage
	order.TotalPrice = math.Round(order.TotalPrice*100) / 100
	return order, nil
}

func (os *OrderService) GetOrders(storeID int64, role string) ([]models.Order, error) {
	orders, err := os.orderRepository.GetOrders(storeID, role)
	if err != nil {
		return nil, err
	}
	for i, order := range orders {
		stage, err := os.orderRepository.GetStageByID(order.StageID)
		if err != nil {
			return nil, err
		}
		orders[i].Stage = *stage
	}
	return orders, nil
}

func (os *OrderService) GetSuccessOrders(storeID int64, role string) ([]models.Order, error) {
	orders, err := os.orderRepository.GetSuccessOrders(storeID, role)
	if err != nil {
		return nil, err
	}
	for i, order := range orders {
		stage, err := os.orderRepository.GetStageByID(order.StageID)
		if err != nil {
			return nil, err
		}
		orders[i].Stage = *stage
	}
	return orders, nil
}
