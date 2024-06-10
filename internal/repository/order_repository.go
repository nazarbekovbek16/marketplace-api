package repository

import (
	"gorm.io/gorm"
	"marketplace-api/internal/models"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (or *OrderRepository) CreateOrder(order *models.Order) error {
	return or.db.Create(&order).Error
}

func (or *OrderRepository) CreateOrderStage(stage *models.Stage) error {
	return or.db.Create(&stage).Error
}
func (or *OrderRepository) UpdateOrderStage(stage *models.Stage) error {
	return or.db.Updates(&stage).Error
}

func (or *OrderRepository) GetOrderByID(userID, orderID int64, role string) (*models.Order, error) {
	var order *models.Order
	if err := or.db.Where("id = ? AND "+role+"_id = ?", orderID, userID).First(&order).Error; err != nil {
		return nil, err
	}
	return order, nil
}

func (or *OrderRepository) GetStageByID(stageID int64) (*models.Stage, error) {
	var orderStages models.Stage
	if err := or.db.Where("id = ?", stageID).Find(&orderStages).Error; err != nil {
		return nil, err
	}
	return &orderStages, nil
}

func (or *OrderRepository) UpdateOrderStatus(order *models.Order) error {
	return or.db.Where("id = ?", order.ID).Updates(order).Error
}

func (or *OrderRepository) GetOrders(userID int64, role string) ([]models.Order, error) {
	var orders []models.Order
	if err := or.db.Where(role+"_id = ?", userID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (or *OrderRepository) GetSuccessOrders(userID int64, role string) ([]models.Order, error) {
	var orders []models.Order
	if err := or.db.Find(&orders, "stage_id IN (SELECT id FROM stages WHERE stage = 'success' AND status = 'success') AND orders."+role+"_id = ?", userID).Error; err != nil {
		return nil, err
	}
	return orders, nil
}
