// internal/services/order_service.go
package services

import (
	"errors"
	"fmt"
	"time"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/domain/repositories"

	"github.com/google/uuid"
)

type OrderService struct {
	orderRepo   repositories.OrderRepository
	tableRepo   repositories.TableRepository
	financeRepo repositories.FinanceRepository
	productRepo repositories.ProductRepository
}

func NewOrderService(orderRepo repositories.OrderRepository, tableRepo repositories.TableRepository, financeRepo repositories.FinanceRepository, productRepo repositories.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		tableRepo:   tableRepo,
		financeRepo: financeRepo,
		productRepo: productRepo,
	}
}

func (s *OrderService) CreateOrder(order *models.Order, orderItems []models.OrderItem) error {
	// Calcular o valor total do pedido
	var totalAmount float64
	for _, item := range orderItems {
		totalAmount += item.Price * float64(item.Quantity)
	}
	order.TotalAmount = totalAmount

	// Criar o pedido
	if err := s.orderRepo.Create(order); err != nil {
		return err
	}

	// Adicionar os itens ao pedido
	for i := range orderItems {
		orderItems[i].OrderID = order.ID
		if err := s.orderRepo.AddItem(&orderItems[i]); err != nil {
			return err
		}
	}

	return nil
}

func (s *OrderService) GetByID(restaurant_id uuid.UUID, id uuid.UUID) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(restaurant_id, id)
	if err != nil {
		return nil, err
	}

	// Carregar os itens do pedido
	items, err := s.orderRepo.FindItems(restaurant_id, id)
	if err != nil {
		return nil, err
	}
	order.OrderItems = items

	return order, nil
}

func (s *OrderService) GetByTable(restaurant_id uuid.UUID, tableID uuid.UUID) ([]models.Order, error) {
	return s.orderRepo.FindByTable(restaurant_id, tableID)
}

func (s *OrderService) GetActiveByTable(restaurant_id uuid.UUID, tableID uuid.UUID) (*models.Order, error) {
	return s.orderRepo.FindActiveByTable(restaurant_id, tableID)
}

func (s *OrderService) GetByStatus(restaurant_id uuid.UUID, status models.OrderStatus) ([]models.Order, error) {
	return s.orderRepo.FindByStatus(restaurant_id, status)
}

func (s *OrderService) List(restaurant_id uuid.UUID) ([]models.Order, error) {
	return s.orderRepo.List(restaurant_id)
}

func (s *OrderService) UpdateStatus(restaurant_id uuid.UUID, id uuid.UUID, status models.OrderStatus) error {
	return s.orderRepo.UpdateStatus(restaurant_id, id, status)
}

func (s *OrderService) AddItem(restaurant_id uuid.UUID, item *models.OrderItem) error {
	// Verificar se o produto existe
	product, err := s.GetProductByID(restaurant_id, item.ProductID)
	if err != nil {
		return err
	}

	// Definir o preço do item de acordo com o preço atual do produto
	item.Price = product.Price

	// Adicionar o item
	if err := s.orderRepo.AddItem(item); err != nil {
		return err
	}

	// Atualizar o valor total do pedido
	order, err := s.orderRepo.FindByID(restaurant_id, item.OrderID)
	if err != nil {
		return err
	}

	order.TotalAmount += item.Price * float64(item.Quantity)
	return s.orderRepo.Update(order)
}

func (s *OrderService) RemoveItem(restaurant_id uuid.UUID, orderID, itemID uuid.UUID) error {
	// Encontrar o item
	items, err := s.orderRepo.FindItems(restaurant_id, orderID)
	if err != nil {
		return err
	}

	var itemToRemove *models.OrderItem
	for _, item := range items {
		if item.ID == itemID {
			itemToRemove = &item
			break
		}
	}

	if itemToRemove == nil {
		return errors.New("item not found")
	}

	// Remover o item
	if err := s.orderRepo.RemoveItem(restaurant_id, orderID, itemID); err != nil {
		return err
	}

	// Atualizar o valor total do pedido
	order, err := s.orderRepo.FindByID(restaurant_id, orderID)
	if err != nil {
		return err
	}

	order.TotalAmount -= itemToRemove.Price * float64(itemToRemove.Quantity)
	return s.orderRepo.Update(order)
}

func (s *OrderService) GetProductByID(restaurant_id uuid.UUID, id uuid.UUID) (*models.Product, error) {
	return s.productRepo.FindByID(restaurant_id, id)
}

func (s *OrderService) RegisterPayment(order *models.Order, userID uuid.UUID) error {
	// Registrar a transação financeira
	transaction := &models.FinancialTransaction{
		Type:        models.TransactionTypeIncome,
		Category:    models.TransactionCategorySales,
		Amount:      order.TotalAmount,
		Description: fmt.Sprintf("Payment for order #%d", order.ID),
		OrderID:     &order.ID,
		UserID:      userID,
		Date:        time.Now(),
	}

	return s.financeRepo.Create(transaction)
}

// FindDeliveryOrdersByDate retorna todos os pedidos de delivery para uma data específica
func (s *OrderService) FindDeliveryOrdersByDate(restaurantID uuid.UUID, date time.Time) ([]models.Order, error) {
	return s.orderRepo.FindDeliveryOrdersByDate(restaurantID, date)
}

// FindOrdersByDateAndType retorna todos os pedidos de um tipo específico para uma data
func (s *OrderService) FindOrdersByDateAndType(restaurantID uuid.UUID, date time.Time, orderType models.OrderType) ([]models.Order, error) {
	return s.orderRepo.FindOrdersByDateAndType(restaurantID, date, orderType)
}

// FindOrdersByDateRangeAndType retorna pedidos de um tipo específico dentro de um intervalo de datas
func (s *OrderService) FindOrdersByDateRangeAndType(restaurantID uuid.UUID, startDate, endDate time.Time, orderType models.OrderType) ([]models.Order, error) {
	return s.orderRepo.FindOrdersByDateRangeAndType(restaurantID, startDate, endDate, orderType)
}

// FindTodayDeliveryOrders retorna todos os pedidos de delivery para o dia atual
func (s *OrderService) FindTodayDeliveryOrders(restaurantID uuid.UUID) ([]models.Order, error) {
	today := time.Now()
	return s.orderRepo.FindDeliveryOrdersByDate(restaurantID, today)
}

// CountDeliveryOrdersByDate conta o número de pedidos de delivery em uma data específica
// func (s *OrderService) CountDeliveryOrdersByDate(restaurantID uuid.UUID, date time.Time) (int64, error) {
// 	orders, err := s.orderRepo.FindDeliveryOrdersByDate(restaurantID, date)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return int64(len(orders)), nil
// }
