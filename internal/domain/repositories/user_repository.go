package repositories

import (
	"api-jet-manager/internal/domain/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uuid.UUID) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uuid.UUID) error
	List() ([]models.User, error)
	FindByType(userType models.UserType) ([]models.User, error)
	FindByRestaurant(restaurantID uuid.UUID) ([]models.User, error)
}
