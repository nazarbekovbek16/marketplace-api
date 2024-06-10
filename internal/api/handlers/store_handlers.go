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

type StoreHandler struct {
	storeService       *services.StoreService
	productServices    *services.ProductService
	distributorService *services.DistributorService
	cartService        *services.CartService
	orderService       *services.OrderService
}

func NewStoreHandler(storeService *services.StoreService, productServices *services.ProductService, distributorService *services.DistributorService, cartService *services.CartService, orderService *services.OrderService) *StoreHandler {
	return &StoreHandler{storeService: storeService, productServices: productServices, distributorService: distributorService, cartService: cartService, orderService: orderService}
}

func (sh *StoreHandler) GetProfile(c *gin.Context) {
	store, err := sh.storeService.GetStoreByUserID(c.GetInt64("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"store": store, "email": store.User.Email})
}

func (sh *StoreHandler) UpdateProfile(c *gin.Context) {
	var input *models.Store
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}
	userId := c.GetInt64("user_id")

	err := sh.storeService.UpdateStore(userId, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Profile updated successfully"})

}

func (sh *StoreHandler) GetProduct(c *gin.Context) {
	productId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || productId < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
	}

	product, err := sh.productServices.GetProductByID(productId)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusNotFound, gin.H{"message": "the requested resource could not be found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	distributor, err := sh.distributorService.GetDistributorByUserID(product.DistributorID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	product.Distributor = *distributor

	c.JSON(http.StatusOK, gin.H{"product": product})
}

func (sh *StoreHandler) ListProducts(c *gin.Context) {
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
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": v.Errors})
		return
	}

	products, metadata, err := sh.productServices.GetProducts(input.ProductName, input.Filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	for i, product := range products {
		p, err := sh.productServices.GetProductByID(product.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		products[i] = p
	}

	c.JSON(http.StatusOK, gin.H{"products": products, "metadata": metadata})
}

func (sh *StoreHandler) AddToCart(c *gin.Context) {
	productId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || productId < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
	}

	var input struct {
		Quantity int64 `json:"quantity"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	product, err := sh.productServices.GetProductByID(productId)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusNotFound, gin.H{"message": "the requested resource could not be found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if product.Stock < input.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not enough quantity in stock"})
		return
	}
	if product.MinimumQuantity > input.Quantity || input.Quantity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "less than minimum quantity"})
		return
	}
	storeID := c.GetInt64("user_id")
	err = sh.cartService.AddCartItem(storeID, product, input.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "product added to cart successfully"})
}

func (sh *StoreHandler) UpdateCartItem(c *gin.Context) {
	productId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || productId < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
	}

	var input struct {
		Quantity int64 `json:"quantity"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	product, err := sh.productServices.GetProductByID(productId)
	fmt.Println(productId)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusNotFound, gin.H{"message": "the requested resource could not be found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	storeID := c.GetInt64("user_id")
	err = sh.cartService.AddCartItem(storeID, product, input.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "product's quantity in cart updated successfully"})
}

func (sh *StoreHandler) RemoveFromCart(c *gin.Context) {
	productId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || productId < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
		return
	}
	storeID := c.GetInt64("user_id")

	err = sh.cartService.DeleteCartItem(storeID, productId)
	if err != nil {
		if errors.Is(gorm.ErrRecordNotFound, err) {
			c.JSON(http.StatusNotFound, gin.H{"message": "the requested resource could not be found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product removed from cart successfully"})
}

func (sh *StoreHandler) GetCart(c *gin.Context) {
	storeID := c.GetInt64("user_id")

	cart, err := sh.cartService.GetCart(storeID)
	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusOK, gin.H{"cart": models.Cart{
				Items: []models.CartItem{},
			}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	cart.TotalPrice = math.Round(cart.TotalPrice*100) / 100
	c.JSON(http.StatusOK, gin.H{"cart": cart})
}

func (sh *StoreHandler) CreateOrder(c *gin.Context) {
	storeID := c.GetInt64("user_id")

	var input struct {
		City    string `json:"city"`
		Address string `json:"address"`
	}
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	cart, err := sh.cartService.GetCart(storeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "cart is empty"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = sh.orderService.CreatOrder(cart, input.City, input.Address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = sh.cartService.DeleteCart(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "order created successfully"})
}

func (sh *StoreHandler) CancelOrder(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || orderID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
	}
	storeID := c.GetInt64("user_id")
	order, err := sh.orderService.GetOrderByID(storeID, orderID, "store")
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

	err = sh.orderService.ChangeOrderStatus(*order, models.StageStatusError)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	product, err := sh.productServices.GetProductByID(order.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	product.Stock += order.Quantity

	err = sh.productServices.UpdateProduct(product)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order successfully canceled"})
}

func (sh *StoreHandler) GetOrder(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || orderID < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
	}
	storeID := c.GetInt64("user_id")
	order, err := sh.orderService.GetOrderByID(storeID, orderID, "store")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "order not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	product, err := sh.productServices.GetProductByID(order.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	order.Product = *product
	order.TotalPrice = math.Round(order.TotalPrice*100) / 100
	c.JSON(http.StatusOK, gin.H{"order": order})
}

func (sh *StoreHandler) ListOrders(c *gin.Context) {
	storeID := c.GetInt64("user_id")
	orders, err := sh.orderService.GetOrders(storeID, "store")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "order not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	for i, order := range orders {
		product, err := sh.productServices.GetProductByID(order.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		orders[i].Product = *product
		order.Product = *product
		orders[i].TotalPrice = math.Round(order.TotalPrice*100) / 100
	}
	c.JSON(http.StatusOK, gin.H{"orders": orders})
}

func (sh *StoreHandler) GetStatistics(c *gin.Context) {
	storeID := c.GetInt64("user_id")
	orders, err := sh.orderService.GetSuccessOrders(storeID, "store")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, "order not found")
			return
		}
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	var spentOverall float64
	var spentInMonth float64
	for i, order := range orders {
		spentOverall = spentOverall + order.TotalPrice
		if order.Timestamp.Year() == time.Now().Year() && time.Now().Month() == order.Timestamp.Month() {
			spentInMonth = spentInMonth + order.TotalPrice
		}
		product, err := sh.productServices.GetProductByID(order.ProductID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		orders[i].Product = *product
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders, "spent_overall": spentOverall, "spent_in_month": spentInMonth})
}

func (sh *StoreHandler) CreateReview(c *gin.Context) {
	storeID := c.GetInt64("user_id")

	var input models.ReviewInput
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	review := &models.Review{
		DistributorId: input.DistributorId,
		ProductId:     input.ProductId,
		StoreId:       storeID,
		Rating:        input.Rating,
		Text:          input.Text,
	}

	err := sh.productServices.CreatReview(review)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "review created successfully"})
}

func (sh *StoreHandler) GetReview(c *gin.Context) {
	storeID := c.GetInt64("user_id")

	reviews, err := sh.productServices.GetReviewByStoreId(storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"reviews": reviews})
}
func (sh *StoreHandler) GetReviewByProductId(c *gin.Context) {
	productId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || productId < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
	}

	reviews, err := sh.productServices.GetReviewsByProductId(productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"reviews": reviews})
}

func (sh *StoreHandler) DeleteReview(c *gin.Context) {
	reviewId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || reviewId < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id parameter"})
	}

	err = sh.productServices.DeleteByReviewId(reviewId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "order successfully deleted."})
}
