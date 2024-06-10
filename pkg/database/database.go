package database

import (
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"marketplace-api/internal/config"
	"marketplace-api/internal/models"
)

// InitDB initializes the database connection
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dbURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName, cfg.DBPassword)

	db, err := gorm.Open(postgres.Open(dbURI), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate creates missing tables based on the provided models
	err = db.AutoMigrate(
		&models.User{},
		&models.Distributor{},
		&models.Store{},
		&models.Product{},
		&models.Cart{},
		&models.CartItem{},
		&models.Order{},
		&models.Stage{},
		&models.StatusUser{},
		&models.Review{},
	)
	if err != nil {
		return nil, errors.New("failed to start database " + err.Error())
	}

	err = creatAdmin(cfg.AdminEmail, cfg.AdminPassword, db)
	if err != nil {
		return nil, errors.New("failed to create admin user " + err.Error())
	}

	//err = createData(db)
	return db, nil
}

func creatAdmin(adminEmail, adminPassword string, db *gorm.DB) error {
	var admin models.User
	result := db.Where("email = ?", adminEmail).First(&admin)
	fmt.Println("password", adminPassword)
	if result.RowsAffected == 0 {
		admin := models.User{
			Email:    adminEmail,
			Password: adminPassword,
			Role:     "admin",
		}
		err := db.Create(&admin).Error
		if err != nil {
			return err
		}
	}
	return nil
}

//func createData(db *gorm.DB) error {
//	user:=models.User{
//		ID:        2,
//		Email:     "bek@mail.ru",
//		Password:  "QWERTY123",
//		Role:      "store",
//	}
//	if err := db.Create(&user).Error; err != nil {
//		return err
//	}
//
//	if err := db.Create(&models.StatusUser{
//		UserId: user.ID,
//		Status: true,
//	}).Error; err != nil {
//		return err
//	}
//
//	if err := db.Create(&store).Error; err != nil {
//		return err
//	}
//
//	return nil
//}
