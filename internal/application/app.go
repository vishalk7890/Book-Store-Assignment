// package application

// import (
// 	"bookstore/db"
// 	"bookstore/internal/api"

// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"
// 	"time"

// 	"bookstore/internal/application/config"
// )

// // App represents the application
// type App struct {
// 	Config *config.Config
// 	// DB     *db.PostgresDB
// 	// Router http.Handler
// 	// Server *http.Server
// 	// API    api.Service
// }

// func Load() (*App, error) {
// 	cfg, err := config.Load()
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Construct DB connection string
// 	dbConnectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
// 		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

// 	pg, err := db.NewPostgresDB(dbConnectionString)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Initialize repository
// 	repo := api.NewRepository()

// 	apiService := api.NewService(repo)

// 	router := api.NewHandler(apiService)
// 	server := &http.Server{
// 		Addr:    fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
// 		Handler: router,
// 	}

// 	return &App{
// 		Config: cfg,
// 		DB:     pg,
// 		Router: router,
// 		Server: server,
// 		API:    apiService,
// 	}, nil
// }

// // Start starts the application
// func (app *App) Start() {
// 	go func() {
// 		fmt.Printf("Server listening on %s\n", app.Server.Addr)
// 		if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
// 			log.Fatalf("Failed to start server: %v", err)
// 		}
// 	}()

// 	sigChan := make(chan os.Signal, 1)
// 	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
// 	<-sigChan

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	if err := app.Server.Shutdown(ctx); err != nil {
// 		log.Fatalf("Failed to gracefully shutdown server: %v", err)
// 	}
// }

package application

import "bookstore/internal/application/config"

type Application struct {
	config *config.Config
}

func Load() (*Application, error) {
	app := Application{}
	config, err := config.Load()
	if err != nil {
		return &app, err

	}
	app.config = config
	return &app, nil
}
