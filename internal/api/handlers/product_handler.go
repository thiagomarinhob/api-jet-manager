// internal/api/handlers/product_handler.go
package handlers

import (
	"net/http"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Category    string  `json:"category" binding:"required"`
	InStock     bool    `json:"in_stock"`
}

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(productService *services.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validação da categoria
	var category models.ProductCategory
	switch req.Category {
	case string(models.ProductCategoryFood):
		category = models.ProductCategoryFood
	case string(models.ProductCategoryDrink):
		category = models.ProductCategoryDrink
	case string(models.ProductCategoryDessert):
		category = models.ProductCategoryDessert
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product category"})
		return
	}

	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    category,
		InStock:     req.InStock,
	}

	if err := h.productService.Create(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	id := c.Param("product_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID format"})
		return
	}
	product, err := h.productService.GetByID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) List(c *gin.Context) {
	categoryParam := c.Query("category")

	if categoryParam != "" {
		var category models.ProductCategory
		switch categoryParam {
		case string(models.ProductCategoryFood):
			category = models.ProductCategoryFood
		case string(models.ProductCategoryDrink):
			category = models.ProductCategoryDrink
		case string(models.ProductCategoryDessert):
			category = models.ProductCategoryDessert
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product category"})
			return
		}

		products, err := h.productService.GetByCategory(category)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
			return
		}

		c.JSON(http.StatusOK, products)
		return
	}

	products, err := h.productService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) Update(c *gin.Context) {
	id := c.Param("product_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID format"})
		return
	}
	product, err := h.productService.GetByID(parsedID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	// Validação da categoria
	var category models.ProductCategory
	switch req.Category {
	case string(models.ProductCategoryFood):
		category = models.ProductCategoryFood
	case string(models.ProductCategoryDrink):
		category = models.ProductCategoryDrink
	case string(models.ProductCategoryDessert):
		category = models.ProductCategoryDessert
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product category"})
		return
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Category = category
	product.InStock = req.InStock

	if err := h.productService.Update(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	id := c.Param("product_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID format"})
		return
	}
	if err := h.productService.Delete(parsedID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

func (h *ProductHandler) UpdateStock(c *gin.Context) {
	id := c.Param("product_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var req struct {
		InStock bool `json:"in_stock" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID format"})
		return
	}
	if err := h.productService.UpdateStock(parsedID, req.InStock); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product stock updated successfully"})
}
