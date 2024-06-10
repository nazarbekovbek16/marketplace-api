package repository

import (
	"database/sql"
	"errors"
	"gorm.io/gorm"
	"marketplace-api/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser inserts a new user record into the database
func (ur *UserRepository) CreateUser(user *models.User) error {
	if err := ur.db.Create(&user).Error; err != nil {
		return err
	}

	if err := ur.db.Create(&models.StatusUser{
		UserId: user.ID,
		Status: false,
	}).Error; err != nil {
		return err
	}
	return nil
}

// FindByID retrieves a user record from the database by ID
func (ur *UserRepository) FindByID(id int64) (*models.User, error) {
	var user models.User
	if err := ur.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

func (ur *UserRepository) GetStatusById(id int64) bool {
	var statusUser models.StatusUser
	ur.db.Where("user_id = ?", id).First(&statusUser)
	return statusUser.Status
}

func (ur *UserRepository) UpdateStatus(id int64, status bool) error {
	if err := ur.db.Model(&models.StatusUser{}).Where("user_id = ?", id).UpdateColumn("status", status).Error; err != nil {
		return err
	}
	return nil
}

// FindByEmail retrieves a user record from the database by email
func (ur *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := ur.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}

// UpdateEmail updates an existing user email record in the database
func (ur *UserRepository) UpdateEmail(user *models.User) error {
	if err := ur.db.Model(&user).UpdateColumn("email", user.Email).Error; err != nil {
		return err
	}
	return nil
}

// UpdatePassword updates an existing user record in the database
func (ur *UserRepository) UpdatePassword(user *models.User) error {
	if err := ur.db.Model(&user).UpdateColumn("password", user.Password).Error; err != nil {
		return err
	}
	return nil
}

// Delete removes a user record from the database
func (ur *UserRepository) Delete(id int64) error {
	if err := ur.db.Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) ListUsers(companyName string, filters models.Filters) ([]*models.User, models.Metadata, error) {
	rows, err := ur.db.Table("users").
		Select("count(*) OVER()",
			"users.id", "users.email", "users.role",
			"users.created_at").
		Order(filters.SortColumn() + " " + filters.SortDirection()).
		Order("id ASC").
		Limit(filters.Limit()).
		Limit(filters.Offset()).Rows()
	if err != nil {
		return nil, models.Metadata{}, err
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	totalRecords := 0
	var users []*models.User

	for rows.Next() {
		var user models.User

		err := rows.Scan(
			&totalRecords,
			&user.ID,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, models.Metadata{}, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, models.Metadata{}, err
	}
	metadata := models.CalculateMetadata(totalRecords, filters.Page, filters.PageSize)
	return users, metadata, nil
}
