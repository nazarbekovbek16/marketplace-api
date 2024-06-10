// distributor_handler.go

package handlers

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"marketplace-api/internal/models"
	"marketplace-api/internal/services"
	validator "marketplace-api/internal/util"
	"math"
	"net/http"
	"strconv"
	"time"
)

type DistributorHandler struct {
	distributorService *services.DistributorService
	productServices    *services.ProductService
	orderService       *services.OrderService
}

func NewDistributorHandler(distributorService *services.DistributorService, productServices *services.ProductService, orderService *services.OrderService) *DistributorHandler {
	return &DistributorHandler{distributorService: distributorService, productServices: productServices, orderService: orderService}
}

// GetProfile godoc
// @Summary      Get distributor profile
// @Description  Returns distributors profile
// @Tags         distributor
// @Security     BearerToken
// @Accept       json
// @Produce      json
// @Success      200  {object}  models.Distributor
// @Failure      400  {string}  Bad request
// @Failure      404  {string}  Not found
// @Failure      500  {string}  Internal server error
// @Router       /distributor/profile [get]
func (dh *DistributorHandler) GetProfile(c *gin.Context) {
	distributor, err := dh.distributorService.GetDistributorByUserID(c.GetInt64("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"distributor": distributor, "email": distributor.User.Email})
}

func (dh *DistributorHandler) UpdateProfile(c *gin.Context) {
	var input *models.Distributor
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	userId := c.GetInt64("user_id")

	fmt.Println(input.ImgUrl)

	err := dh.distributorService.UpdateDistributor(userId, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Profile updated successfully"})

}

func (dh *DistributorHandler) CreateProduct(c *gin.Context) {
	var input struct {
		ProductName        string   `json:"product_name"`
		ProductDescription string   `json:"product_description"`
		Price              float64  `json:"price"`
		ImgURLs            []string `json:"ImgURLs"`
		MinimumQuantity    int64    `json:"minimum_quantity"`
		Stock              int64    `json:"stock"`
		City               string   `json:"city"`
		Category           string   `json:"category"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	product := &models.Product{
		ProductName:        input.ProductName,
		ProductDescription: input.ProductDescription,
		Price:              input.Price,
		ImgURLs:            input.ImgURLs,
		MinimumQuantity:    input.MinimumQuantity,
		DistributorID:      c.GetInt64("user_id"),
		Stock:              input.Stock,
		City:               input.City,
		Category:           input.Category,
	}
	fmt.Println(product.ImgURLs)
	err := dh.productServices.CreateProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Product created successfully"})
}

func (dh *DistributorHandler) UpdateProduct(c *gin.Context) {
	productId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || productId < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
	}
	var input struct {
		ProductName        string   `json:"product_name"`
		ProductDescription string   `json:"product_description"`
		Price              float64  `json:"price"`
		ImgURLs            []string `json:"ImgURLs"`
		MinimumQuantity    int64    `json:"minimum_quantity"`
		Stock              int64    `json:"stock"`
		City               string   `json:"city"`
		Category           string   `json:"category"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	product, err := dh.productServices.GetProductByID(productId)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusNotFound, gin.H{"message": "the requested resource could not be found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if product.DistributorID != c.GetInt64("user_id") {
		c.JSON(http.StatusNotFound, gin.H{"message": "the requested resource could not be found"})
		return
	}

	product = &models.Product{
		ID:                 productId,
		ProductName:        input.ProductName,
		ProductDescription: input.ProductDescription,
		Price:              input.Price,
		ImgURLs:            input.ImgURLs,
		MinimumQuantity:    input.MinimumQuantity,
		DistributorID:      c.GetInt64("user_id"),
		Stock:              input.Stock,
		City:               input.City,
		Category:           input.Category,
	}
	fmt.Println(product.ImgURLs)

	err = dh.productServices.UpdateProduct(product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Product updated successfully"})
}

func (dh *DistributorHandler) GetProduct(c *gin.Context) {
	productId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || productId < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
	}

	product, err := dh.productServices.GetProductByID(productId)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusNotFound, gin.H{"message": "the requested resource could not be found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if product.DistributorID != c.GetInt64("user_id") {
		c.JSON(http.StatusNotFound, gin.H{"message": "the requested resource could not be found"})
		return
	}
	distributor, err := dh.distributorService.GetDistributorByUserID(product.DistributorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	product.Distributor = *distributor

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func (dh *DistributorHandler) ListProducts(c *gin.Context) {
	var input struct {
		ProductName string
		models.Filters
	}
	v := validator.New()

	input.ProductName = c.DefaultQuery("product_name", "")

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

	input.Filters.SortSafelist = []string{"id", "product_name", "price", "created_at", "-id", "-product_name", "-price", "-created_at"}

	if models.ValidateFilters(v, input.Filters); !v.Valid() {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	products, metadata, err := dh.productServices.GetProductsByDistributorID(input.ProductName, input.Filters, c.GetInt64("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": v.Errors})
		return
	}
	for i, product := range products {
		p, err := dh.productServices.GetProductByID(product.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		products[i] = p
	}
	c.JSON(http.StatusOK, gin.H{"products": products, "metadata": metadata})
}

func (dh *DistributorHandler) DeleteProduct(c *gin.Context) {
	productId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || productId < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
	}

	product, err := dh.productServices.GetProductByID(productId)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusNotFound, gin.H{"message": "the requested resource could not be found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if product.DistributorID != c.GetInt64("user_id") {
		c.JSON(http.StatusNotFound, gin.H{"message": "the requested resource could not be found"})
		return
	}

	err = dh.productServices.DeleteProduct(productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

func (dh *DistributorHandler) UpdateOrder(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || orderID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
	}
	var input struct {
		StageStatus string `json:"stage_status"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	distributorID := c.GetInt64("user_id")
	order, err := dh.orderService.GetOrderByID(distributorID, orderID, "distributor")
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusNotFound, "order doesnt exist")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if order.Status == models.OrderStatusClosed {
		c.JSON(http.StatusBadRequest, "order already canceled")
		return
	}
	err = dh.orderService.ChangeOrderStatus(*order, input.StageStatus)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order stage status changed"})
}

func (dh *DistributorHandler) GetOrder(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || orderID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
	}
	distributorID := c.GetInt64("user_id")
	order, err := dh.orderService.GetOrderByID(distributorID, orderID, "distributor")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "order not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	product, err := dh.productServices.GetProductByID(order.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	order.Product = *product
	order.TotalPrice = math.Round(order.TotalPrice*100) / 100
	c.JSON(http.StatusOK, gin.H{"order": order})
}

func (dh *DistributorHandler) ListOrders(c *gin.Context) {
	fmt.Println("here")
	distributorID := c.GetInt64("user_id")
	orders, err := dh.orderService.GetOrders(distributorID, "distributor")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "order not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	for i, order := range orders {
		product, err := dh.productServices.GetProductByID(order.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		orders[i].Product = *product
		orders[i].TotalPrice = math.Round(order.TotalPrice*100) / 100
	}
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (dh *DistributorHandler) GetStatistics(c *gin.Context) {
	distributorID := c.GetInt64("user_id")
	orders, err := dh.orderService.GetSuccessOrders(distributorID, "distributor")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "order not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var soldOverall float64
	var soldInMonth float64
	for i, order := range orders {
		soldOverall = soldOverall + order.TotalPrice
		if order.Timestamp.Year() == time.Now().Year() && time.Now().Month() == order.Timestamp.Month() {
			soldInMonth = soldInMonth + order.TotalPrice
		}
		product, err := dh.productServices.GetProductByID(order.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		orders[i].Product = *product
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders, "sold_overall": soldOverall, "sold_in_month": soldInMonth})
}

func (dh *DistributorHandler) GetReviews(c *gin.Context) {
	distributorID := c.GetInt64("user_id")

	reviews, err := dh.productServices.GetReviewsByDistributorId(distributorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func (dh *DistributorHandler) GetReviewByProductId(c *gin.Context) {
	productId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || productId < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
	}

	reviews, err := dh.productServices.GetReviewsByProductId(productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}
