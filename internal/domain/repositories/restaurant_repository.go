package repositories

import (
	"api-jet-manager/internal/domain/models"

	"github.com/google/uuid"
)

type RestaurantRepository interface {
	Create(restaurant *models.Restaurant) error
	FindByID(id uuid.UUID) (*models.Restaurant, error)
	Update(restaurant *models.Restaurant) error
	Delete(id uuid.UUID) error
	List() ([]models.Restaurant, error)
	FindByStatus(status models.SubscriptionStatus) ([]models.Restaurant, error)
	FindByName(name string) ([]models.Restaurant, error)
	UpdateStatus(id uuid.UUID, status models.SubscriptionStatus) error
}
