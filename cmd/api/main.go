package main

import (
	"context"
	"github.com/PakornBank/go-backend-example/cmd/api/di"
	"github.com/PakornBank/go-backend-example/cmd/api/routes"
	"github.com/PakornBank/go-backend-example/internal/common/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	container := di.NewContainer(cfg)

	r := gin.Default()
	routes.SetupRoutes(r, container)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server running on port %s\n", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Shut down the server
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	// Close database connections
	if db, err := container.GetDB(); err == nil {
		sqlDB, err := db.DB()
		if err == nil {
			sqlDB.Close()
		}
	}

	log.Println("Server exiting")
}
