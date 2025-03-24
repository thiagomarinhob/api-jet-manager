// internal/domain/repositories/product_repository.go
package repositories

import (
	"api-jet-manager/internal/domain/models"

	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindByID(restaurantID, id uuid.UUID) (*models.Product, error)
	Update(product *models.Product) error
	Delete(restaurantID, id uuid.UUID) error

	// Métodos de listagem simples
	FindByRestaurant(restaurantID uuid.UUID) ([]models.Product, error)
	FindByCategory(restaurantID uuid.UUID, category models.ProductCategory) ([]models.Product, error)
	UpdateStock(restaurantID, id uuid.UUID, inStock bool) error

	// Método de paginação e filtragem
	// Retorna: produtos, contagem total e erro
	FindWithFilters(
		restaurantID uuid.UUID,
		offset int,
		limit int,
		category *models.ProductCategory,
		inStock *bool,
		nameSearch string,
		sortBy string,
		sortOrder string,
	) ([]models.Product, int64, error)
}
