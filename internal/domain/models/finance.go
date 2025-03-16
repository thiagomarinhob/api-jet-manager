package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeIncome  TransactionType = "income"
	TransactionTypeExpense TransactionType = "expense"
)

type TransactionCategory string

const (
	// Receitas
	TransactionCategorySales TransactionCategory = "sales"
	TransactionCategoryOther TransactionCategory = "other_income"

	// Despesas
	TransactionCategoryIngredients TransactionCategory = "ingredients"
	TransactionCategoryUtilities   TransactionCategory = "utilities"
	TransactionCategorySalaries    TransactionCategory = "salaries"
	TransactionCategoryRent        TransactionCategory = "rent"
	TransactionCategoryEquipment   TransactionCategory = "equipment"
	TransactionCategoryMaintenance TransactionCategory = "maintenance"
)

type FinancialTransaction struct {
	ID            uuid.UUID           `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	RestaurantID  uuid.UUID           `json:"restaurant_id" gorm:"type:uuid;not null"`
	Restaurant    *Restaurant         `json:"restaurant,omitempty" gorm:"foreignKey:RestaurantID"`
	Type          TransactionType     `gorm:"size:20;not null" json:"type"`
	Category      TransactionCategory `gorm:"size:30;not null" json:"category"`
	Amount        float64             `gorm:"not null" json:"amount"`
	Description   string              `gorm:"size:255" json:"description"`
	OrderID       *uuid.UUID          `json:"order_id,omitempty" gorm:"type:uuid"`
	Order         *Order              `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	UserID        uuid.UUID           `json:"user_id" gorm:"type:uuid;not null"`
	User          *User               `json:"user,omitempty" gorm:"foreignKey:UserID"`
	PaymentMethod string              `gorm:"size:30" json:"payment_method"`
	Date          time.Time           `gorm:"not null" json:"date"`
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
}

func (ft *FinancialTransaction) BeforeCreate(tx *gorm.DB) error {
	if ft.ID == uuid.Nil {
		ft.ID = uuid.New()
	}
	return nil
}
