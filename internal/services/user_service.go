// internal/services/user_service.go

package services

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"marketplace-api/internal/models"
	"marketplace-api/internal/repository"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type UserService struct {
	userRepository        *repository.UserRepository
	distributorRepository *repository.DistributorRepository
	storeRepository       *repository.StoreRepository
}

func NewUserService(userRepository *repository.UserRepository, distributorRepository *repository.DistributorRepository, storeRepository *repository.StoreRepository) *UserService {
	return &UserService{
		userRepository:        userRepository,
		distributorRepository: distributorRepository,
		storeRepository:       storeRepository,
	}
}

func (us *UserService) RegisterUser(input *models.RegisterInput) error {
	user := &models.User{
		Email:    input.Email,
		Password: input.Password,
		Role:     input.Role,
	}
	err := us.userRepository.CreateUser(user)
	if err != nil {
		return err
	}

	switch user.Role {
	case models.RoleDistributor:
		distributor := &models.Distributor{
			ID:          user.ID,
			Name:        input.Name,
			CompanyName: input.CompanyName,
			PhoneNumber: input.PhoneNumber,
			City:        input.City,
			BIN:         input.BIN,
			UserID:      user.ID,
			User:        *user,
		}
		if err := us.distributorRepository.CreateDistributor(distributor); err != nil {
			return err
		}
	case models.RoleStore:
		store := &models.Store{
			ID:          user.ID,
			Name:        input.Name,
			CompanyName: input.CompanyName,
			PhoneNumber: input.PhoneNumber,
			City:        input.City,
			BIN:         input.BIN,
			UserID:      user.ID,
			User:        *user,
		}
		if err := us.storeRepository.CreateStore(store); err != nil {
			return err
		}
	}
	return nil
}

func (us *UserService) ValidateCredentials(credentials models.LoginCredentials) (*models.User, error) {
	// Validate user credentials (e.g., check username/password against database)
	user, err := us.userRepository.FindByEmail(credentials.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}
	// Compare hashed password
	if user.VerifyPassword(credentials.Password) != nil {
		fmt.Println(user.VerifyPassword(credentials.Password).Error())
		return nil, ErrInvalidCredentials
	}
	if user.Role != credentials.Role {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (us *UserService) ListUsers(email string, filters models.Filters) ([]*models.Store, []*models.Distributor, models.Metadata, error) {
	users, metadata, err := us.userRepository.ListUsers(email, filters)
	if err != nil {
		return nil, nil, models.Metadata{}, err
	}
	fmt.Println(users)
	var distributors []*models.Distributor
	for _, usr := range users {
		distributor, err := us.distributorRepository.GetDistributorByUserID(usr.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return nil, nil, models.Metadata{}, err
		}
		distributor.User = *usr
		distributors = append(distributors, distributor)
	}
	var stores []*models.Store
	for _, usr := range users {
		store, err := us.storeRepository.GetStoreByUserID(usr.ID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return nil, nil, models.Metadata{}, err
		}
		store.User = *usr
		stores = append(stores, store)
	}

	return stores, distributors, metadata, nil
}

func (us *UserService) DeleteUser(id int64) error {
	user, err := us.userRepository.FindByID(id)
	if err != nil {
		return err
	}
	if user.Role == models.RoleAdmin {
		return errors.New("cant delete admin")
	}
	err = us.userRepository.Delete(user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (us *UserService) UpdatePassword(id int64, currentPassword, newPassword string) error {
	user, err := us.userRepository.FindByID(id)
	if err != nil {
		return err
	}
	if user.VerifyPassword(currentPassword) != nil {
		return errors.New("invalid password")
	}

	user.Password = newPassword
	err = user.HashPassword()
	if err != nil {
		return err
	}
	if err := us.userRepository.UpdatePassword(user); err != nil {
		return err
	}
	return nil
}

func (us *UserService) UpdateEmail(id int64, email, password string) error {
	user, err := us.userRepository.FindByID(id)
	if err != nil {
		return err
	}
	if user.VerifyPassword(password) != nil {
		return errors.New("invalid password")
	}
	user.Email = email

	if err := us.userRepository.UpdateEmail(user); err != nil {
		return err
	}
	return nil
}

func (us *UserService) UpdateStatus(id int64, status bool) error {
	if err := us.userRepository.UpdateStatus(id, status); err != nil {
		return err
	}
	return nil
}

func (us *UserService) GetStatusById(id int64) bool {
	return us.userRepository.GetStatusById(id)
}
