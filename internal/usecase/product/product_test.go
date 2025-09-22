package product

import (
	"context"
	"errors"
	"prodcrud/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_CreateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(Mock)
		service := NewService(mockRepo)

		p := &models.Product{
			Name:        "Test Product",
			Price:       1000,
			Quantity:    10,
			Description: "Test Product Description",
		}

		mockRepo.On("CreateProduct", mock.Anything, p).Return(nil).Once()

		err := service.CreateProduct(context.Background(), p)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
	t.Run("failed", func(t *testing.T) {
		mockRepo := new(Mock)
		service := NewService(mockRepo)
		p := &models.Product{
			Name:        "Test Product",
			Price:       1000,
			Quantity:    10,
			Description: "Test Product Description",
		}
		mockRepo.On("CreateProduct", mock.Anything, p).Return(
			errors.New("failed to create product usc")).Once()

		err := service.CreateProduct(context.Background(), p)
		assert.Error(t, err)
		assert.EqualError(t, err, "failed to create product usc")
		mockRepo.AssertExpectations(t)
	})
	t.Run("error with negative price", func(t *testing.T) {
		mockRepo := new(Mock)
		service := NewService(mockRepo)
		p := &models.Product{
			Name:        "Test Product",
			Price:       -1000,
			Quantity:    10,
			Description: "Test Product Description",
		}
		err := service.CreateProduct(context.Background(), p)
		assert.Error(t, err)
		assert.EqualError(t, err, "price cannot be negative or zero")
	})
}

func TestService_GetAllProducts(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(Mock)
		service := NewService(mockRepo)
		mockRepo.On("GetAllProducts", mock.Anything).Return([]*models.Product{
			{
				ID:          1,
				Name:        "Test Product",
				Price:       1000,
				Quantity:    10,
				Description: "Test Product Description",
			},
			{
				ID:          2,
				Name:        "Test Product 2",
				Price:       1000,
				Quantity:    10,
				Description: "Test Product Description 2",
			},
		}, nil).Once()
		products, err := service.GetAllProducts(context.Background())
		assert.NoError(t, err)
		assert.Len(t, products, 2)
		mockRepo.AssertExpectations(t)
	})
	t.Run("failed", func(t *testing.T) {
		mockRepo := new(Mock)
		service := NewService(mockRepo)
		mockRepo.On("GetAllProducts", mock.Anything).Return([]*models.Product{}, errors.New("failed to get products usc")).Once()
		products, err := service.GetAllProducts(context.Background())
		assert.Error(t, err)
		assert.Nil(t, products)
	})
}

func TestService_GetProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(Mock)
		service := NewService(mockRepo)
		mockRepo.On("GetProduct", mock.Anything, int64(1)).Return(&models.Product{
			ID:          1,
			Name:        "Test Product",
			Price:       1000,
			Quantity:    10,
			Description: "Test Product Description",
		}, nil).Once()
		product, err := service.GetProduct(context.Background(), 1)
		assert.NoError(t, err)
		assert.NotNil(t, product)
		mockRepo.AssertExpectations(t)
	})
	t.Run("failed", func(t *testing.T) {
		mockRepo := new(Mock)
		service := NewService(mockRepo)
		mockRepo.On("GetProduct", mock.Anything, int64(1)).Return(
			(*models.Product)(nil), errors.New("failed to get product usc")).Once()
		product, err := service.GetProduct(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, product)
	})
	t.Run("not found", func(t *testing.T) {
		mockRepo := new(Mock)
		service := NewService(mockRepo)
		mockRepo.On("GetProduct", mock.Anything, int64(1)).Return(
			(*models.Product)(nil), ErrProductNotFound).Once()
		product, err := service.GetProduct(context.Background(), 1)
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.EqualError(t, err, "product not found")
		mockRepo.AssertExpectations(t)
	})
}

