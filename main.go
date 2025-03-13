package main

import (
	"context"
	"fmt"
	"go-gin-e-comm/admin"
	"go-gin-e-comm/auth"
	"go-gin-e-comm/cart"
	"go-gin-e-comm/configs"
	"go-gin-e-comm/db"
	"go-gin-e-comm/logger"
	"go-gin-e-comm/products"
	"go-gin-e-comm/user"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := configs.LoadConfig("./configs")
	if err != nil {
		log.Fatal("Failed to load config.yaml", err)
	}
	logger.InitSlogLogger(cfg)
	mongoDB, err := db.InitMongoDB(cfg.Database.Uri)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB", err)
	}
	userRepo := user.NewUserRepository(mongoDB)
	productsRepo := products.NewProductRepository(mongoDB)
	cartRepo := cart.NewCartRepository(mongoDB)

	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	authGroup := r.Group("/auth")
	auth.SetupRoutes(authGroup, userRepo, cfg.JWT.Secret)

	protected := r.Group("/", auth.AuthenticationMiddleware(userRepo, cfg.JWT.Secret))

	adminGroup := protected.Group("/admin", auth.AdminMiddleware(userRepo))
	admin.SetupRoutes(adminGroup, productsRepo)

	userGroup := protected.Group("/user")
	user.SetupRoutes(userGroup, userRepo)

	productsGroup := protected.Group("/products")
	products.SetupRoutes(productsGroup, productsRepo)

	cartGroup := protected.Group("/cart")
	cart.SetupRoutes(cartGroup, cartRepo)

	port := fmt.Sprintf(":%d", cfg.Server.Port)
	srv := &http.Server{
		Addr:    port,
		Handler: r,
	}

	go func() {
		slog.Info(
			"Server is running",
			slog.String("port", port),
			slog.String("mode", cfg.Server.Mode),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}

	slog.Info("Server exited gracefully")
}
