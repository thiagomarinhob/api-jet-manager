package services

import (
	"errors"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/domain/repositories"

	"github.com/google/uuid"
)

type ProductCategoryService struct {
	categoryRepo repositories.ProductCategoryRepository
}

func NewProductCategoryService(categoryRepo repositories.ProductCategoryRepository) *ProductCategoryService {
	return &ProductCategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *ProductCategoryService) Create(category *models.ProductCategory) error {
	// Verificar se já existe categoria com o mesmo nome no restaurante
	existing, err := s.categoryRepo.FindByName(category.RestaurantID, category.Name)
	if err != nil {
		return err
	}

	if existing != nil {
		return errors.New("a category with this name already exists")
	}

	return s.categoryRepo.Create(category)
}

func (s *ProductCategoryService) FindByID(restaurant_id uuid.UUID, id uuid.UUID) (*models.ProductCategory, error) {
	return s.categoryRepo.FindByID(restaurant_id, id)
}

func (s *ProductCategoryService) Update(restaurant_id uuid.UUID, category *models.ProductCategory) error {
	existing, err := s.categoryRepo.FindByID(restaurant_id, category.ID)
	if err != nil {
		return err
	}

	// Verificar se já existe outra categoria com o mesmo nome
	if category.Name != existing.Name {
		nameExists, err := s.categoryRepo.FindByName(category.RestaurantID, category.Name)
		if err != nil {
			return err
		}

		if nameExists != nil {
			return errors.New("a category with this name already exists")
		}
	}

	return s.categoryRepo.Update(category)
}

func (s *ProductCategoryService) Delete(restaurant_id uuid.UUID, id uuid.UUID) error {
	// Em vez de deletar completamente, apenas desativa a categoria
	return s.categoryRepo.Delete(restaurant_id, id)
}

func (s *ProductCategoryService) FindByRestaurant(restaurantID uuid.UUID) ([]models.ProductCategory, error) {
	return s.categoryRepo.FindByRestaurant(restaurantID)
}

func (s *ProductCategoryService) FindActive(restaurantID uuid.UUID) ([]models.ProductCategory, error) {
	return s.categoryRepo.FindActive(restaurantID)
}

func (s *ProductCategoryService) FindByName(restaurantID uuid.UUID, name string) (*models.ProductCategory, error) {
	return s.categoryRepo.FindByName(restaurantID, name)
}

// ListWithPagination retorna categorias paginadas com opções de filtragem e ordenação
func (s *ProductCategoryService) FindWithFilters(
	restaurantID uuid.UUID,
	page int,
	pageSize int,
	active *bool,
	nameSearch string,
	sortBy string,
	sortOrder string,
) ([]models.ProductCategory, int64, error) {
	// Calcular offset
	offset := (page - 1) * pageSize

	// Chamar o repositório para buscar as categorias paginadas
	return s.categoryRepo.FindWithFilters(
		restaurantID,
		offset,
		pageSize,
		active,
		nameSearch,
		sortBy,
		sortOrder,
	)
}
