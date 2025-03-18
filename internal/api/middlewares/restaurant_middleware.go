package middlewares

import (
	"fmt"
	"net/http"

	"api-jet-manager/internal/domain/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RestaurantMiddleware garante que usuários só possam acessar dados do seu próprio restaurante
func RestaurantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtém o user_id e type do contexto (definidos pelo AuthMiddleware)
		_, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userType, exists := c.Get("user_type")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Superadmin pode acessar qualquer restaurante
		if userType.(models.UserType) == models.UserTypeSuperAdmin {
			c.Next()
			return
		}

		// Para outros usuários, verifica o restaurante associado
		restaurantID, exists := c.Get("restaurant_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Converte o restaurantID para string, independente do tipo real
		var restaurantIDStr string

		// Tenta diferentes conversões possíveis
		switch v := restaurantID.(type) {
		case string:
			restaurantIDStr = v
		case uuid.UUID:
			restaurantIDStr = v.String()
		case fmt.Stringer: // Para qualquer tipo que implementa String()
			restaurantIDStr = v.String()
		default:
			// Para qualquer outro tipo, tenta imprimir como string
			restaurantIDStr = fmt.Sprintf("%v", restaurantID)
		}

		// Verifica se o ID do restaurante está na URL
		requestedRestaurantID := c.Param("restaurant_id")
		if requestedRestaurantID == "" {
			// Se não estiver na URL, tenta obter do query parameter
			requestedRestaurantID = c.Query("restaurant_id")
		}

		// Se não encontrar o restaurant_id nem na URL nem nos query parameters,
		// verifica no body da requisição para métodos POST e PUT
		if requestedRestaurantID == "" && (c.Request.Method == "POST" || c.Request.Method == "PUT") {
			var requestBody map[string]interface{}
			if err := c.ShouldBindJSON(&requestBody); err == nil {
				if id, ok := requestBody["restaurant_id"].(string); ok {
					requestedRestaurantID = id
				}
				// Restaura o body para que possa ser lido novamente pelos handlers
				body, _ := c.Request.GetBody()
				c.Request.Body = body
			}
		}

		// Se ainda não encontrou o restaurant_id, usa o do usuário
		if requestedRestaurantID == "" {
			c.Set("requested_restaurant_id", restaurantIDStr)
			c.Next()
			return
		}

		// Verifica se o usuário pertence ao restaurante requisitado
		if restaurantIDStr != requestedRestaurantID {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "you don't have access to this restaurant"})
			return
		}

		// Armazena o ID do restaurante requisitado no contexto
		c.Set("requested_restaurant_id", requestedRestaurantID)
		c.Next()
	}
}

// Middleware para verificar permissões específicas em um restaurante
func RestaurantRoleMiddleware(allowedTypes ...models.UserType) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Obtém o tipo de usuário do contexto
		userType, exists := c.Get("user_type")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		// Superadmin sempre tem acesso
		if userType.(models.UserType) == models.UserTypeSuperAdmin {
			c.Next()
			return
		}

		// Verifica se o tipo do usuário está entre os permitidos
		userTypeStr := userType.(models.UserType)
		for _, allowedType := range allowedTypes {
			if userTypeStr == allowedType {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}
