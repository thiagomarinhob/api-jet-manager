// internal/services/product_service.go
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

func (s *ProductService) List() ([]models.Product, error) {
	return s.productRepo.List()
}

func (s *ProductService) GetByCategory(category models.ProductCategory) ([]models.Product, error) {
	return s.productRepo.FindByCategory(category)
}

func (s *ProductService) UpdateStock(id uuid.UUID, inStock bool) error {
	return s.productRepo.UpdateStock(id, inStock)
}
