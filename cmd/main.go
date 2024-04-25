package main

import (
	"bookstore/db"
	"bookstore/internal/api"
	"bookstore/internal/application"
	"bookstore/internal/application/config"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Config{}

	pg, _ := db.NewPostgresDB(&cfg)

	// Initialize API service
	apiService := api.NewService(api.NewRepository())

	// Initialize the application

	app := application.App(cfg, pg, apiService)

	// Initialize Gin engine
	router := gin.Default()

	// Setup routes
	setupRouter(router, app)

	// Start the Gin server
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupRouter(app *application.Application) *gin.Engine {
	r := gin.Default()

	// Define your endpoints here

	bookStore := api.NewRepository(app, db)
	bookstoreService := api.NewService(app, bookStore)
	bookStoreHandler := api.NewHandler(app, bookstoreService)
	r.GET("/create", bookStoreHandler.CreateAccount)

	// Add other endpoints...

	return r
}
