// distributor_service.go

package services

import (
	"marketplace-api/internal/models"
	"marketplace-api/internal/repository"
)

type DistributorService struct {
	distributorRepository *repository.DistributorRepository
	userRepository        *repository.UserRepository
}

func NewDistributorService(distributorRepository *repository.DistributorRepository, userRepository *repository.UserRepository) *DistributorService {
	return &DistributorService{distributorRepository: distributorRepository, userRepository: userRepository}
}

func (ds *DistributorService) CreateDistributor(distributor *models.Distributor) error {
	return ds.distributorRepository.CreateDistributor(distributor)
}

func (ds *DistributorService) GetDistributorByID(id int64) (*models.Distributor, error) {
	return ds.distributorRepository.GetDistributorByID(id)
}

func (ds *DistributorService) GetStoreByID(id int64) (*models.Store, error) {
	return ds.distributorRepository.GetStoreByID(id)
}

func (ds *DistributorService) UpdateDistributor(distributorID int64, updatedDistributor *models.Distributor) error {
	return ds.distributorRepository.UpdateDistributor(distributorID, updatedDistributor)
}

func (ds *DistributorService) DeleteDistributor(distributorID int64) error {
	return ds.distributorRepository.DeleteDistributor(distributorID)
}

func (ds *DistributorService) GetDistributorByUserID(userID int64) (*models.Distributor, error) {
	user, err := ds.userRepository.FindByID(userID)
	if err != nil {
		return nil, err
	}
	distributor, err := ds.distributorRepository.GetDistributorByUserID(userID)
	if err != nil {
		return nil, err
	}
	distributor.User = *user
	return distributor, nil
}
