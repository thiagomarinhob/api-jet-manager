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

type PostgresProductCategoryRepository struct {
	DB *gorm.DB
}

func NewPostgresProductCategoryRepository(db *database.PostgresDB) *PostgresProductCategoryRepository {
	return &PostgresProductCategoryRepository{
		DB: db.DB,
	}
}

func (r *PostgresProductCategoryRepository) Create(category *models.ProductCategory) error {
	// O restaurant_id já deve estar definido no objeto category antes de chamar este método
	return r.DB.Create(category).Error
}

func (r *PostgresProductCategoryRepository) FindByID(restaurantID, id uuid.UUID) (*models.ProductCategory, error) {
	var category models.ProductCategory
	if err := r.DB.Where("restaurant_id = ? AND id = ?", restaurantID, id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, err
	}
	return &category, nil
}

func (r *PostgresProductCategoryRepository) Update(category *models.ProductCategory) error {
	// Assumindo que o restaurant_id já está definido no objeto category
	return r.DB.Save(category).Error
}

func (r *PostgresProductCategoryRepository) Delete(restaurantID, id uuid.UUID) error {
	return r.DB.Model(&models.ProductCategory{}).Where("restaurant_id = ? AND id = ?", restaurantID, id).Update("active", false).Error
}

func (r *PostgresProductCategoryRepository) FindByRestaurant(restaurantID uuid.UUID) ([]models.ProductCategory, error) {
	var categories []models.ProductCategory
	if err := r.DB.Where("restaurant_id = ?", restaurantID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *PostgresProductCategoryRepository) FindActive(restaurantID uuid.UUID) ([]models.ProductCategory, error) {
	var categories []models.ProductCategory
	if err := r.DB.Where("restaurant_id = ? AND active = true", restaurantID).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *PostgresProductCategoryRepository) FindByName(restaurantID uuid.UUID, name string) (*models.ProductCategory, error) {
	var category models.ProductCategory
	if err := r.DB.Where("restaurant_id = ? AND name = ?", restaurantID, name).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Retorna nil sem erro se não encontrar, para verificação de unicidade
		}
		return nil, err
	}
	return &category, nil
}

func (r *PostgresProductCategoryRepository) FindWithFilters(
	restaurantID uuid.UUID,
	offset int,
	limit int,
	active *bool,
	nameSearch string,
	sortBy string,
	sortOrder string,
) ([]models.ProductCategory, int64, error) {
	// Construir a query base
	query := r.DB.Model(&models.ProductCategory{}).Where("restaurant_id = ?", restaurantID)

	// Aplicar filtros se fornecidos
	if active != nil {
		query = query.Where("active = ?", *active)
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
	if isValidCategorySortField(sortBy) {
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
	var categories []models.ProductCategory
	if err := query.Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	return categories, totalItems, nil
}

// Função auxiliar para validar campos de ordenação
func isValidCategorySortField(field string) bool {
	validFields := map[string]bool{
		"name":       true,
		"active":     true,
		"created_at": true,
		"updated_at": true,
	}

	return validFields[field]
}
