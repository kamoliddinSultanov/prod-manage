package main

import (
	"log"
	"net"
	"net/http"
	"os"
	healthRepo "prodcrud/internal/repository/health"
	"prodcrud/internal/rest"
	healthHandler "prodcrud/internal/rest/handlers/health"
	healthService "prodcrud/internal/usecase/health"

	"prodcrud/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/dig"
)

func main() {
	var (
		port = "7777"
		host = "0.0.0.0"
		dsn  = "postgres://user:pass@localhost:5434/prod_db?sslmode=disable"
	)

	if err := execute(host, port, dsn); err != nil {
		log.Println(err)
		os.Exit(1)
	}

}

func execute(host, port, dsn string) error {
	deps := []interface{}{
		func() (*pgxpool.Pool, error) {
			return db.NewDB(dsn)
		},
		gin.New,
		healthRepo.NewRepo,
		healthService.NewService,
		healthHandler.NewHandler,
		rest.NewServer,
		func(server *rest.Server) *http.Server {
			return &http.Server{
				Addr:    net.JoinHostPort(host, port),
				Handler: server,
			}
		},
	}

	container := dig.New()
	for _, dep := range deps {
		if err := container.Provide(dep); err != nil {
			return err
		}
	}
	err := container.Invoke(func(server *rest.Server) {
		server.Init()
	})
	if err != nil {
		return err
	}
	return container.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})
}
