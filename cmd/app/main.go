package main

import (
	"log"
	"net"
	"net/http"
	"os"
	healthRepo "prodcrud/internal/repository/health"
	productRepo "prodcrud/internal/repository/product"
	"prodcrud/internal/rest"
	healthHandler "prodcrud/internal/rest/handlers/health"
	productHandler "prodcrud/internal/rest/handlers/product"
	healthService "prodcrud/internal/usecase/health"
	productService "prodcrud/internal/usecase/product"
	"prodcrud/pkg/migration"

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
		productRepo.NewRepo,
		productService.NewService,
		productHandler.NewHandler,
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
