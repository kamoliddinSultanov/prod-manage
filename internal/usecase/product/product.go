package product

import (
	"context"
	"errors"
	"fmt"
	"prodcrud/internal/models"
	"prodcrud/internal/repository/product"
)

type ServiceInterface interface {
	CreateProduct(ctx context.Context, p *models.Product) error
	GetAllProducts(ctx context.Context) ([]*models.Product, error)
	GetProduct(ctx context.Context, id int64) (*models.Product, error)
	UpdateProduct(ctx context.Context, p *models.Product) error
	DeleteProduct(ctx context.Context, id int64) error
	RestoreProduct(ctx context.Context, id int64) error
}

type Service struct {
	repo product.Repository
}

func NewService(repo product.Repository) ServiceInterface {
	return &Service{repo: repo}
}
func (s *Service) CreateProduct(ctx context.Context, p *models.Product) error {
	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Price <= 0 {
		return errors.New("price cannot be negative or zero")
	}
	if p.Quantity <= 0 {
		return errors.New("quantity cannot be negative or zero")
	}
	if p.Description == "" {
		return errors.New("description is required")
	}

	if err := s.repo.CreateProduct(ctx, p); err != nil {
		return errors.New("failed to create product usc")
	}
	return nil
}

func (s *Service) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	products, err := s.repo.GetAllProducts(ctx)
	if err != nil {
		return nil, errors.New("failed to get products usc")
	}
	return products, nil
}

func (s *Service) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	prod, err := s.repo.GetProduct(ctx, id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, fmt.Errorf("failed to get product usc: %w", err)
	}
	return prod, nil
}

func (s *Service) UpdateProduct(ctx context.Context, p *models.Product) error {
	upd, err := s.repo.GetProduct(ctx, p.ID)
	if err != nil {
		return fmt.Errorf("failed to get product usc: %w", err)
	}

	if p.Name != "" {
		upd.Name = p.Name
	}
	if p.Price != 0 {
		upd.Price = p.Price
	}
	if p.Quantity != 0 {
		upd.Quantity = p.Quantity
	}
	if p.Description != "" {
		upd.Description = p.Description
	}

	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Price <= 0 {
		return errors.New("price cannot be negative or zero")
	}
	if p.Quantity <= 0 {
		return errors.New("quantity cannot be negative or zero")
	}
	if p.Description == "" {
		return errors.New("description is required")
	}

	if err := s.repo.UpdateProduct(ctx, upd); err != nil {
		return errors.New("failed to update product usc")
	}
	return nil
}

func (s *Service) DeleteProduct(ctx context.Context, id int64) error {
	if err := s.repo.DeleteProduct(ctx, id); err != nil {
		return errors.New("failed to delete product usc")
	}
	return nil
}

func (s *Service) RestoreProduct(ctx context.Context, id int64) error {
	if err := s.repo.RestoreProduct(ctx, id); err != nil {
		return errors.New("failed to restore product usc")
	}
	return nil
}

var (
	ErrProductNotFound = errors.New("product not found")
)
