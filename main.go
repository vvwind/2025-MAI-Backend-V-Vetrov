package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"rest-api/internal/config"
	"rest-api/internal/controller"
	"rest-api/internal/middleware"
	"rest-api/internal/repository"
	"rest-api/internal/service"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load configuration

	time.Sleep(5 * time.Second)
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	poolConfig, err := pgxpool.ParseConfig(cfg.Database.URL)
	if err != nil {
		log.Fatalf("Unable to parse database config: %v", err)
	}

	poolConfig.MaxConns = cfg.Database.MaxConnections
	poolConfig.MinConns = cfg.Database.MinConnections
	poolConfig.MaxConnLifetime = cfg.Database.MaxConnLifetime

	dbPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}
	defer dbPool.Close()

	//Verify database connection
	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	//Initialize repository with caching
	productPGRepo := repository.NewPostgresProductRepository(dbPool, rdb)
	userPGRepo := repository.NewPostgresUserRepository(dbPool)

	// Initialize services
	productService := service.NewProductService(productPGRepo)
	userService := service.NewUserService(userPGRepo)

	// Initialize controllers
	marketplaceController := controller.NewMarketplaceController(productService, userService)

	// Create router
	router := mux.NewRouter()

	// Register middleware
	router.Use(middleware.RecoveryMiddleware)
	router.Use(middleware.LoggingMiddleware)

	// Register routes
	marketplaceController.RegisterRoutes(router)

	// Start server
	log.Printf("Server starting on port %s...", cfg.Server.Port)
	if err := http.ListenAndServe(cfg.Server.Port, router); err != nil {
		log.Fatalf("Could not start server: %v", err)
	}
}
