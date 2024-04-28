package main

import (
	"bookstore/health"
	"bookstore/internal/api"
	"bookstore/internal/application"
	"bookstore/internal/application/config"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.Load()

	//app, err := application.Load()
	if err != nil {
		panic(err)
	}
	app := application.NewAppMock()
	r := setupRouter(app, config)
	if err := r.Run(fmt.Sprintf(":%d", config.AppPort)); err != nil {
		panic(err)
	}

}

func setupRouter(app *application.Application, config *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	// db, err := api.NewPostgresDB(&config.Config{
	// 	DBHost:     "localhost",
	// 	DBPort:     5432,
	// 	DBUser:     "postgres",
	// 	DBPassword: "admin",
	// 	DBName:     "bookstore",
	// })
	db, err := api.NewPostgresDB(config)
	if err != nil {
		panic(err)
	}
	bookstoreRepo := api.NewRepository(app, *db)
	bookStoreService := api.NewService(app, bookstoreRepo)
	bookStoreHandler := api.NewHandler(app, bookStoreService)
	r.GET("/health", health.Check)
	r.GET("/books", bookStoreHandler.GetAllBooks)
	r.POST("/accounts", bookStoreHandler.CreateAccount)
	r.POST("/orders", bookStoreHandler.PlaceOrder)
	r.GET("/order/history", bookStoreHandler.GetOrderHistory)
	r.GET("/users/:email", bookStoreHandler.GetUserIDByEmail)
	return r

}
