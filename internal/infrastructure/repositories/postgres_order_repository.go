// internal/infrastructure/repositories/postgres_order_repository.go
package repositories

import (
	"errors"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/infrastructure/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresOrderRepository struct {
	DB *gorm.DB
}

func NewPostgresOrderRepository(db *database.PostgresDB) *PostgresOrderRepository {
	return &PostgresOrderRepository{
		DB: db.DB,
	}
}

func (r *PostgresOrderRepository) Create(order *models.Order) error {
	return r.DB.Create(order).Error
}

func (r *PostgresOrderRepository) FindByID(id uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.DB.Where("id = ?", id).First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

func (r *PostgresOrderRepository) Update(order *models.Order) error {
	return r.DB.Save(order).Error
}

func (r *PostgresOrderRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.Order{}, id).Error
}

func (r *PostgresOrderRepository) List() ([]models.Order, error) {
	var orders []models.Order
	if err := r.DB.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *PostgresOrderRepository) FindByTable(tableID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	if err := r.DB.Where("table_id = ?", tableID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *PostgresOrderRepository) FindActiveByTable(tableID uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.DB.Where("table_id = ? AND status NOT IN (?, ?)",
		tableID, models.OrderStatusPaid, models.OrderStatusCancelled).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no active order for this table")
		}
		return nil, err
	}
	return &order, nil
}

func (r *PostgresOrderRepository) FindByStatus(status models.OrderStatus) ([]models.Order, error) {
	var orders []models.Order
	if err := r.DB.Where("status = ?", status).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *PostgresOrderRepository) UpdateStatus(id uuid.UUID, status models.OrderStatus) error {
	return r.DB.Model(&models.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *PostgresOrderRepository) AddItem(item *models.OrderItem) error {
	return r.DB.Create(item).Error
}

func (r *PostgresOrderRepository) RemoveItem(orderID uuid.UUID, itemID uuid.UUID) error {
	return r.DB.Where("order_id = ? AND id = ?", orderID, itemID).Delete(&models.OrderItem{}).Error
}

func (r *PostgresOrderRepository) UpdateItem(item *models.OrderItem) error {
	return r.DB.Save(item).Error
}

func (r *PostgresOrderRepository) FindItems(orderID uuid.UUID) ([]models.OrderItem, error) {
	var items []models.OrderItem
	if err := r.DB.Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
