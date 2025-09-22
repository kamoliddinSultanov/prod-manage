package product

import (
	"context"
	"prodcrud/internal/models"

	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) CreateProduct(ctx context.Context, p *models.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *Mock) GetAllProducts(ctx context.Context) ([]*models.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *Mock) GetProduct(ctx context.Context, id int64) (*models.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *Mock) UpdateProduct(ctx context.Context, p *models.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *Mock) DeleteProduct(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *Mock) RestoreProduct(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
