package rest

import (
	"net/http"
	"prodcrud/internal/rest/handlers/health"

	"github.com/gin-gonic/gin"
)

type Server struct {
	mux    *gin.Engine
	health *health.Handler
}

func NewServer(mux *gin.Engine, healthHandler *health.Handler) *Server {
	return &Server{
		mux:    mux,
		health: healthHandler,
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
	}
}
