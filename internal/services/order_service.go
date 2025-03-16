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

func (s *OrderService) GetByID(id uuid.UUID) (*models.Order, error) {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Carregar os itens do pedido
	items, err := s.orderRepo.FindItems(id)
	if err != nil {
		return nil, err
	}
	order.OrderItems = items

	return order, nil
}

func (s *OrderService) GetByTable(tableID uuid.UUID) ([]models.Order, error) {
	return s.orderRepo.FindByTable(tableID)
}

func (s *OrderService) GetActiveByTable(tableID uuid.UUID) (*models.Order, error) {
	return s.orderRepo.FindActiveByTable(tableID)
}

func (s *OrderService) GetByStatus(status models.OrderStatus) ([]models.Order, error) {
	return s.orderRepo.FindByStatus(status)
}

func (s *OrderService) List() ([]models.Order, error) {
	return s.orderRepo.List()
}

func (s *OrderService) UpdateStatus(id uuid.UUID, status models.OrderStatus) error {
	return s.orderRepo.UpdateStatus(id, status)
}

func (s *OrderService) AddItem(item *models.OrderItem) error {
	// Verificar se o produto existe
	product, err := s.GetProductByID(item.ProductID)
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
	order, err := s.orderRepo.FindByID(item.OrderID)
	if err != nil {
		return err
	}

	order.TotalAmount += item.Price * float64(item.Quantity)
	return s.orderRepo.Update(order)
}

func (s *OrderService) RemoveItem(orderID, itemID uuid.UUID) error {
	// Encontrar o item
	items, err := s.orderRepo.FindItems(orderID)
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
	if err := s.orderRepo.RemoveItem(orderID, itemID); err != nil {
		return err
	}

	// Atualizar o valor total do pedido
	order, err := s.orderRepo.FindByID(orderID)
	if err != nil {
		return err
	}

	order.TotalAmount -= itemToRemove.Price * float64(itemToRemove.Quantity)
	return s.orderRepo.Update(order)
}

func (s *OrderService) GetProductByID(id uuid.UUID) (*models.Product, error) {
	return s.productRepo.FindByID(id)
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