func TestService_UpdateProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(Mock)
		service := NewService(mockRepo)

		existProd := &models.Product{
			ID:          1,
			Name:        "Test Product",
			Price:       1000,
			Quantity:    10,
			Description: "Test Product Description",
		}
		updatedProd := &models.Product{
			ID:          1,
			Name:        "Test Product Updated",
			Price:       2000000,
			Quantity:    20,
			Description: "Test Product Description Updated",
		}
		mockRepo.On("GetProduct", mock.Anything, int64(1)).Return(existProd, nil).Once()
		mockRepo.On("UpdateProduct", mock.Anything, updatedProd).Return(nil).Once()
		err := service.UpdateProduct(context.Background(), updatedProd)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNumberOfCalls(t, "UpdateProduct", 1)
		mockRepo.AssertCalled(t, "UpdateProduct", mock.Anything, updatedProd)
	})
	t.Run("failed", func(t *testing.T) {
		mockRepo := new(Mock)
		service := NewService(mockRepo)
		updatedProd := &models.Product{
			ID:          1,
			Name:        "Test Product Updated",
			Price:       2000000,
			Quantity:    20,
			Description: "Test Product Description Updated",
		}
		mockRepo.On("GetProduct", mock.Anything, int64(1)).Return(
			(*models.Product)(nil), errors.New("failed to get product usc")).Once()
		err := service.UpdateProduct(context.Background(), updatedProd)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNumberOfCalls(t, "UpdateProduct", 0)
		mockRepo.AssertNotCalled(t, "UpdateProduct", mock.Anything, updatedProd)
	})
	t.Run("update failed", func(t *testing.T) {
		mockRepo := new(Mock)
		service := NewService(mockRepo)
		existProd := &models.Product{
			ID:          1,
			Name:        "Test Product",
			Price:       1000,
			Quantity:    10,
			Description: "Test Product Description",
		}
		updatedProd := &models.Product{
			ID:          1,
			Name:        "Test Product Updated",
			Price:       2000000,
			Quantity:    20,
			Description: "Test Product Description Updated",
		}
		mockRepo.On("GetProduct", mock.Anything, int64(1)).Return(existProd, nil).Once()
		mockRepo.On("UpdateProduct", mock.Anything, updatedProd).Return(
			errors.New("failed to update product usc")).Once()
		err := service.UpdateProduct(context.Background(), updatedProd)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNumberOfCalls(t, "UpdateProduct", 1)
		mockRepo.AssertCalled(t, "UpdateProduct", mock.Anything, updatedProd)
	})
}

func TestService_DeleteProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(Mock)
		service := NewService(mockRepo)
		mockRepo.On("DeleteProduct", mock.Anything, int64(1)).Return(nil).Once()
		err := service.DeleteProduct(context.Background(), 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNumberOfCalls(t, "DeleteProduct", 1)
		mockRepo.AssertCalled(t, "DeleteProduct", mock.Anything, int64(1))
	})
	t.Run("failed", func(t *testing.T) {
		mockRepo := new(Mock)
		service := NewService(mockRepo)
		mockRepo.On("DeleteProduct", mock.Anything, int64(1)).Return(
			errors.New("failed to delete product usc")).Once()
		err := service.DeleteProduct(context.Background(), 1)
		assert.Error(t, err)
	})
}

func TestService_RestoreProduct(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(Mock)
		service := NewService(mockRepo)
		mockRepo.On("RestoreProduct", mock.Anything, int64(1)).Return(nil).Once()
		err := service.RestoreProduct(context.Background(), 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNumberOfCalls(t, "RestoreProduct", 1)
		mockRepo.AssertCalled(t, "RestoreProduct", mock.Anything, int64(1))
	})
	t.Run("failed", func(t *testing.T) {
		mockRepo := new(Mock)
		service := NewService(mockRepo)
		mockRepo.On("RestoreProduct", mock.Anything, int64(1)).Return(
			errors.New("failed to restore product usc")).Once()
		err := service.RestoreProduct(context.Background(), 1)
		assert.Error(t, err)
	})
}
