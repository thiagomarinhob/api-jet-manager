// internal/api/middlewares/auth_middleware.go
package middlewares

import (
	"net/http"
	"strings"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/infrastructure/auth"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService *auth.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			return
		}

		tokenString := parts[1]
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Armazena informações do usuário no contexto para uso posterior
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("user_type", claims.UserType)
		c.Set("restaurant_id", claims.RestaurantID)

		c.Next()
	}
}

// Middleware para verificar se o usuário tem um tipo específico
func UserTypeMiddleware(allowedTypes ...models.UserType) gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userTypeValue, ok := userType.(models.UserType)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		// Superadmin sempre tem acesso
		if userTypeValue == models.UserTypeSuperAdmin {
			c.Next()
			return
		}

		// Verifica se o tipo do usuário está entre os permitidos
		for _, allowedType := range allowedTypes {
			if userTypeValue == allowedType {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}

// Middleware que restringe acesso apenas para superadmins
func SuperAdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userType, exists := c.Get("user_type")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		userTypeValue, ok := userType.(models.UserType)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		if userTypeValue != models.UserTypeSuperAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "only superadmins can access this resource"})
			return
		}

		c.Next()
	}
}
