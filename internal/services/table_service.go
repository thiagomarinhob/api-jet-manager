// internal/services/table_service.go
package services

import (
	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/domain/repositories"

	"github.com/google/uuid"
)

type TableService struct {
	tableRepo repositories.TableRepository
}

func NewTableService(tableRepo repositories.TableRepository) *TableService {
	return &TableService{
		tableRepo: tableRepo,
	}
}

func (s *TableService) Create(table *models.Table) error {
	return s.tableRepo.Create(table)
}

func (s *TableService) GetByID(restaurant_id uuid.UUID, id uuid.UUID) (*models.Table, error) {
	return s.tableRepo.FindByID(restaurant_id, id)
}

func (s *TableService) GetByNumber(restaurant_id uuid.UUID, number int) (*models.Table, error) {
	return s.tableRepo.FindByNumber(restaurant_id, number)
}

func (s *TableService) Update(table *models.Table) error {
	return s.tableRepo.Update(table)
}

func (s *TableService) Delete(restaurant_id uuid.UUID, id uuid.UUID) error {
	return s.tableRepo.Delete(restaurant_id, id)
}

func (s *TableService) List(restaurant_id uuid.UUID) ([]models.Table, error) {
	return s.tableRepo.List(restaurant_id)
}

func (s *TableService) UpdateStatus(restaurant_id uuid.UUID, id uuid.UUID, status models.TableStatus) error {
	return s.tableRepo.UpdateStatus(restaurant_id, id, status)
}

func (s *TableService) SetCurrentOrder(restaurant_id uuid.UUID, id uuid.UUID, orderID *uuid.UUID) error {
	return s.tableRepo.SetCurrentOrder(restaurant_id, id, orderID)
}
