package repositories

import (
	"api-jet-manager/internal/domain/models"

	"github.com/google/uuid"
)

type UserRepository interface {
	Create(user *models.User) error
	FindByID(restaurantID, id uuid.UUID) (*models.User, error)
	FindByEmail(restaurantID uuid.UUID, email string) (*models.User, error)
	// Método adicional para busca global por email (útil para login)
	FindByEmailGlobal(email string) (*models.User, error)
	Update(user *models.User) error
	Delete(restaurantID, id uuid.UUID) error
	List(restaurantID uuid.UUID) ([]models.User, error)
	FindByType(restaurantID uuid.UUID, userType models.UserType) ([]models.User, error)
	FindByTypeGlobal(userType models.UserType) ([]models.User, error)
	FindByRestaurant(restaurantID uuid.UUID) ([]models.User, error)
}
