package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	healthRepo "prodcrud/internal/repository/health"
	"prodcrud/internal/repository/product"
	"prodcrud/internal/rest"
	healthHandler "prodcrud/internal/rest/handlers/health"
	productHandler "prodcrud/internal/rest/handlers/product"
	healthService "prodcrud/internal/usecase/health"
	productService "prodcrud/internal/usecase/product"
	"prodcrud/pkg/migration"
	"time"

	"prodcrud/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go.uber.org/dig"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	var (
		port = os.Getenv("PORT")
		host = os.Getenv("HOST")
		dsn  = os.Getenv("DSN")
		file = os.Getenv("FILE")
	)

	if err := migration.Migrate(file, dsn); err != nil {
		log.Fatalf("Error running migration: %s", err.Error())
	}

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
		productHandler.NewHandler,
		func(server *rest.Server) *http.Server {
			return &http.Server{
				Addr:              net.JoinHostPort(host, port),
				Handler:           server,
				ReadHeaderTimeout: 5 * time.Second,
			}
		},
	}

	container := dig.New()
	for _, dep := range deps {
		if err := container.Provide(dep); err != nil {
			return fmt.Errorf("failed to provide dependencies: %w", err)
		}
	}

	if err := container.Provide(product.NewRepo, dig.As(new(product.Repository))); err != nil {
		return fmt.Errorf("failed to provide dependency: %w", err)
	}

	if err := container.Provide(productService.NewService, dig.As(new(productService.ServiceInterface))); err != nil {
		return fmt.Errorf("failed to provide dependency: %w", err)
	}

	err := container.Invoke(func(server *rest.Server) {
		server.Init()
	})
	if err != nil {
		return fmt.Errorf("failed to init server: %w", err)
	}
	//nolint:wrapcheck //dig.Invoke returns error
	return container.Invoke(func(server *http.Server) error {
		return server.ListenAndServe()
	})
}
