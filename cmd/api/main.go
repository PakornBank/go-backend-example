package main

import (
	"log"

	"github.com/PakornBank/go-backend-example/internal/common/config"
	"github.com/PakornBank/go-backend-example/internal/di"
	"github.com/PakornBank/go-backend-example/internal/router"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	container := di.NewContainer(cfg)

	r := gin.Default()
	router.SetupRoutes(r, container)

	log.Printf("Server running on port %s\n", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatal("failed to start server: ", err)
	}
}
