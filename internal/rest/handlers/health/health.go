package health

import (
	"context"
	"net/http"
	"prodcrud/internal/usecase/health"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *health.Service
}

func NewHandler(service *health.Service) *Handler {
	return &Handler{service: service}
}

// Check godoc
//
//	@Summary		Check service health
//	@Description	Ping the database and cache to verify service health
//	@Tags			health
//	@Accept			json
//	@Produce		json
//	@Failure		500	{object} map[string]string
//	@Success		200	{object} map[string]string
//	@Router			/products/health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	if err := h.service.Check(context.Background()); err != nil {
		http.Error(c.Writer, "service not working", http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "service is working"})
}
