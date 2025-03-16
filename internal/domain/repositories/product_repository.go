// internal/domain/repositories/product_repository.go
package repositories

import (
	"api-jet-manager/internal/domain/models"

	"github.com/google/uuid"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindByID(id uuid.UUID) (*models.Product, error)
	Update(product *models.Product) error
	Delete(id uuid.UUID) error
	List() ([]models.Product, error)
	FindByCategory(category models.ProductCategory) ([]models.Product, error)
	UpdateStock(id uuid.UUID, inStock bool) error
}
