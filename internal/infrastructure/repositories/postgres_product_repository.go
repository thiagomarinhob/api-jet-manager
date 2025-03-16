// internal/infrastructure/repositories/postgres_product_repository.go
package repositories

import (
	"errors"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/infrastructure/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresProductRepository struct {
	DB *gorm.DB
}

func NewPostgresProductRepository(db *database.PostgresDB) *PostgresProductRepository {
	return &PostgresProductRepository{
		DB: db.DB,
	}
}

func (r *PostgresProductRepository) Create(product *models.Product) error {
	return r.DB.Create(product).Error
}

func (r *PostgresProductRepository) FindByID(id uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := r.DB.Where("id = ?", id).First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (r *PostgresProductRepository) Update(product *models.Product) error {
	return r.DB.Save(product).Error
}

func (r *PostgresProductRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.Product{}, id).Error
}

func (r *PostgresProductRepository) List() ([]models.Product, error) {
	var products []models.Product
	if err := r.DB.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *PostgresProductRepository) FindByCategory(category models.ProductCategory) ([]models.Product, error) {
	var products []models.Product
	if err := r.DB.Where("category = ?", category).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *PostgresProductRepository) UpdateStock(id uuid.UUID, inStock bool) error {
	return r.DB.Model(&models.Product{}).Where("id = ?", id).Update("in_stock", inStock).Error
}
