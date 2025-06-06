package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SelectionType define como as opções podem ser selecionadas
type SelectionType string

const (
	SingleSelection    SelectionType = "single"               // Apenas uma opção
	MultipleNoRepeat   SelectionType = "multiple_no_repeat"   // Múltiplas sem repetição
	MultipleWithRepeat SelectionType = "multiple_with_repeat" // Múltiplas com repetição
)

// PriceMethod define como o preço é calculado para seleções múltiplas
type PriceMethod string

const (
	Sum     PriceMethod = "sum"     // Soma das opções selecionadas
	Average PriceMethod = "average" // Média das opções selecionadas
	Highest PriceMethod = "highest" // Opção mais cara
	Lowest  PriceMethod = "lowest"  // Opção mais barata
)

// Option representa uma opção individual dentro de um complemento
type Option struct {
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	AddonID     uuid.UUID `gorm:"type:uuid;not null" json:"addon_id"`
	Addon       *Addon    `json:"addon,omitempty" gorm:"foreignKey:AddonID"`
	Name        string    `gorm:"size:100;not null" json:"name"`
	Price       float64   `gorm:"not null" json:"price"`
	Active      bool      `gorm:"default:true" json:"active"`
	MaxQuantity int       `gorm:"default:1" json:"max_quantity"` // Quantidade máxima permitida por opção (para MultipleWithRepeat)
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (o *Option) BeforeCreate(tx *gorm.DB) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return nil
}

// Addon representa um grupo de opções que podem ser adicionadas a um produto
type Addon struct {
	ID            uuid.UUID     `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Title         string        `gorm:"size:100;not null" json:"title"`
	ProductID     uuid.UUID     `gorm:"type:uuid;not null" json:"product_id"`
	Product       *Product      `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Options       []Option      `json:"options,omitempty" gorm:"foreignKey:AddonID"`
	SelectionType SelectionType `gorm:"size:30;not null" json:"selection_type"`
	MinSelections int           `gorm:"default:0" json:"min_selections"` // Mínimo de seleções
	MaxSelections int           `gorm:"default:1" json:"max_selections"` // Máximo de seleções
	Required      bool          `gorm:"default:false" json:"required"`   // Se o cliente é obrigado a escolher
	PriceMethod   PriceMethod   `gorm:"size:20;not null" json:"price_method"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

func (a *Addon) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
