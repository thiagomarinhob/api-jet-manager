package repositories

import (
	"api-jet-manager/internal/domain/models"

	"github.com/google/uuid"
)

type ProductCategoryRepository interface {
	Create(category *models.ProductCategory) error
	FindByID(restaurantID, id uuid.UUID) (*models.ProductCategory, error)
	Update(category *models.ProductCategory) error
	Delete(restaurantID, id uuid.UUID) error

	// Métodos de listagem
	FindByRestaurant(restaurantID uuid.UUID) ([]models.ProductCategory, error)
	FindActive(restaurantID uuid.UUID) ([]models.ProductCategory, error)

	// Busca por nome
	FindByName(restaurantID uuid.UUID, name string) (*models.ProductCategory, error)

	// Método de paginação e filtragem
	FindWithFilters(
		restaurantID uuid.UUID,
		offset int,
		limit int,
		active *bool,
		nameSearch string,
		sortBy string,
		sortOrder string,
	) ([]models.ProductCategory, int64, error)
}
