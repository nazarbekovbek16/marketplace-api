package services

import (
	"fmt"
	"marketplace-api/internal/models"
	"marketplace-api/internal/repository"
)

type StoreService struct {
	storeRepository *repository.StoreRepository
	userRepository  *repository.UserRepository
}

func NewStoreService(storeRepository *repository.StoreRepository, userRepository *repository.UserRepository) *StoreService {
	return &StoreService{storeRepository: storeRepository, userRepository: userRepository}
}

func (ss *StoreService) CreateStore(store *models.Store) error {
	return ss.storeRepository.CreateStore(store)
}

func (ss *StoreService) GetStoreByID(id int64) (*models.Store, error) {
	return ss.storeRepository.GetStoreByID(id)
}

func (ss *StoreService) UpdateStore(storeID int64, updatedStore *models.Store) error {
	return ss.storeRepository.UpdateStore(storeID, updatedStore)
}

func (ss *StoreService) DeleteDistributor(storeID int64) error {
	return ss.storeRepository.DeleteStore(storeID)
}

func (ss *StoreService) GetStoreByUserID(userID int64) (*models.Store, error) {
	user, err := ss.userRepository.FindByID(userID)
	if err != nil {
		return nil, err
	}
	fmt.Println(user.Email)
	store, err := ss.storeRepository.GetStoreByUserID(userID)
	if err != nil {
		return nil, err
	}
	store.User = *user
	return store, nil
}
