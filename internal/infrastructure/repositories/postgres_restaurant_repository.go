package repositories

import (
	"errors"
	"fmt"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/infrastructure/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresRestaurantRepository struct {
	DB *gorm.DB
}

func NewPostgresRestaurantRepository(db *database.PostgresDB) *PostgresRestaurantRepository {
	return &PostgresRestaurantRepository{
		DB: db.DB,
	}
}

func (r *PostgresRestaurantRepository) Create(restaurant *models.Restaurant) error {
	return r.DB.Create(restaurant).Error
}

func (r *PostgresRestaurantRepository) FindByID(id uuid.UUID) (*models.Restaurant, error) {
	var restaurant models.Restaurant
	fmt.Print("id repository", id)
	if err := r.DB.Where("id = ?", id).First(&restaurant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("restaurant not found")
		}
		return nil, err
	}
	return &restaurant, nil
}

func (r *PostgresRestaurantRepository) Update(restaurant *models.Restaurant) error {
	return r.DB.Save(restaurant).Error
}

func (r *PostgresRestaurantRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.Restaurant{}, id).Error
}

func (r *PostgresRestaurantRepository) List() ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	if err := r.DB.Find(&restaurants).Error; err != nil {
		return nil, err
	}
	return restaurants, nil
}

func (r *PostgresRestaurantRepository) FindByStatus(status models.SubscriptionStatus) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	if err := r.DB.Where("status = ?", status).Find(&restaurants).Error; err != nil {
		return nil, err
	}
	return restaurants, nil
}

func (r *PostgresRestaurantRepository) FindByName(name string) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	if err := r.DB.Where("name LIKE ?", "%"+name+"%").Find(&restaurants).Error; err != nil {
		return nil, err
	}
	return restaurants, nil
}

func (r *PostgresRestaurantRepository) UpdateStatus(id uuid.UUID, status models.SubscriptionStatus) error {
	return r.DB.Model(&models.Restaurant{}).Where("id = ?", id).Update("status", status).Error
}
