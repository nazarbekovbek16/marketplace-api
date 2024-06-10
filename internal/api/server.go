package api

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"marketplace-api/internal/api/handlers"
	"marketplace-api/internal/api/middleware"
	"marketplace-api/internal/api/routes"
	"marketplace-api/internal/config"
	"marketplace-api/internal/repository"
	"marketplace-api/internal/services"
)

type Server struct {
	config *config.Config
	router *gin.Engine
	db     *gorm.DB
	logger *logrus.Logger
}

func NewServer(router *gin.Engine, db *gorm.DB, logger *logrus.Logger, config *config.Config) *Server {
	server := &Server{
		router: router,
		db:     db,
		config: config,
		logger: logger,
	}
	// Initialize repository layer
	userRepository := repository.NewUserRepository(db)
	distributorRepository := repository.NewDistributorRepository(db)
	storeRepository := repository.NewStoreRepository(db)
	productRepository := repository.NewProductRepository(db)
	cartRepository := repository.NewCartRepository(db)
	orderRepository := repository.NewOrderRepository(db)
	// Initialize service layer
	userService := services.NewUserService(userRepository, distributorRepository, storeRepository)
	distributorService := services.NewDistributorService(distributorRepository, userRepository)
	productService := services.NewProductService(productRepository, distributorRepository)
	storeService := services.NewStoreService(storeRepository, userRepository)
	cartService := services.NewCartService(cartRepository, productRepository, distributorRepository)
	orderService := services.NewOrderService(orderRepository, productRepository)
	// Initialize handler layer
	authHandler := handlers.NewAuthHandler(userService, distributorService, nil, config.JWTSecret, logger)
	distributorHandler := handlers.NewDistributorHandler(distributorService, productService, orderService)
	storeHandler := handlers.NewStoreHandler(storeService, productService, distributorService, cartService, orderService)
	adminHandler := handlers.NewAdminHandler(userService, distributorService, storeService, logger)
	//productHandler := handlers.ProductHandler{}
	// Register routes
	handler := handlers.NewHandlers(authHandler, distributorHandler, storeHandler, adminHandler)
	router.Use(middleware.CorsMiddleware())
	APIRouter := router.Group("/api")
	APIRouter.Static("images/", "./images/")
	routes.RegisterRoutes(APIRouter, *handler, config)
	return server
}
