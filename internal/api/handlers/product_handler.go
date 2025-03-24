package handlers

import (
	"net/http"
	"strconv"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	CategoryID  string  `json:"category_id" binding:"required"`
	InStock     bool    `json:"in_stock"`
	ImageURL    string  `json:"image_url"`
	Type        string  `json:"type"` // Mantido para compatibilidade
}

type ProductHandler struct {
	productService         *services.ProductService
	productCategoryService *services.ProductCategoryService
}

// Estrutura para a resposta paginada
type PaginatedResponse struct {
	Items       interface{} `json:"items"`
	TotalItems  int64       `json:"total_items"`
	TotalPages  int         `json:"total_pages"`
	CurrentPage int         `json:"current_page"`
	PageSize    int         `json:"page_size"`
	HasNext     bool        `json:"has_next"`
	HasPrev     bool        `json:"has_prev"`
}

func NewProductHandler(productService *services.ProductService, productCategoryService *services.ProductCategoryService) *ProductHandler {
	return &ProductHandler{
		productService:         productService,
		productCategoryService: productCategoryService,
	}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Obtém o restaurant_id do contexto
	restaurantID, exists := c.Get("requested_restaurant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant ID not found"})
		return
	}

	// Converte para UUID
	restID, ok := restaurantID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID format"})
		return
	}

	restaurantUUID, err := uuid.Parse(restID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	// Converte category_id para UUID
	categoryUUID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	// Verifica se a categoria existe e pertence ao restaurante
	category, err := h.productCategoryService.FindByID(restaurantUUID, categoryUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category not found"})
		return
	}

	// Verifica se a categoria pertence ao restaurante
	if category.RestaurantID != restaurantUUID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category does not belong to this restaurant"})
		return
	}

	// Configurar tipo de produto (mantido para compatibilidade)
	var productType models.ProductType
	if req.Type != "" {
		switch req.Type {
		case string(models.ProductTypeFood):
			productType = models.ProductTypeFood
		case string(models.ProductTypeDrink):
			productType = models.ProductTypeDrink
		case string(models.ProductTypeDessert):
			productType = models.ProductTypeDessert
		default:
			productType = models.ProductTypeFood // Valor padrão
		}
	} else {
		productType = models.ProductTypeFood // Valor padrão
	}

	product := &models.Product{
		RestaurantID: restaurantUUID,
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		CategoryID:   categoryUUID,
		Category:     category,
		Type:         productType,
		InStock:      req.InStock,
		ImageURL:     req.ImageURL,
	}

	if err := h.productService.Create(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, product)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	// Obtém o ID do produto da URL
	productID := c.Param("product_id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})
		return
	}

	// Converte para UUID
	prodUUID, err := uuid.Parse(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	restaurant_id := c.Param("restaurant_id")
	if restaurant_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	restaurant_uuid, err := uuid.Parse(restaurant_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	product, err := h.productService.GetByID(restaurant_uuid, prodUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) List(c *gin.Context) {
	// Obtém o restaurant_id do contexto
	restaurantID, exists := c.Get("requested_restaurant_id")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant ID not found"})
		return
	}

	// Converte para UUID
	restID, ok := restaurantID.(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID format"})
		return
	}

	restaurantUUID, err := uuid.Parse(restID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	// Parâmetros de paginação
	pageStr := c.DefaultQuery("page", "1")
	pageSizeStr := c.DefaultQuery("page_size", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// Filtragem por categoria
	categoryIDParam := c.Query("category_id")
	var category *models.ProductCategory

	if categoryIDParam != "" {
		catUUID, err := uuid.Parse(categoryIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID format"})
			return
		}

		// Retrieve the ProductCategory object
		category, err = h.productCategoryService.FindByID(restaurantUUID, catUUID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "category not found"})
			return
		}
	}

	// Filtragem por tipo (mantido para compatibilidade)
	// typeParam := c.Query("type")
	// var productType *models.ProductType

	// if typeParam != "" {
	// 	var pType models.ProductType
	// 	switch typeParam {
	// 	case string(models.ProductTypeFood):
	// 		pType = models.ProductTypeFood
	// 	case string(models.ProductTypeDrink):
	// 		pType = models.ProductTypeDrink
	// 	case string(models.ProductTypeDessert):
	// 		pType = models.ProductTypeDessert
	// 	default:
	// 		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product type"})
	// 		return
	// 	}
	// 	productType = &pType
	// }

	// Filtragem por disponibilidade em estoque
	inStockParam := c.Query("in_stock")
	var inStock *bool
	if inStockParam != "" {
		inStockValue := inStockParam == "true"
		inStock = &inStockValue
	}

	// Filtragem por nome (pesquisa parcial)
	nameSearch := c.Query("name")

	// Opções de ordenação
	sortBy := c.DefaultQuery("sort_by", "name")
	sortOrder := c.DefaultQuery("sort_order", "asc")
	// Buscar produtos com paginação
	products, totalItems, err := h.productService.ListWithPagination(
		restaurantUUID,
		page,
		pageSize,
		category,
		inStock,
		nameSearch,
		sortBy,
		sortOrder,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch products: " + err.Error()})
		return
	}

	// Calcular total de páginas
	totalPages := calculateTotalPages(totalItems, pageSize)

	// Construir resposta paginada
	response := PaginatedResponse{
		Items:       products,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
		HasNext:     page < totalPages,
		HasPrev:     page > 1,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandler) Update(c *gin.Context) {
	// Obtém o ID do produto da URL
	productID := c.Param("product_id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})
		return
	}

	// Converte para UUID
	prodUUID, err := uuid.Parse(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	restaurant_id := c.Param("restaurant_id")
	if restaurant_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	restaurant_uuid, err := uuid.Parse(restaurant_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	var req ProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.productService.GetByID(restaurant_uuid, prodUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	// Converte category_id para UUID
	categoryUUID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	// Verifica se a categoria existe
	_, err = h.productCategoryService.FindByID(restaurant_uuid, categoryUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category not found"})
		return
	}

	// Verificar se a categoria pertence ao mesmo restaurante que o produto
	if product.RestaurantID != product.Category.RestaurantID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category does not belong to this restaurant"})
		return
	}

	// Atualizar tipo se fornecido (mantido para compatibilidade)
	if req.Type != "" {
		switch req.Type {
		case string(models.ProductTypeFood):
			product.Type = models.ProductTypeFood
		case string(models.ProductTypeDrink):
			product.Type = models.ProductTypeDrink
		case string(models.ProductTypeDessert):
			product.Type = models.ProductTypeDessert
		}
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.CategoryID = categoryUUID
	product.InStock = req.InStock
	if req.ImageURL != "" {
		product.ImageURL = req.ImageURL
	}

	if err := h.productService.Update(product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *ProductHandler) Delete(c *gin.Context) {
	// Obtém o ID do produto da URL
	productID := c.Param("product_id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})
		return
	}

	// Converte para UUID
	prodUUID, err := uuid.Parse(productID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	restaurant_id := c.Param("restaurant_id")
	if restaurant_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	restaurant_uuid, err := uuid.Parse(restaurant_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	if err := h.productService.Delete(restaurant_uuid, prodUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}

func (h *ProductHandler) UpdateStock(c *gin.Context) {
	// Obtém o ID do produto da URL
	productID := c.Param("product_id")
	if productID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product ID is required"})
		return
	}

	// Converte para UUID
	prodUUID, err := uuid.Parse(productID)
	if err != nil {
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

	restaurant_id := c.Param("restaurant_id")
	if restaurant_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	restaurant_uuid, err := uuid.Parse(restaurant_id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	if err := h.productService.UpdateStock(restaurant_uuid, prodUUID, req.InStock); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product stock updated successfully"})
}

// Função auxiliar para calcular o número total de páginas
func calculateTotalPages(totalItems int64, pageSize int) int {
	if totalItems == 0 {
		return 1
	}

	totalPages := int(totalItems) / pageSize
	if int(totalItems)%pageSize > 0 {
		totalPages++
	}
	return totalPages
}
