package user

import (
	"github.com/PakornBank/go-backend-example/internal/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:generate mockgen -destination=./handler_mock.go -package=user github.com/PakornBank/go-backend-example/cmd/api/handler/user Handler

// Handler defines the interface for user-related HTTP requests.
type Handler interface {
	GetProfile(c *gin.Context)
}

// handler handles user-related HTTP requests.
type handler struct {
	service user.Service
}

// NewHandler creates a new instance of handler with the provided service.
func NewHandler(s user.Service) Handler {

	return &handler{service: s}
}

// GetProfile handles the request to retrieve the profile of the authenticated user.
func (h *handler) GetProfile(c *gin.Context) {
	id, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	res, err := h.service.GetUserByID(c.Request.Context(), id.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
