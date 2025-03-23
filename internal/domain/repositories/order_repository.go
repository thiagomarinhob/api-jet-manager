package repositories

import (
	"api-jet-manager/internal/domain/models"
	"time"

	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(order *models.Order) error
	FindByID(restaurantID, id uuid.UUID) (*models.Order, error)
	Update(order *models.Order) error
	Delete(restaurantID, id uuid.UUID) error
	List(restaurantID uuid.UUID) ([]models.Order, error)
	FindByTable(restaurantID, tableID uuid.UUID) ([]models.Order, error)
	FindActiveByTable(restaurantID, tableID uuid.UUID) (*models.Order, error)
	FindByStatus(restaurantID uuid.UUID, status models.OrderStatus) ([]models.Order, error)
	UpdateStatus(restaurantID, id uuid.UUID, status models.OrderStatus) error
	AddItem(item *models.OrderItem) error
	RemoveItem(restaurantID, orderID, itemID uuid.UUID) error
	UpdateItem(item *models.OrderItem) error
	FindItems(restaurantID, orderID uuid.UUID) ([]models.OrderItem, error)
	FindOrdersByDateAndType(restaurantID uuid.UUID, date time.Time, orderType models.OrderType) ([]models.Order, error)
	FindDeliveryOrdersByDate(restaurantID uuid.UUID, date time.Time) ([]models.Order, error)
	FindOrdersByDateRangeAndType(restaurantID uuid.UUID, startDate, endDate time.Time, orderType models.OrderType) ([]models.Order, error)
}
