package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductType é mantido para compatibilidade com o código existente
type ProductType string

const (
	ProductTypeFood    ProductType = "food"
	ProductTypeDrink   ProductType = "drink"
	ProductTypeDessert ProductType = "dessert"
)

// Product representa um produto no menu do restaurante
type Product struct {
	ID           uuid.UUID        `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	RestaurantID uuid.UUID        `gorm:"type:uuid;not null" json:"restaurant_id"`
	Restaurant   *Restaurant      `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID"`
	Name         string           `gorm:"size:100;not null" json:"name"`
	Description  string           `gorm:"size:255" json:"description"`
	Price        float64          `gorm:"not null" json:"price"`
	CategoryID   uuid.UUID        `gorm:"type:uuid;not null" json:"category_id"`
	Category     *ProductCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Type         ProductType      `gorm:"size:20" json:"type"` // Campo mantido para compatibilidade
	InStock      bool             `gorm:"default:true" json:"in_stock"`
	ImageURL     string           `gorm:"size:255" json:"image_url"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
