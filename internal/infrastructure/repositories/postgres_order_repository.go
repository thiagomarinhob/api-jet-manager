// internal/infrastructure/repositories/postgres_order_repository.go
package repositories

import (
	"errors"
	"fmt"
	"time"

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
	// O restaurant_id já deve estar definido no objeto order antes de chamar este método
	return r.DB.Create(order).Error
}

func (r *PostgresOrderRepository) FindByID(restaurantID, id uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.DB.Where("restaurant_id = ? AND id = ?", restaurantID, id).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}

func (r *PostgresOrderRepository) Update(order *models.Order) error {
	// Assumindo que o restaurant_id já está definido no objeto order
	return r.DB.Save(order).Error
}

func (r *PostgresOrderRepository) Delete(restaurantID, id uuid.UUID) error {
	return r.DB.Where("restaurant_id = ?", restaurantID).Delete(&models.Order{}, id).Error
}

func (r *PostgresOrderRepository) List(restaurantID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	if err := r.DB.Where("restaurant_id = ?", restaurantID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *PostgresOrderRepository) FindByTable(restaurantID, tableID uuid.UUID) ([]models.Order, error) {
	var orders []models.Order
	if err := r.DB.Where("restaurant_id = ? AND table_id = ?", restaurantID, tableID).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *PostgresOrderRepository) FindActiveByTable(restaurantID, tableID uuid.UUID) (*models.Order, error) {
	var order models.Order
	if err := r.DB.Where("restaurant_id = ? AND table_id = ? AND status NOT IN (?, ?)",
		restaurantID, tableID, models.OrderStatusPaid, models.OrderStatusCancelled).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no active order for this table")
		}
		return nil, err
	}
	return &order, nil
}

func (r *PostgresOrderRepository) FindByStatus(restaurantID uuid.UUID, status models.OrderStatus) ([]models.Order, error) {
	var orders []models.Order
	if err := r.DB.Where("restaurant_id = ? AND status = ?", restaurantID, status).Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *PostgresOrderRepository) UpdateStatus(restaurantID, id uuid.UUID, status models.OrderStatus) error {
	return r.DB.Model(&models.Order{}).Where("restaurant_id = ? AND id = ?", restaurantID, id).Update("status", status).Error
}

func (r *PostgresOrderRepository) AddItem(item *models.OrderItem) error {
	// Ao adicionar um item, precisamos garantir que o order_id pertence ao restaurante correto
	// Isso geralmente é feito no service layer antes de chamar este método
	return r.DB.Create(item).Error
}

func (r *PostgresOrderRepository) RemoveItem(restaurantID, orderID, itemID uuid.UUID) error {
	// Primeiro, verifique se o pedido pertence ao restaurante
	var order models.Order
	if err := r.DB.Where("restaurant_id = ? AND id = ?", restaurantID, orderID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("order not found for this restaurant")
		}
		return err
	}

	// Se o pedido pertence ao restaurante, então podemos remover o item
	return r.DB.Where("order_id = ? AND id = ?", orderID, itemID).Delete(&models.OrderItem{}).Error
}

func (r *PostgresOrderRepository) UpdateItem(item *models.OrderItem) error {
	// Assumindo que garantimos no service layer que este item pertence a um pedido do restaurante correto
	return r.DB.Save(item).Error
}

func (r *PostgresOrderRepository) FindItems(restaurantID, orderID uuid.UUID) ([]models.OrderItem, error) {
	// Primeiro, verifique se o pedido pertence ao restaurante
	var order models.Order
	if err := r.DB.Where("restaurant_id = ? AND id = ?", restaurantID, orderID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("order not found for this restaurant")
		}
		return nil, err
	}

	// Se o pedido pertence ao restaurante, então podemos buscar os itens
	var items []models.OrderItem
	if err := r.DB.Where("order_id = ?", orderID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// FindOrdersByDateAndType busca pedidos por data e tipo específico
// Este método já estava implementado corretamente para multitenancy
func (r *PostgresOrderRepository) FindOrdersByDateAndType(restaurantID uuid.UUID, date time.Time, orderType models.OrderType) ([]models.Order, error) {
	var orders []models.Order

	// Cria o início e fim do dia para filtrar pedidos
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Query buscando pedidos pelo restaurante, tipo e intervalo de data
	result := r.DB.
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Preload("Restaurant").
		Preload("User").
		Where("restaurant_id = ? AND type = ? AND created_at >= ? AND created_at < ?",
			restaurantID, orderType, startOfDay, endOfDay).
		Find(&orders)

	if result.Error != nil {
		return nil, fmt.Errorf("erro ao buscar pedidos por data e tipo: %w", result.Error)
	}

	return orders, nil
}

// FindDeliveryOrdersByDate é um método helper específico para pedidos de delivery
// Este método já estava implementado corretamente para multitenancy
func (r *PostgresOrderRepository) FindDeliveryOrdersByDate(restaurantID uuid.UUID, date time.Time) ([]models.Order, error) {
	return r.FindOrdersByDateAndType(restaurantID, date, models.OrderTypeDelivery)
}

// FindOrdersByDateRangeAndType busca pedidos em um intervalo de datas e por tipo
// Este método já estava implementado corretamente para multitenancy
func (r *PostgresOrderRepository) FindOrdersByDateRangeAndType(restaurantID uuid.UUID, startDate, endDate time.Time, orderType models.OrderType) ([]models.Order, error) {
	var orders []models.Order

	// Ajusta horários para início e fim dos dias
	startOfDay := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, startDate.Location())
	endOfDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 999999999, endDate.Location())

	// Query com range de datas e restaurante
	result := r.DB.
		Preload("OrderItems").
		Preload("OrderItems.Product").
		Preload("Restaurant").
		Preload("User").
		Where("restaurant_id = ? AND type = ? AND created_at >= ? AND created_at <= ?",
			restaurantID, orderType, startOfDay, endOfDay).
		Find(&orders)

	if result.Error != nil {
		return nil, fmt.Errorf("erro ao buscar pedidos no intervalo de datas: %w", result.Error)
	}

	return orders, nil
}
