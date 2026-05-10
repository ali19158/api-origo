package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/online-shop/internal/config"
	"github.com/online-shop/internal/database"
	"github.com/online-shop/internal/handler"
	"github.com/online-shop/internal/logger"
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

	// Sentry
	sentryCfg := config.InitSentry()
	defer config.FlushSentry()

	// logger
	logger.Init(context.Background())

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
	//orderRepo := repository.NewOrderRepository(pool)
	categoryRepo := repository.NewCategoryRepository(pool)
	mediaRepo := repository.NewMediaRepository(pool)
	attrRepo := repository.NewAttributeRepository(pool)

	// Services
	mediaSvc := service.NewMediaService(mediaRepo, cfg.AdminURL)
	attrSvc := service.NewAttributeService(attrRepo)
	userSvc := service.NewUserService(userRepo, cfg.JWT)
	productSvc := service.NewProductService(productRepo, mediaSvc, attrSvc)
	//orderSvc := service.NewOrderService(orderRepo, productRepo)
	categorySvc := service.NewCategoryService(categoryRepo, mediaSvc)
	webhookSvc := service.NewWebhookService(cfg.TelegramToken, cfg.TelegramChatID)

	// Handlers
	userH := handler.NewUserHandler(userSvc)
	productH := handler.NewProductHandler(productSvc)
	//orderH := handler.NewOrderHandler(orderSvc)
	categoryH := handler.NewCategoryHandler(categorySvc)
	webhookH := handler.NewWebhookHandler(webhookSvc)

	webhookSecret := ""
	if sentryCfg != nil {
		webhookSecret = sentryCfg.WebhookSecret
	}

	// Router
	r := router.New(
		cfg.JWT.Secret, webhookSecret, userH, productH, categoryH, webhookH,
	)

	var sHandler http.Handler = r
	if sentryCfg != nil && sentryCfg.Handler != nil {
		sHandler = sentryCfg.Handler.Handle(r)
	}

	// Start server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Default.Info("server starting on %s", addr)
	if err := http.ListenAndServe(addr, sHandler); err != nil {
		logger.Default.Fatal("server error: %v", err)
	}
}
