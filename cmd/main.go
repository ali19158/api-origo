package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/online-shop/internal/config"
	"github.com/online-shop/internal/database"
	"github.com/online-shop/internal/handler"
	"github.com/online-shop/internal/repository"
	"github.com/online-shop/internal/router"
	"github.com/online-shop/internal/service"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Connect to PostgreSQL
	pool, err := database.NewPostgresPool(cfg.Database.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()
	log.Println("connected to PostgreSQL")

	// Repositories
	userRepo := repository.NewUserRepository(pool)
	productRepo := repository.NewProductRepository(pool)
	orderRepo := repository.NewOrderRepository(pool)
	categoryRepo := repository.NewCategoryRepository(pool)

	// Services
	userSvc := service.NewUserService(userRepo, cfg.JWT)
	productSvc := service.NewProductService(productRepo)
	orderSvc := service.NewOrderService(orderRepo, productRepo)
	categorySvc := service.NewCategoryService(categoryRepo)

	// Handlers
	userH := handler.NewUserHandler(userSvc)
	productH := handler.NewProductHandler(productSvc)
	orderH := handler.NewOrderHandler(orderSvc)
	categoryH := handler.NewCategoryHandler(categorySvc)

	// Router
	r := router.New(cfg.JWT.Secret, userH, productH, orderH, categoryH)

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("server starting on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
