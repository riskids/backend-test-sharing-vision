package main

import (
	"fmt"
	"log"

	"article-microservice/internal/config"
	httpDelivery "article-microservice/internal/delivery/http"
	"article-microservice/internal/repository"
	"article-microservice/internal/service"
	"article-microservice/pkg/mysql"

	"github.com/go-playground/validator/v10"
)

func main() {
	// Load configuration from .env
	cfg := config.LoadConfig()

	// Connect to MySQL database
	db, err := mysql.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Successfully connected to MySQL database")

	// Initialize validator
	validate := validator.New()

	// Initialize layers with Dependency Injection
	// Repository Layer (Data Access)
	articleRepo := repository.NewArticleRepository(db)

	// Service Layer (Business Logic)
	articleService := service.NewArticleService(articleRepo)

	// Delivery Layer (HTTP Handlers & Routing)
	router := httpDelivery.SetupRouter(articleService, validate)

	// Start the server
	port := ":" + cfg.AppPort
	fmt.Printf("Server is running on port %s\n", port)

	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}