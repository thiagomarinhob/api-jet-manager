// internal/api/handlers/auth_handler.go
package handlers

import (
	"fmt"
	"net/http"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Name         string     `json:"name" binding:"required"`
	Email        string     `json:"email" binding:"required,email"`
	Password     string     `json:"password" binding:"required,min=6"`
	Type         string     `json:"type"`
	RestaurantID *uuid.UUID `json:"restaurant_id"`
}

type AuthHandler struct {
	authService       *services.AuthService
	restaurantService *services.RestaurantService
}

func NewAuthHandler(authService *services.AuthService, restaurantService *services.RestaurantService) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		restaurantService: restaurantService,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, user, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Verificar se o restaurante está ativo (exceto para superadmin)
	if user.Type != models.UserTypeSuperAdmin && user.RestaurantID != nil {
		isActive, err := h.restaurantService.IsRestaurantActive(*user.RestaurantID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error verifying restaurant status"})
			return
		}
		if !isActive {
			c.JSON(http.StatusForbidden, gin.H{"error": "your restaurant subscription is inactive or expired"})
			return
		}
	}

	// Resposta adaptada para o modelo multitenancy
	response := gin.H{
		"token": token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"type":  user.Type,
		},
	}

	// Adicionar info do restaurante se disponível
	if user.RestaurantID != nil {
		response["user"].(gin.H)["restaurant_id"] = user.RestaurantID

		restaurant, err := h.restaurantService.GetByID(*user.RestaurantID)
		if err == nil {
			response["restaurant"] = restaurant
		}
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Definir tipo de usuário
	userType := models.UserType(req.Type)

	// Verificar permissões para criação de usuários
	currentUserType, exists := c.Get("user_type")

	// Se não existir contexto de usuário, só permite criar staff
	if !exists && userType != models.UserTypeStaff {
		userType = models.UserTypeStaff
	} else if exists {
		// Verificar permissões com base no tipo do usuário atual
		currentType := currentUserType.(models.UserType)

		// SuperAdmin pode criar qualquer tipo
		if currentType == models.UserTypeSuperAdmin {
			// Permite o tipo solicitado
		} else if currentType == models.UserTypeAdmin {
			// Admin só pode criar manager ou staff
			if userType != models.UserTypeManager && userType != models.UserTypeStaff {
				userType = models.UserTypeStaff
			}
		} else if currentType == models.UserTypeManager {
			// Manager só pode criar staff
			userType = models.UserTypeStaff
		} else {
			// Staff não pode criar usuários
			c.JSON(http.StatusForbidden, gin.H{"error": "you don't have permission to create users"})
			return
		}
	}

	// Verificar associação com restaurante
	var restaurantID *uuid.UUID

	// Se for superadmin, não precisa de restaurante
	if userType == models.UserTypeSuperAdmin {
		restaurantID = nil
	} else {
		// Se um restaurant_id foi fornecido
		if req.RestaurantID != nil {
			restaurantID = req.RestaurantID
			fmt.Println("restaurant", *restaurantID)

			// Verificar se o restaurante existe
			_, err := h.restaurantService.GetByID(*restaurantID)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant not found"})
				return
			}

			// Se não for superadmin, verificar se o usuário atual pertence ao mesmo restaurante
			if exists && currentUserType.(models.UserType) != models.UserTypeSuperAdmin {
				currentRestaurantID, hasRestaurant := c.Get("restaurant_id")
				if !hasRestaurant || currentRestaurantID != *restaurantID {
					c.JSON(http.StatusForbidden, gin.H{"error": "you can only create users for your own restaurant"})
					return
				}
			}
		} else {
			// Se nenhum restaurant_id foi fornecido, usar o do usuário atual
			if exists && currentUserType.(models.UserType) != models.UserTypeSuperAdmin {
				currentRestaurantID, hasRestaurant := c.Get("restaurant_id")
				if hasRestaurant {
					restaurantID = new(uuid.UUID)
					*restaurantID = currentRestaurantID.(uuid.UUID)
				} else {
					c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id is required"})
					return
				}
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "restaurant_id is required"})
				return
			}
		}
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		Password:     req.Password,
		Type:         userType,
		RestaurantID: restaurantID,
	}

	if err := h.authService.Register(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":            user.ID,
		"name":          user.Name,
		"email":         user.Email,
		"type":          user.Type,
		"restaurant_id": user.RestaurantID,
	})
}

// RegisterSuperAdmin - cria o primeiro usuário superadmin
func (h *AuthHandler) RegisterSuperAdmin(c *gin.Context) {
	// Verificar se já existe um superadmin
	exists, err := h.authService.SuperAdminExists()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check superadmin existence"})
		return
	}

	if exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "superadmin already exists"})
		return
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		Password:     req.Password,
		Type:         models.UserTypeSuperAdmin,
		RestaurantID: nil,
	}

	if err := h.authService.Register(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"type":  user.Type,
	})
}

// GetProfile - obtém o perfil do usuário atual
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	restaurantIDValue, existsRestaurant := c.Get("restaurant_id")
	if !existsRestaurant {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Verifica se userID é um ponteiro para UUID ou uma string
	var userID uuid.UUID
	switch v := userIDValue.(type) {
	case *uuid.UUID:
		userID = *v
	case uuid.UUID:
		userID = v
	case string:
		var err error
		userID, err = uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id format"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id type"})
		return
	}

	// Verifica se restaurantID é um ponteiro para UUID ou uma string
	var restaurantID uuid.UUID
	switch v := restaurantIDValue.(type) {
	case *uuid.UUID:
		restaurantID = *v
	case uuid.UUID:
		restaurantID = v
	case string:
		var err error
		restaurantID, err = uuid.Parse(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant id format"})
			return
		}
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid restaurant id type"})
		return
	}

	user, err := h.authService.FindUserByID(restaurantID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// O resto do código permanece o mesmo
	response := gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"type":  user.Type,
	}

	if user.RestaurantID != nil {
		response["restaurant_id"] = user.RestaurantID

		restaurant, err := h.restaurantService.GetByID(*user.RestaurantID)
		if err == nil {
			response["restaurant"] = restaurant
		}
	}

	c.JSON(http.StatusOK, response)
}

// UpdateProfile - atualiza o perfil do usuário atual
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	restaurantID, existsRestaurant := c.Get("requested_restaurant_id")
	if !existsRestaurant {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.FindUserByID(restaurantID.(uuid.UUID), userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		user.Password = req.Password
	}

	if err := h.authService.UpdateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"type":  user.Type,
	})
}
