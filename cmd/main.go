package main

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "marketplace-api/docs"
	"marketplace-api/internal/api"
	"marketplace-api/internal/config"
	"marketplace-api/internal/logger"
	"marketplace-api/pkg/database"
)

// @title           Duken-API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:4000
// @BasePath  /api/

// @securityDefinitions.basic.apikey  BearerToken
// @in header
// @name Authorization

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	logger.InitLogger()
	log := logger.GetLogger()

	cfg := config.LoadConfig()

	db, err := database.InitDB(cfg) // Initialize the database
	if err != nil {
		log.Fatalf("Failed to connect to database: %s", err.Error())
	}
	d, _ := db.DB()
	defer func(d *sql.DB) {
		err = d.Close()
		if err != nil {
			log.Fatalf("Failed to close database: %s", err.Error())
		}
	}(d)
	log.Info("database connection pool established")

	router := gin.Default() // Create a Gin router

	api.NewServer(router, db, log, cfg) // Initialize API server

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	err = router.Run(cfg.ServerAddress)
	if err != nil {
		log.Fatalf("Failed to run server: %s", err.Error())
	}
}
