package services

import (
	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/domain/repositories"

	"github.com/google/uuid"
)

type ProductService struct {
	productRepo repositories.ProductRepository
}

func NewProductService(productRepo repositories.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) Create(product *models.Product) error {
	return s.productRepo.Create(product)
}

func (s *ProductService) GetByID(id uuid.UUID) (*models.Product, error) {
	return s.productRepo.FindByID(id)
}

func (s *ProductService) Update(product *models.Product) error {
	return s.productRepo.Update(product)
}

func (s *ProductService) Delete(id uuid.UUID) error {
	return s.productRepo.Delete(id)
}

func (s *ProductService) List(restaurantID uuid.UUID) ([]models.Product, error) {
	return s.productRepo.FindByRestaurant(restaurantID)
}

func (s *ProductService) GetByCategory(restaurantID uuid.UUID, category models.ProductCategory) ([]models.Product, error) {
	return s.productRepo.FindByCategory(restaurantID, category)
}

func (s *ProductService) UpdateStock(id uuid.UUID, inStock bool) error {
	return s.productRepo.UpdateStock(id, inStock)
}

// ListWithPagination retorna produtos paginados com opções de filtragem e ordenação
func (s *ProductService) ListWithPagination(
	restaurantID uuid.UUID,
	page int,
	pageSize int,
	category *models.ProductCategory,
	inStock *bool,
	nameSearch string,
	sortBy string,
	sortOrder string,
) ([]models.Product, int64, error) {
	// Calcular offset
	offset := (page - 1) * pageSize

	// Chamar o repositório para buscar os produtos paginados
	return s.productRepo.FindWithFilters(
		restaurantID,
		offset,
		pageSize,
		category,
		inStock,
		nameSearch,
		sortBy,
		sortOrder,
	)
}
