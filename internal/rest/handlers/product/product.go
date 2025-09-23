package product

import (
	"errors"
	"net/http"
	"prodcrud/internal/models"
	"prodcrud/internal/usecase/product"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service product.ServiceInterface
}

func NewHandler(service product.ServiceInterface) *Handler {
	return &Handler{service: service}
}

// CreateProduct godoc
//
//	@Summary		Create a new product
//	@Description	Create a new product with the provided details
//	@Tags			products
//
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.ProductResponse	true	"Product details"
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Success		201		{object}	map[string]string
//	@Router			/products/ [post]
func (h *Handler) CreateProduct(c *gin.Context) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateProduct(c, &p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":     "a product has been created",
		"name":        p.Name,
		"price":       p.Price,
		"quantity":    p.Quantity,
		"description": p.Description,
	})
}

// GetProduct godoc
//
//	@Summary		Get a product by ID
//	@Description	Get the details of a product by its ID
//	@Tags			products
//
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int64	true	"Product ID"
//	@Failure		400	{object}	map[string]string
//	@Failure		404	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Success		200	{object}	models.Product
//	@Router			/products/{id} [get]
func (h *Handler) GetProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}
	prod, err := h.service.GetProduct(c, id)
	if err != nil {
		if errors.Is(err, product.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prod)
}

// GetAllProducts godoc
//
//	@Summary		Get all products
//	@Description	Get a list of all products with optional filters
//	@Tags			products
//
//	@Accept			json
//	@Produce		json
//	@Failure		500	{object}	map[string]string
//	@Success		200	{object}	[]models.Product
//	@Router			/products/ [get]
func (h *Handler) GetAllProducts(c *gin.Context) {
	products, err := h.service.GetAllProducts(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// UpdateProduct godoc
//
//	@Summary		Update an existing product
//	@Description	Update the details of an existing product by ID
//	@Tags			products
//
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int64			true	"Product ID"
//	@Param			request	body		models.ProductResponse	true	"Product details"
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Success		200		{object}	models.Product
//	@Router			/products/{id} [put]
func (h *Handler) UpdateProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	p.ID = id

	if err := h.service.UpdateProduct(c, &p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "a product has been updated"})
}

// DeleteProduct godoc
//
//	@Summary		Delete a product
//	@Description	Delete an existing product by ID
//	@Tags			products
//
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int64	true	"Product ID"
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Success		204	{object}	map[string]string
//	@Router			/products/{id} [delete]
func (h *Handler) DeleteProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}
	if err := h.service.DeleteProduct(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "a product has been deleted"})
}

// RestoreProduct godoc
//
//	@Summary		Restore a deleted product
//	@Description	Restore a soft-deleted product by ID
//	@Tags			products
//
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int64	true	"Product ID"
//	@Failure		400	{object}	map[string]string
//	@Failure		500	{object}	map[string]string
//	@Success		200	{object}	map[string]string
//	@Router			/products/{id}/restore [put]
func (h *Handler) RestoreProduct(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}
	if err := h.service.RestoreProduct(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "a product has been restored"})
}
