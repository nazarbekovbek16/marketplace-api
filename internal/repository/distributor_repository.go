// distributor_repository.go

package repository

import (
	"errors"
	"gorm.io/gorm"
	"marketplace-api/internal/models"
)

type DistributorRepository struct {
	db *gorm.DB
}

func NewDistributorRepository(db *gorm.DB) *DistributorRepository {
	return &DistributorRepository{db: db}
}

func (sr *DistributorRepository) CreateDistributor(distributor *models.Distributor) error {
	if err := sr.db.Create(distributor).Error; err != nil {
		return err
	}
	return nil
}

func (sr *DistributorRepository) GetDistributorByID(id int64) (*models.Distributor, error) {
	var distributor models.Distributor
	if err := sr.db.First(&distributor, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &distributor, nil
}

func (sr *DistributorRepository) GetStoreByID(id int64) (*models.Store, error) {
	var store models.Store
	if err := sr.db.First(&store, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &store, nil
}

func (sr *DistributorRepository) UpdateDistributor(distributorID int64, updatedDistributor *models.Distributor) error {
	updatedDistributor.ID = distributorID
	updatedDistributor.UserID = distributorID
	if err := sr.db.Where("id = ?", distributorID).Updates(&updatedDistributor).Error; err != nil {
		return err
	}
	if updatedDistributor.ImgUrl == "" {
		return sr.db.Model(&updatedDistributor).Where("id = ?", distributorID).Update("img_url", "").Error
	}
	return nil
}

func (sr *DistributorRepository) DeleteDistributor(distributorID int64) error {
	distributor, err := sr.GetDistributorByID(distributorID)
	if err != nil {
		return err
	}

	if err := sr.db.Delete(distributor).Error; err != nil {
		return err
	}
	return nil
}

func (sr *DistributorRepository) GetDistributorByUserID(userID int64) (*models.Distributor, error) {
	var distributor models.Distributor
	if err := sr.db.Where("user_id = ?", userID).First(&distributor).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &distributor, nil
}

func (sr *DistributorRepository) GetEmail(userID int64) (string, error) {
	var user models.User
	err := sr.db.Where("id = ?", userID).First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Email, nil
}
