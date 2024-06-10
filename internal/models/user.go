package models

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

const (
	RoleAdmin       = "admin"
	RoleStore       = "store"
	RoleDistributor = "distributor"
)

// User model info
type User struct {
	ID       int64  `gorm:"primary_key" json:"id"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `gorm:"not null" json:"-"`
	Role     string `gorm:"not null" json:"role"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StatusUser struct {
	ID     int64 `gorm:"primary_key" json:"id"`
	UserId int64 `gorm:"not null" json:"user_id"`
	Status bool  `gorm:"not null" json:"status"`
}

type StatusUserInput struct {
	Status string `json:"status" binding:"required"`
}

// LoginCredentials model info
type LoginCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// RegisterInput model info
type RegisterInput struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Role        string `json:"role" binding:"required"`
	Name        string `json:"name"`
	CompanyName string `json:"company_name"`
	PhoneNumber string `json:"phone_number"`
	City        string `json:"city"`
	BIN         string `json:"bin"`
}

func (u *User) BeforeSave(*gorm.DB) error {
	if err := u.HashPassword(); err != nil {
		return err
	}
	u.SetCreatedAt()
	return nil
}
func (u *User) BeforeUpdate(db *gorm.DB) error {
	u.SetUpdatedAt()
	return nil
}

// LoginRequest model info
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

// HashPassword hashes the user's password using bcrypt
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// VerifyPassword verifies if the provided password matches the hashed password
func (u *User) VerifyPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// Validate checks if the user object is valid
func (u *User) Validate() error {
	if u.Email == "" {
		return errors.New("email is required")
	}
	// Add more validation rules as needed
	return nil
}

// SetCreatedAt sets the CreatedAt field to the current time
func (u *User) SetCreatedAt() {
	u.CreatedAt = time.Now()
}

// SetUpdatedAt sets the UpdatedAt field to the current time
func (u *User) SetUpdatedAt() {
	u.UpdatedAt = time.Now()
}
