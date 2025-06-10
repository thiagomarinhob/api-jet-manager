package repositories

import (
	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/infrastructure/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresAddonRepository struct {
	DB *gorm.DB
}

func NewPostgresAddonRepository(db *database.PostgresDB) *PostgresAddonRepository {
	return &PostgresAddonRepository{
		DB: db.DB,
	}
}

func (r *PostgresAddonRepository) CreateAddon(addon *models.Addon) error {
	if err := r.DB.Create(addon).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostgresAddonRepository) GetAddonByID(restaurantID, id uuid.UUID) (*models.Addon, error) {
	var addon models.Addon
	if err := r.DB.Where("restaurant_id = ? AND id = ?", restaurantID, id).First(&addon).Preload("Options").Error; err != nil {
		return nil, err
	}
	return &addon, nil
}
