package repositories

import (
	"api-jet-manager/internal/domain/models"

	"github.com/google/uuid"
)

type AddonRepository interface {
	// Addon operations
	CreateAddon(addon *models.Addon) error
	GetAddonByID(restaurant_id, id uuid.UUID) (*models.Addon, error)
	GetAddonsByProductID(restaurant_id, productID uuid.UUID) ([]models.Addon, error)
	UpdateAddon(restaurant_id uuid.UUID, addon *models.Addon) error
	DeleteAddon(restaurant_id, id uuid.UUID) error

	// Option operations
	CreateOption(option *models.Option) error
	GetOptionByID(id uuid.UUID) (*models.Option, error)
	GetOptionsByAddonID(addonID uuid.UUID) ([]models.Option, error)
	UpdateOption(option *models.Option) error
	DeleteOption(id uuid.UUID) error
}
