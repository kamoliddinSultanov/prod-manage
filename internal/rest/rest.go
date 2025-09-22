package rest

import (
	"net/http"
	"prodcrud/internal/rest/handlers/health"
	"prodcrud/internal/rest/handlers/product"

	"github.com/gin-gonic/gin"
)

type Server struct {
	mux     *gin.Engine
	health  *health.Handler
	product *product.Handler
}

func NewServer(mux *gin.Engine, healthHandler *health.Handler, productHandler *product.Handler) *Server {
	return &Server{
		mux:     mux,
		health:  healthHandler,
		product: productHandler,
	}
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}
func (s *Server) Init() {
	s.mux.Use(gin.Logger())
	s.mux.Use(gin.Recovery())

	gr := s.mux.Group("/products")
	{
		gr.GET("/health", s.health.HealthCheck)

		gr.GET("/", s.product.GetAllProducts)
		gr.GET("/:id", s.product.GetProduct)
		gr.POST("/", s.product.CreateProduct)
		gr.PUT("/:id", s.product.UpdateProduct)
		gr.DELETE("/:id", s.product.DeleteProduct)
		gr.PUT("/:id/restore", s.product.RestoreProduct)
	}
}
