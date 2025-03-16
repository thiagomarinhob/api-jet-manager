package repositories

import (
	"api-jet-manager/internal/domain/models"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(order *models.Order) error
	FindByID(id uuid.UUID) (*models.Order, error)
	Update(order *models.Order) error
	Delete(id uuid.UUID) error
	List() ([]models.Order, error)
	FindByTable(tableID uuid.UUID) ([]models.Order, error)
	FindActiveByTable(tableID uuid.UUID) (*models.Order, error)
	FindByStatus(status models.OrderStatus) ([]models.Order, error)
	UpdateStatus(id uuid.UUID, status models.OrderStatus) error
	AddItem(item *models.OrderItem) error
	RemoveItem(orderID uuid.UUID, itemID uuid.UUID) error
	UpdateItem(item *models.OrderItem) error
	FindItems(orderID uuid.UUID) ([]models.OrderItem, error)
}
