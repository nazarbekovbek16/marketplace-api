//store_repository.go

package repository

import (
	"errors"
	"gorm.io/gorm"
	"marketplace-api/internal/models"
)

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) *StoreRepository {
	return &StoreRepository{db: db}
}

func (sr *StoreRepository) CreateStore(store *models.Store) error {
	if err := sr.db.Create(&store).Error; err != nil {
		return err
	}
	return nil
}

func (sr *StoreRepository) GetStoreByID(id int64) (*models.Store, error) {
	var store models.Store
	if err := sr.db.Where("id = ?", id).First(&store).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &store, nil
}

func (sr *StoreRepository) UpdateStore(storeID int64, updatedStore *models.Store) error {

	updatedStore.ID = storeID
	updatedStore.UserID = storeID
	if err := sr.db.Where("id = ?", storeID).Updates(&updatedStore).Error; err != nil {
		return err
	}
	if updatedStore.ImgUrl == "" {
		return sr.db.Model(&updatedStore).Where("id = ?", storeID).Update("img_url", "").Error
	}
	return nil

}

func (sr *StoreRepository) DeleteStore(storeID int64) error {
	store, err := sr.GetStoreByUserID(storeID)
	if err != nil {
		return err
	}

	if err = sr.db.Delete(store).Error; err != nil {
		return err
	}
	return nil
}

func (sr *StoreRepository) GetStoreByUserID(userID int64) (*models.Store, error) {
	var store models.Store
	if err := sr.db.Where("user_id = ?", userID).First(&store).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &store, nil
}

func (sr *StoreRepository) GetEmail(userID int64) (string, error) {
	var user models.User
	err := sr.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Email, nil
}
