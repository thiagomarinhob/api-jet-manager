package repositories

import (
	"api-jet-manager/internal/domain/models"

	"github.com/google/uuid"
)

type TableRepository interface {
	Create(table *models.Table) error
	FindByID(restauranteID, id uuid.UUID) (*models.Table, error)
	FindByNumber(restauranteID uuid.UUID, number int) (*models.Table, error)
	Update(table *models.Table) error
	Delete(restauranteID, id uuid.UUID) error
	List(restauranteID uuid.UUID) ([]models.Table, error)
	UpdateStatus(restauranteID, id uuid.UUID, status models.TableStatus) error
	SetCurrentOrder(restauranteID, id uuid.UUID, orderID *uuid.UUID) error
}
