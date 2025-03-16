package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive   SubscriptionStatus = "active"
	SubscriptionStatusInactive SubscriptionStatus = "inactive"
	SubscriptionStatusTrial    SubscriptionStatus = "trial"
)

// Restaurant representa um estabelecimento no sistema SaaS
type Restaurant struct {
	ID               uuid.UUID          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name             string             `gorm:"size:100;not null" json:"name"`
	Description      string             `gorm:"size:255" json:"description"`
	Address          string             `gorm:"size:255" json:"address"`
	Phone            string             `gorm:"size:20" json:"phone"`
	Email            string             `gorm:"size:100" json:"email"`
	Logo             string             `gorm:"size:255" json:"logo"`
	SubscriptionPlan string             `gorm:"size:50" json:"subscription_plan"`
	Status           SubscriptionStatus `gorm:"size:20;not null;default:'trial'" json:"status"`
	TrialEndsAt      *time.Time         `json:"trial_ends_at"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
	Users            []User             `json:"users,omitempty" gorm:"foreignKey:RestaurantID"`
	Tables           []Table            `json:"tables,omitempty" gorm:"foreignKey:RestaurantID"`
	Products         []Product          `json:"products,omitempty" gorm:"foreignKey:RestaurantID"`
	Orders           []Order            `json:"orders,omitempty" gorm:"foreignKey:RestaurantID"`
}

func (r *Restaurant) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
