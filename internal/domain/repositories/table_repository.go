package repositories

import (
	"api-jet-manager/internal/domain/models"

	"github.com/google/uuid"
)

type TableRepository interface {
	Create(table *models.Table) error
	FindByID(id uuid.UUID) (*models.Table, error)
	FindByNumber(number int) (*models.Table, error)
	Update(table *models.Table) error
	Delete(id uuid.UUID) error
	List() ([]models.Table, error)
	UpdateStatus(id uuid.UUID, status models.TableStatus) error
	SetCurrentOrder(id uuid.UUID, orderID *uuid.UUID) error
}
