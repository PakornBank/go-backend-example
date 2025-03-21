package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Handler defines the interface for health check HTTP requests.
type Handler interface {
	Check(c *gin.Context)
}

// handler handles health check HTTP requests.
type handler struct {
	db *gorm.DB
}

// NewHandler creates a new instance of handler with the provided DB connection.
func NewHandler(db *gorm.DB) Handler {
	return &handler{db: db}
}

// Check performs a health check on the application and its dependencies.
func (h *handler) Check(c *gin.Context) {
	// Check database connection
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Database connection error",
		})
		return
	}

	// Test database connectivity
	if err := sqlDB.Ping(); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Database ping failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Service is healthy",
	})
}
