package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"marketplace-api/internal/models"
	"marketplace-api/internal/services"
	validator "marketplace-api/internal/util"
	"net/http"
	"strconv"
	"time"
)

type AdminHandler struct {
	userService        *services.UserService
	distributorService *services.DistributorService
	storeService       *services.StoreService
	log                *logrus.Logger
}

func NewAdminHandler(userService *services.UserService, distributorService *services.DistributorService, storeService *services.StoreService, log *logrus.Logger) *AdminHandler {
	return &AdminHandler{userService: userService, distributorService: distributorService, storeService: storeService, log: log}
}

func (ah *AdminHandler) GetAllUsers(c *gin.Context) {
	var input struct {
		CompanyName string
		models.Filters
	}
	v := validator.New()

	input.CompanyName = c.DefaultQuery("company_name", "")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	input.Filters.Page = page
	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if err != nil {
		pageSize = 10
	}
	input.Filters.PageSize = pageSize
	input.Filters.Sort = c.DefaultQuery("sort", "id")

	input.Filters.SortSafelist = []string{"id", "created_at", "-id", "-id", "-created_at"}

	if models.ValidateFilters(v, input.Filters); !v.Valid() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	type User struct {
		ID          int64     `json:"id"`
		CompanyName string    `json:"company_name"`
		Name        string    `json:"name"`
		Email       string    `json:"email"`
		BIN         string    `json:"bin"`
		Role        string    `json:"role"`
		IsActive    bool      `json:"is_active"`
		Created     time.Time `json:"created"`
		ImgUrl      string    `json:"img_url"`
	}
	var users []User

	stores, distributors, metadata, err := ah.userService.ListUsers(input.CompanyName, input.Filters)

	for _, store := range stores {
		var usr User
		usr.Name = store.Name
		usr.CompanyName = store.CompanyName
		usr.Email = store.User.Email
		usr.Role = store.User.Role
		usr.BIN = store.BIN
		usr.IsActive = ah.userService.GetStatusById(store.UserID)
		usr.ID = store.ID
		usr.ImgUrl = store.ImgUrl
		usr.Created = store.User.CreatedAt
		users = append(users, usr)
	}

	for _, distributor := range distributors {
		var usr User
		usr.Name = distributor.Name
		usr.CompanyName = distributor.CompanyName
		usr.Email = distributor.User.Email
		usr.Role = distributor.User.Role
		usr.IsActive = ah.userService.GetStatusById(distributor.UserID)
		usr.ID = distributor.ID
		usr.BIN = distributor.BIN
		usr.Created = distributor.User.CreatedAt
		usr.ImgUrl = distributor.ImgUrl
		users = append(users, usr)
	}
	if users == nil {
		c.JSON(http.StatusOK, gin.H{"users": []User{}, "metadata": metadata})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users, "metadata": metadata})
}

func (ah *AdminHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	err = ah.userService.DeleteUser(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user doesnt exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user with id %d successfully deleted", id)})
}

func (ah *AdminHandler) ActivateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	err = ah.userService.UpdateStatus(id, true)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user doesnt exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user with id %d successfully change status", id)})
}

func (ah *AdminHandler) DeactivateUser(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}

	err = ah.userService.UpdateStatus(id, false)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "user doesnt exist"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user with id %d successfully change status", id)})
}
