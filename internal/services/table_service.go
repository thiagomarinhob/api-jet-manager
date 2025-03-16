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

func (s *TableService) GetByID(id uuid.UUID) (*models.Table, error) {
	return s.tableRepo.FindByID(id)
}

func (s *TableService) GetByNumber(number int) (*models.Table, error) {
	return s.tableRepo.FindByNumber(number)
}

func (s *TableService) Update(table *models.Table) error {
	return s.tableRepo.Update(table)
}

func (s *TableService) Delete(id uuid.UUID) error {
	return s.tableRepo.Delete(id)
}

func (s *TableService) List() ([]models.Table, error) {
	return s.tableRepo.List()
}

func (s *TableService) UpdateStatus(id uuid.UUID, status models.TableStatus) error {
	return s.tableRepo.UpdateStatus(id, status)
}

func (s *TableService) SetCurrentOrder(id uuid.UUID, orderID *uuid.UUID) error {
	return s.tableRepo.SetCurrentOrder(id, orderID)
}
