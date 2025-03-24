package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPreparing OrderStatus = "preparing"
	OrderStatusReady     OrderStatus = "ready"
	OrderStatusDelivered OrderStatus = "delivered"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusCancelled OrderStatus = "cancelled"
)

type OrderType string

const (
	OrderTypeInHouse  OrderType = "in_house" // Pedido para consumo no local
	OrderTypeDelivery OrderType = "delivery" // Pedido para entrega
	OrderTypeTakeaway OrderType = "takeaway" // Pedido para retirada
)

type Order struct {
	ID              uuid.UUID   `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	RestaurantID    uuid.UUID   `json:"restaurant_id" gorm:"type:uuid;not null"`
	Restaurant      *Restaurant `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID"`
	TableID         *uuid.UUID  `json:"table_id" gorm:"type:uuid"`
	Table           *Table      `json:"table,omitempty" gorm:"foreignKey:TableID"`
	UserID          uuid.UUID   `json:"user_id" gorm:"type:uuid;not null"`
	User            *User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Code            string      `gorm:"size:20" json:"code"`
	CustomerName    string      `gorm:"size:100" json:"customer_name"`
	CustomerPhone   string      `gorm:"size:20" json:"customer_phone"`
	CustomerEmail   string      `gorm:"size:100" json:"customer_email"`
	Type            OrderType   `gorm:"size:20;not null;default:'in_house'" json:"type"`
	Status          OrderStatus `gorm:"size:20;not null;default:'pending'" json:"status"`
	OrderItems      []OrderItem `json:"order_items,omitempty" gorm:"foreignKey:OrderID"`
	TotalAmount     float64     `gorm:"not null;default:0" json:"total_amount"`
	Notes           string      `gorm:"size:255" json:"notes"`
	DeliveryAddress string      `gorm:"size:255" json:"delivery_address"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
	PaidAt          *time.Time  `json:"paid_at"`
	DeliveredAt     *time.Time  `json:"delivered_at"`
}

type OrderItem struct {
	ID        uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	OrderID   uuid.UUID `json:"order_id" gorm:"type:uuid;not null"`
	Order     *Order    `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:uuid;not null"`
	Product   *Product  `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Quantity  int       `gorm:"not null;default:1" json:"quantity"`
	Price     float64   `gorm:"not null" json:"price"` // Preço no momento da venda
	Notes     string    `gorm:"size:255" json:"notes"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (oi *OrderItem) BeforeCreate(tx *gorm.DB) error {
	if oi.ID == uuid.Nil {
		oi.ID = uuid.New()
	}
	return nil
}

func (o *Order) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}
