package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TableStatus string

const (
	TableStatusFree     TableStatus = "free"
	TableStatusOccupied TableStatus = "occupied"
	TableStatusReserved TableStatus = "reserved"
)

type Table struct {
	ID             uuid.UUID   `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	RestaurantID   uuid.UUID   `json:"restaurant_id" gorm:"type:uuid;not null"`
	Restaurant     *Restaurant `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID"`
	Number         int         `gorm:"not null" json:"number"`
	Capacity       int         `gorm:"not null" json:"capacity"`
	Status         TableStatus `gorm:"size:20;not null;default:'free'" json:"status"`
	CurrentOrderID *uuid.UUID  `json:"current_order_id" gorm:"type:uuid"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

// Chave composta para garantir que números de mesa são únicos por restaurante
func (Table) TableName() string {
	return "tables"
}

func (t *Table) BeforeCreate(tx *gorm.DB) error {
	// Verifica se já existe uma mesa com o mesmo número no mesmo restaurante
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	var count int64
	tx.Model(&Table{}).Where("restaurant_id = ? AND number = ?", t.RestaurantID, t.Number).Count(&count)
	if count > 0 {
		return fmt.Errorf("já existe uma mesa com o número %d neste restaurante", t.Number)
	}
	return nil
}
