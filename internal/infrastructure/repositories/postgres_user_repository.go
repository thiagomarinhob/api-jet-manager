// internal/infrastructure/repositories/postgres_user_repository.go
package repositories

import (
	"errors"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/infrastructure/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	DB *gorm.DB
}

func NewPostgresUserRepository(db *database.PostgresDB) *PostgresUserRepository {
	return &PostgresUserRepository{
		DB: db.DB,
	}
}

func (r *PostgresUserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *PostgresUserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("id = ?", id).First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *PostgresUserRepository) Update(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *PostgresUserRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.User{}, id).Error
}

func (r *PostgresUserRepository) List() ([]models.User, error) {
	var users []models.User
	if err := r.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *PostgresUserRepository) FindByType(userType models.UserType) ([]models.User, error) {
	var users []models.User
	if err := r.DB.Where("type = ?", userType).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *PostgresUserRepository) FindByRestaurant(restaurantID uuid.UUID) ([]models.User, error) {
	var users []models.User
	if err := r.DB.Where("restaurant_id = ?", restaurantID).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}
