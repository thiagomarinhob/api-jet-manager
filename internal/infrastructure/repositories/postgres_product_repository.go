// internal/infrastructure/repositories/postgres_product_repository.go
package repositories

import (
	"errors"
	"fmt"
	"strings"

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

func (r *PostgresProductRepository) FindByID(restaurantID, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	if err := r.DB.Where("restaurant_id = ? AND id = ?", restaurantID, id).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

func (r *PostgresProductRepository) Update(product *models.Product) error {
	// Assumindo que o restaurant_id já está definido no objeto product
	return r.DB.Save(product).Error
}

func (r *PostgresProductRepository) Delete(restaurantID, id uuid.UUID) error {
	return r.DB.Where("restaurant_id = ?", restaurantID).Delete(&models.Product{}, "id = ?", id).Error
}

func (r *PostgresProductRepository) FindByRestaurant(restaurantID uuid.UUID) ([]models.Product, error) {
	var products []models.Product
	if err := r.DB.Where("restaurant_id = ?", restaurantID).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *PostgresProductRepository) FindByCategory(restaurantID uuid.UUID, category models.ProductCategory) ([]models.Product, error) {
	var products []models.Product
	if err := r.DB.Where("restaurant_id = ? AND category = ?", restaurantID, category).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *PostgresProductRepository) UpdateStock(restaurantID, id uuid.UUID, inStock bool) error {
	return r.DB.Model(&models.Product{}).Where("restaurant_id = ? AND id = ?", restaurantID, id).Update("in_stock", inStock).Error
}

func (r *PostgresProductRepository) FindWithFilters(
	restaurantID uuid.UUID,
	offset int,
	limit int,
	category *models.ProductCategory,
	inStock *bool,
	nameSearch string,
	sortBy string,
	sortOrder string,
) ([]models.Product, int64, error) {
	// Construir a query base
	query := r.DB.Model(&models.Product{}).Where("restaurant_id = ?", restaurantID)

	// Pré-carregar a categoria
	query = query.Preload("Category")

	// Aplicar filtros se fornecidos
	if category != nil {
		query = query.Where("category = ?", *category)
	}

	if inStock != nil {
		query = query.Where("in_stock = ?", *inStock)
	}

	if nameSearch != "" {
		query = query.Where("name ILIKE ?", "%"+nameSearch+"%")
	}

	// Contar total de itens antes de aplicar paginação
	var totalItems int64
	if err := query.Count(&totalItems).Error; err != nil {
		return nil, 0, err
	}

	// Validar e aplicar ordenação
	if isValidSortField(sortBy) {
		// Garantir que sortOrder é válido
		if strings.ToLower(sortOrder) != "asc" && strings.ToLower(sortOrder) != "desc" {
			sortOrder = "asc"
		}

		query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))
	} else {
		// Ordenação padrão se o campo for inválido
		query = query.Order("name asc")
	}

	// Aplicar paginação
	query = query.Offset(offset).Limit(limit)

	// Executar a consulta
	var products []models.Product
	if err := query.Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, totalItems, nil
}

func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"name":       true,
		"price":      true,
		"category":   true,
		"in_stock":   true,
		"created_at": true,
		"updated_at": true,
	}

	return validFields[field]
}
