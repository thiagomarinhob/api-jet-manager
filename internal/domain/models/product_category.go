package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ProductCategory representa uma categoria de produto personalizada
type ProductCategory struct {
	ID           uuid.UUID   `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	RestaurantID uuid.UUID   `gorm:"type:uuid;not null" json:"restaurant_id"`
	Restaurant   *Restaurant `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID"`
	Name         string      `gorm:"size:100;not null" json:"name"`
	Description  string      `gorm:"size:255" json:"description"`
	Active       bool        `gorm:"default:true" json:"active"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
}

func (pc *ProductCategory) BeforeCreate(tx *gorm.DB) error {
	if pc.ID == uuid.Nil {
		pc.ID = uuid.New()
	}
	return nil
}
