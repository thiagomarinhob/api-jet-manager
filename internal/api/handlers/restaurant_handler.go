package handlers

import (
	"net/http"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RestaurantRequest struct {
	Name             string `json:"name" binding:"required"`
	Description      string `json:"description"`
	Address          string `json:"address"`
	Phone            string `json:"phone"`
	Email            string `json:"email" binding:"required,email"`
	Logo             string `json:"logo"`
	SubscriptionPlan string `json:"subscription_plan"`
	Status           string `json:"status"`
}

type RestaurantHandler struct {
	restaurantService *services.RestaurantService
	userService       *services.AuthService
}

func NewRestaurantHandler(restaurantService *services.RestaurantService, userService *services.AuthService) *RestaurantHandler {
	return &RestaurantHandler{
		restaurantService: restaurantService,
		userService:       userService,
	}
}

// Create - apenas superadmins podem criar restaurantes
func (h *RestaurantHandler) Create(c *gin.Context) {
	var req RestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar o status
	var status models.SubscriptionStatus
	switch req.Status {
	case string(models.SubscriptionStatusActive):
		status = models.SubscriptionStatusActive
	case string(models.SubscriptionStatusInactive):
		status = models.SubscriptionStatusInactive
	case string(models.SubscriptionStatusTrial):
		status = models.SubscriptionStatusTrial
	case "":
		status = models.SubscriptionStatusTrial
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription status"})
		return
	}

	restaurant := &models.Restaurant{
		Name:             req.Name,
		Description:      req.Description,
		Address:          req.Address,
		Phone:            req.Phone,
		Email:            req.Email,
		Logo:             req.Logo,
		SubscriptionPlan: req.SubscriptionPlan,
		Status:           status,
	}

	if err := h.restaurantService.Create(restaurant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, restaurant)
}

// GetByID - superadmins podem ver qualquer restaurante, outros apenas o seu
func (h *RestaurantHandler) GetByID(c *gin.Context) {
	id := c.Param("restaurant_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	restaurantID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	restaurant, err := h.restaurantService.GetByID(restaurantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "restaurant not found"})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

// List - superadmins veem todos os restaurantes, outros apenas o seu
func (h *RestaurantHandler) List(c *gin.Context) {
	// Verifique se é um superadmin
	userType, _ := c.Get("user_type")
	restaurantID, _ := c.Get("restaurant_id")

	// Se for superadmin, liste todos os restaurantes
	if userType.(models.UserType) == models.UserTypeSuperAdmin {
		// Filtro opcional por status
		statusFilter := c.Query("status")
		nameFilter := c.Query("name")

		var restaurants []models.Restaurant
		var err error

		if statusFilter != "" {
			var status models.SubscriptionStatus
			switch statusFilter {
			case string(models.SubscriptionStatusActive):
				status = models.SubscriptionStatusActive
			case string(models.SubscriptionStatusInactive):
				status = models.SubscriptionStatusInactive
			case string(models.SubscriptionStatusTrial):
				status = models.SubscriptionStatusTrial
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status filter"})
				return
			}
			restaurants, err = h.restaurantService.GetByStatus(status)
		} else if nameFilter != "" {
			restaurants, err = h.restaurantService.SearchByName(nameFilter)
		} else {
			restaurants, err = h.restaurantService.List()
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch restaurants"})
			return
		}

		c.JSON(http.StatusOK, restaurants)
		return
	}

	// Se não for superadmin, retorna apenas o restaurante do usuário
	if restaurantID == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "you don't have a restaurant assigned"})
		return
	}

	restaurant, err := h.restaurantService.GetByID(restaurantID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "restaurant not found"})
		return
	}

	c.JSON(http.StatusOK, []models.Restaurant{*restaurant})
}

// Update - superadmins podem atualizar qualquer restaurante, admins apenas o seu
func (h *RestaurantHandler) Update(c *gin.Context) {
	id := c.Param("restaurant_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	var req RestaurantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	restaurantID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	restaurant, err := h.restaurantService.GetByID(restaurantID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "restaurant not found"})
		return
	}

	// Apenas superadmins podem alterar o status e o plano
	userType, _ := c.Get("user_type")
	if userType.(models.UserType) == models.UserTypeSuperAdmin {
		if req.Status != "" {
			switch req.Status {
			case string(models.SubscriptionStatusActive):
				restaurant.Status = models.SubscriptionStatusActive
			case string(models.SubscriptionStatusInactive):
				restaurant.Status = models.SubscriptionStatusInactive
			case string(models.SubscriptionStatusTrial):
				restaurant.Status = models.SubscriptionStatusTrial
			default:
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription status"})
				return
			}
		}

		if req.SubscriptionPlan != "" {
			restaurant.SubscriptionPlan = req.SubscriptionPlan
		}
	}

	// Campos que qualquer admin pode alterar
	if req.Name != "" {
		restaurant.Name = req.Name
	}
	if req.Description != "" {
		restaurant.Description = req.Description
	}
	if req.Address != "" {
		restaurant.Address = req.Address
	}
	if req.Phone != "" {
		restaurant.Phone = req.Phone
	}
	if req.Email != "" {
		restaurant.Email = req.Email
	}
	if req.Logo != "" {
		restaurant.Logo = req.Logo
	}

	if err := h.restaurantService.Update(restaurant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, restaurant)
}

// Delete - apenas superadmins podem excluir restaurantes
func (h *RestaurantHandler) Delete(c *gin.Context) {
	id := c.Param("restaurant_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	restaurantID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	if err := h.restaurantService.Delete(restaurantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "restaurant deleted successfully"})
}

// UpdateStatus - apenas superadmins podem atualizar o status
func (h *RestaurantHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("restaurant_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var status models.SubscriptionStatus
	switch req.Status {
	case string(models.SubscriptionStatusActive):
		status = models.SubscriptionStatusActive
	case string(models.SubscriptionStatusInactive):
		status = models.SubscriptionStatusInactive
	case string(models.SubscriptionStatusTrial):
		status = models.SubscriptionStatusTrial
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription status"})
		return
	}

	restaurantID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant ID"})
		return
	}

	if err := h.restaurantService.UpdateStatus(restaurantID, status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "restaurant status updated successfully"})
}
