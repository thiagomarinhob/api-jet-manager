package handlers

import (
	"net/http"
	"strconv"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProductCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type ProductCategoryHandler struct {
	categoryService *services.ProductCategoryService
}

func NewProductCategoryHandler(categoryService *services.ProductCategoryService) *ProductCategoryHandler {
	return &ProductCategoryHandler{
		categoryService: categoryService,
	}
}

// Create cria uma nova categoria de produto
func (h *ProductCategoryHandler) Create(c *gin.Context) {
	var req ProductCategoryRequest
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

	// Verifica se já existe uma categoria com o mesmo nome
	existingCategory, _ := h.categoryService.FindByName(restaurantUUID, req.Name)
	if existingCategory != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "a category with this name already exists"})
		return
	}

	category := &models.ProductCategory{
		RestaurantID: restaurantUUID,
		Name:         req.Name,
		Description:  req.Description,
		Active:       req.Active,
	}

	if !category.Active {
		category.Active = true // Define como true por padrão se não for especificado
	}

	if err := h.categoryService.Create(category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

// GetByID retorna uma categoria pelo ID
func (h *ProductCategoryHandler) GetByID(c *gin.Context) {
	// Obtém o ID da categoria da URL
	categoryID := c.Param("category_id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID is required"})
		return
	}

	// Converte para UUID
	catUUID, err := uuid.Parse(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
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

	category, err := h.categoryService.FindByID(restaurant_uuid, catUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	c.JSON(http.StatusOK, category)
}

// List lista todas as categorias de um restaurante com paginação e filtros
func (h *ProductCategoryHandler) List(c *gin.Context) {
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

	// Calcular o offset com base na página e tamanho da página
	offset := (page - 1) * pageSize

	// Filtragem por status (ativo/inativo)
	activeParam := c.Query("active")
	var active *bool
	if activeParam != "" {
		activeValue := activeParam == "true"
		active = &activeValue
	}

	// Filtragem por nome (pesquisa parcial)
	nameSearch := c.Query("name")

	// Opções de ordenação
	sortBy := c.DefaultQuery("sort_by", "name")
	sortOrder := c.DefaultQuery("sort_order", "asc")

	// Buscar categorias com paginação e filtros
	categories, totalItems, err := h.categoryService.FindWithFilters(
		restaurantUUID,
		offset,
		pageSize,
		active,
		nameSearch,
		sortBy,
		sortOrder,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch categories: " + err.Error()})
		return
	}

	// Calcular total de páginas
	totalPages := calculateTotalPages(totalItems, pageSize)

	// Construir resposta paginada
	response := PaginatedResponse{
		Items:       categories,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
		HasNext:     page < totalPages,
		HasPrev:     page > 1,
	}

	c.JSON(http.StatusOK, response)
}

// ListActive lista apenas as categorias ativas de um restaurante
func (h *ProductCategoryHandler) ListActive(c *gin.Context) {
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

	categories, err := h.categoryService.FindActive(restaurantUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch active categories"})
		return
	}

	c.JSON(http.StatusOK, categories)
}

// Update atualiza uma categoria existente
func (h *ProductCategoryHandler) Update(c *gin.Context) {
	// Obtém o ID da categoria da URL
	categoryID := c.Param("category_id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID is required"})
		return
	}

	// Converte para UUID
	catUUID, err := uuid.Parse(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	var req ProductCategoryRequest
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

	// Recupera a categoria existente
	category, err := h.categoryService.FindByID(restaurant_uuid, catUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	// Verifica se já existe outra categoria com o mesmo nome
	if req.Name != category.Name {
		existingCategory, _ := h.categoryService.FindByName(category.RestaurantID, req.Name)
		if existingCategory != nil && existingCategory.ID != category.ID {
			c.JSON(http.StatusBadRequest, gin.H{"error": "another category with this name already exists"})
			return
		}
	}

	// Atualiza os campos
	category.Name = req.Name
	category.Description = req.Description
	category.Active = req.Active

	if err := h.categoryService.Update(restaurant_uuid, category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, category)
}

// Delete remove uma categoria
func (h *ProductCategoryHandler) Delete(c *gin.Context) {
	// Obtém o ID da categoria da URL
	categoryID := c.Param("category_id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID is required"})
		return
	}

	// Converte para UUID
	catUUID, err := uuid.Parse(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
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

	// Verifica se existem produtos associados a esta categoria
	// Isso deveria ser verificado em um serviço ou repository dedicado
	// Por ora, vamos supor que o serviço já trata isso internamente

	if err := h.categoryService.Delete(restaurant_uuid, catUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}

// UpdateStatus atualiza apenas o status (ativo/inativo) de uma categoria
func (h *ProductCategoryHandler) UpdateStatus(c *gin.Context) {
	// Obtém o ID da categoria da URL
	categoryID := c.Param("category_id")
	if categoryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "category ID is required"})
		return
	}

	// Converte para UUID
	catUUID, err := uuid.Parse(categoryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category ID"})
		return
	}

	var req struct {
		Active bool `json:"active" binding:"required"`
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

	// Recupera a categoria existente
	category, err := h.categoryService.FindByID(restaurant_uuid, catUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

	// Atualiza apenas o campo active
	category.Active = req.Active

	if err := h.categoryService.Update(restaurant_uuid, category); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "category status updated successfully"})
}
