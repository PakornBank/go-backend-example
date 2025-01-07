package main

import (
	"log"

	"github.com/PakornBank/go-backend-example/internal/common/config"
	"github.com/PakornBank/go-backend-example/internal/common/database"
	"github.com/PakornBank/go-backend-example/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	db, err := database.NewDataBase(cfg)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	r := gin.Default()
	router.SetupRoutes(r, db, cfg)

	log.Printf("Server running on port %s\n", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
