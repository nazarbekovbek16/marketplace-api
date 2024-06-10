package routes

import (
	"github.com/gin-gonic/gin"
	"marketplace-api/internal/api/handlers"
	"marketplace-api/internal/api/middleware"
	"marketplace-api/internal/config"
	"marketplace-api/internal/models"
)

func RegisterRoutes(router *gin.RouterGroup, handlers handlers.Handlers, cfg *config.Config) {

	//Authentication routes
	authRouters := router.Group("/auth")
	authRouters.POST("/register", handlers.AuthHandler.Register)
	authRouters.POST("/login", handlers.AuthHandler.Login)
	authRouters.POST("/logout", handlers.AuthHandler.LogOut)

	router.GET("/store/user/:id", handlers.AuthHandler.GetStoreByID)
	router.GET("/distributor/user/:id", handlers.AuthHandler.GetDistributorByID)

	uploadRouters := router.Group("/upload")
	uploadRouters.Use(middleware.AuthMiddleware(cfg.JWTSecret, ""))
	uploadRouters.POST("/image", handlers.UploadImage)
	uploadRouters.POST("/images", handlers.UploadImages)

	//User routes
	userRouters := router.Group("/user")
	userRouters.Use(middleware.AuthMiddleware(cfg.JWTSecret, ""))
	//Change Email
	userRouters.PUT("/email", handlers.AuthHandler.ChangeEmail)
	//Change Password
	userRouters.PUT("/password", handlers.AuthHandler.ChangePassword)
	//Delete Account
	userRouters.DELETE("/", handlers.AuthHandler.DeleteAccount)

	//Admin
	adminRouters := router.Group("/admin")
	adminRouters.Use(middleware.AuthMiddleware(cfg.JWTSecret, "admin"))
	adminRouters.GET("/panel", handlers.AdminHandler.GetAllUsers)
	adminRouters.DELETE("/delete/user/:id", handlers.AdminHandler.DeleteUser)
	adminRouters.POST("/activate/user/:id", handlers.AdminHandler.ActivateUser)
	adminRouters.POST("/deactivate/user/:id", handlers.AdminHandler.DeactivateUser)

	//Distributors routes
	distributorRouters := router.Group("/distributor")
	distributorRouters.Use(
		middleware.AuthMiddleware(cfg.JWTSecret, "distributor"),
	)
	//Profile routes
	distributorRouters.GET("/profile", handlers.DistributorHandler.GetProfile)
	distributorRouters.PUT("/profile", handlers.DistributorHandler.UpdateProfile)
	//products routes
	distributorRouters.POST("/products", handlers.DistributorHandler.CreateProduct)
	distributorRouters.PUT("/products/:id", handlers.DistributorHandler.UpdateProduct)
	distributorRouters.GET("/products/:id", handlers.DistributorHandler.GetProduct)
	distributorRouters.GET("/products", handlers.DistributorHandler.ListProducts)
	distributorRouters.DELETE("/products/:id", handlers.DistributorHandler.DeleteProduct)
	//orders routes
	distributorRouters.PUT("/orders/:id", handlers.DistributorHandler.UpdateOrder)
	distributorRouters.GET("/orders/:id", handlers.DistributorHandler.GetOrder)
	distributorRouters.GET("/orders", handlers.DistributorHandler.ListOrders)
	distributorRouters.GET("/orders/sold", handlers.DistributorHandler.GetStatistics)
	//review routes
	distributorRouters.GET("/reviews", handlers.DistributorHandler.GetReviews)
	distributorRouters.GET("/reviews/product/:id", handlers.DistributorHandler.GetReviewByProductId)

	//Stores routes
	storeRouters := router.Group("/store")
	storeRouters.Use(middleware.AuthMiddleware(cfg.JWTSecret, models.RoleStore))

	//Profile routes
	storeRouters.GET("/profile", handlers.StoreHandler.GetProfile)
	storeRouters.PUT("/profile", handlers.StoreHandler.UpdateProfile)
	//products routes
	storeRouters.GET("/products/:id", handlers.StoreHandler.GetProduct)
	storeRouters.GET("/products", handlers.StoreHandler.ListProducts)
	//carts routes
	storeRouters.POST("/carts/products/:id", handlers.StoreHandler.AddToCart)
	storeRouters.PUT("/carts/products/:id", handlers.StoreHandler.UpdateCartItem)
	storeRouters.DELETE("/carts/products/:id", handlers.StoreHandler.RemoveFromCart)
	storeRouters.GET("/carts", handlers.StoreHandler.GetCart)
	//orders routes
	storeRouters.POST("/orders", handlers.StoreHandler.CreateOrder)
	storeRouters.PUT("/orders/:id", handlers.StoreHandler.CancelOrder)
	storeRouters.GET("/orders/:id", handlers.StoreHandler.GetOrder)
	storeRouters.GET("/orders", handlers.StoreHandler.ListOrders)
	storeRouters.GET("/orders/purchased", handlers.StoreHandler.GetStatistics)
	//review
	storeRouters.POST("/reviews", handlers.StoreHandler.CreateReview)
	storeRouters.GET("/reviews", handlers.StoreHandler.GetReview)
	storeRouters.GET("/reviews/product/:id", handlers.StoreHandler.GetReviewByProductId)
	storeRouters.DELETE("/reviews/:id", handlers.StoreHandler.DeleteReview)
}
