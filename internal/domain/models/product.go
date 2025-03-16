package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductCategory string

const (
	ProductCategoryFood    ProductCategory = "food"
	ProductCategoryDrink   ProductCategory = "drink"
	ProductCategoryDessert ProductCategory = "dessert"
)

type Product struct {
	ID           uuid.UUID       `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	RestaurantID uuid.UUID       `json:"restaurant_id" gorm:"type:uuid;not null"`
	Restaurant   *Restaurant     `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID"`
	Name         string          `gorm:"size:100;not null" json:"name"`
	Description  string          `gorm:"size:255" json:"description"`
	Price        float64         `gorm:"not null" json:"price"`
	Category     ProductCategory `gorm:"size:20;not null" json:"category"`
	InStock      bool            `gorm:"default:true" json:"in_stock"`
	ImageURL     string          `gorm:"size:255" json:"image_url"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

func (p *Product) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
