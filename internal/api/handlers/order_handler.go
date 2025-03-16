package handlers

import (
	"net/http"
	"time"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type OrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
	Notes     string    `json:"notes"`
}

type OrderRequest struct {
	TableID    *uuid.UUID         `json:"table_id"`
	OrderItems []OrderItemRequest `json:"order_items" binding:"required,dive"`
	Notes      string             `json:"notes"`
}

type OrderHandler struct {
	orderService *services.OrderService
	tableService *services.TableService
}

func NewOrderHandler(orderService *services.OrderService, tableService *services.TableService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		tableService: tableService,
	}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Verificar se a mesa existe e está livre
	if req.TableID == nil {
		table, err := h.tableService.GetByID(*req.TableID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "table not found"})
			return
		}

		if table.Status != models.TableStatusFree && table.Status != models.TableStatusReserved {
			c.JSON(http.StatusBadRequest, gin.H{"error": "table is not available"})
			return
		}
	}

	order := &models.Order{
		TableID: req.TableID,
		UserID:  userID.(uuid.UUID),
		Status:  models.OrderStatusPending,
		Notes:   req.Notes,
	}

	// Processar itens do pedido
	orderItems := make([]models.OrderItem, 0, len(req.OrderItems))
	for _, item := range req.OrderItems {
		product, err := h.orderService.GetProductByID(item.ProductID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product not found: " + item.ProductID.String()})
			return
		}

		if !product.InStock {
			c.JSON(http.StatusBadRequest, gin.H{"error": "product out of stock: " + product.Name})
			return
		}

		orderItems = append(orderItems, models.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
			Notes:     item.Notes,
		})
	}

	if err := h.orderService.CreateOrder(order, orderItems); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Atualizar status da mesa se o pedido estiver associado a uma mesa
	if req.TableID == nil {
		if err := h.tableService.UpdateStatus(*req.TableID, models.TableStatusOccupied); err != nil {
			// Log do erro, mas não falha a criação do pedido
			c.JSON(http.StatusCreated, gin.H{
				"order":   order,
				"warning": "order created but failed to update table status",
			})
			return
		}

		if err := h.tableService.SetCurrentOrder(*req.TableID, &order.ID); err != nil {
			c.JSON(http.StatusCreated, gin.H{
				"order":   order,
				"warning": "order created but failed to link order to table",
			})
			return
		}
	}

	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("order_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	orderUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}
	order, err := h.orderService.GetByID(orderUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) List(c *gin.Context) {
	tableID := c.Query("table_id")
	status := c.Query("status")

	var orders []models.Order
	var err error

	// Filtrar por mesa
	if tableID == "" {
		tableUUID, err := uuid.Parse(tableID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
			return
		}
		orders, err = h.orderService.GetByTable(tableUUID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch orders"})
			return
		}
	} else if status != "" {
		// Filtrar por status
		var orderStatus models.OrderStatus
		switch status {
		case string(models.OrderStatusPending):
			orderStatus = models.OrderStatusPending
		case string(models.OrderStatusPreparing):
			orderStatus = models.OrderStatusPreparing
		case string(models.OrderStatusReady):
			orderStatus = models.OrderStatusReady
		case string(models.OrderStatusDelivered):
			orderStatus = models.OrderStatusDelivered
		case string(models.OrderStatusPaid):
			orderStatus = models.OrderStatusPaid
		case string(models.OrderStatusCancelled):
			orderStatus = models.OrderStatusCancelled
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
			return
		}

		orders, err = h.orderService.GetByStatus(orderStatus)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch orders"})
			return
		}
	} else {
		// Listar todos
		orders, err = h.orderService.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch orders"})
			return
		}
	}

	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("order_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var status models.OrderStatus
	switch req.Status {
	case string(models.OrderStatusPending):
		status = models.OrderStatusPending
	case string(models.OrderStatusPreparing):
		status = models.OrderStatusPreparing
	case string(models.OrderStatusReady):
		status = models.OrderStatusReady
	case string(models.OrderStatusDelivered):
		status = models.OrderStatusDelivered
	case string(models.OrderStatusPaid):
		status = models.OrderStatusPaid
	case string(models.OrderStatusCancelled):
		status = models.OrderStatusCancelled
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	orderID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}
	order, err := h.orderService.GetByID(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	// Lógica especial para marcar como pago
	if status == models.OrderStatusPaid && order.Status != models.OrderStatusPaid {
		now := time.Now()
		order.PaidAt = &now

		// Criar transação financeira
		userID := c.GetString("user_id")
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
			return
		}
		if err := h.orderService.RegisterPayment(order, userUUID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register payment"})
			return
		}

		// Liberar mesa se o pedido estiver vinculado a uma
		if order.TableID == nil {
			if err := h.tableService.UpdateStatus(*order.TableID, models.TableStatusFree); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"message": "order status updated but failed to free table",
				})
				return
			}

			// Remover associação com o pedido atual
			var nilOrderID *uuid.UUID
			if err := h.tableService.SetCurrentOrder(*order.TableID, nilOrderID); err != nil {
				c.JSON(http.StatusOK, gin.H{
					"message": "order status updated but failed to unlink order from table",
				})
				return
			}
		}
	}

	if err := h.orderService.UpdateStatus(orderID, status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order status updated successfully"})
}

func (h *OrderHandler) AddItem(c *gin.Context) {
	id := c.Param("order_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	orderUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	var req OrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.orderService.GetByID(orderUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if order.Status == models.OrderStatusPaid || order.Status == models.OrderStatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot add items to a paid or cancelled order"})
		return
	}

	product, err := h.orderService.GetProductByID(req.ProductID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product not found"})
		return
	}

	if !product.InStock {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product out of stock"})
		return
	}

	item := &models.OrderItem{
		OrderID:   orderUUID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Price:     product.Price,
		Notes:     req.Notes,
	}

	if err := h.orderService.AddItem(item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *OrderHandler) RemoveItem(c *gin.Context) {
	orderID := c.Param("order_id")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	itemID := c.Param("item_id")
	if itemID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	itemUUID, err := uuid.Parse(itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	order, err := h.orderService.GetByID(orderUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	if order.Status == models.OrderStatusPaid || order.Status == models.OrderStatusCancelled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot remove items from a paid or cancelled order"})
		return
	}

	if err := h.orderService.RemoveItem(orderUUID, itemUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item removed successfully"})
}
