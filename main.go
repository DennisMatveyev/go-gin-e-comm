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
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := configs.MustLoadConfig("./configs")
	log := logger.Setup(cfg)
	mongoDB := db.MustInitMongoDB(cfg.Database.Uri, log)
	defer mongoDB.Client().Disconnect(context.Background())

	userRepo := user.NewUserRepository(mongoDB, log)
	productsRepo := products.NewProductRepository(mongoDB, log)
	cartRepo := cart.NewCartRepository(mongoDB, log)

	gin.SetMode(cfg.Server.Mode)
	r := gin.Default()

	authGroup := r.Group("/auth")
	auth.SetupRoutes(authGroup, userRepo, log, cfg.JWT.Secret)

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
		log.Info("Server is running", "port", port, "mode", cfg.Server.Mode)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server failed to start", "error", err.Error())
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err.Error())
		os.Exit(1)
	}

	log.Info("Server exited gracefully")
}
