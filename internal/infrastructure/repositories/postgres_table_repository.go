// internal/infrastructure/repositories/postgres_table_repository.go
package repositories

import (
	"errors"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/infrastructure/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresTableRepository struct {
	DB *gorm.DB
}

func NewPostgresTableRepository(db *database.PostgresDB) *PostgresTableRepository {
	return &PostgresTableRepository{
		DB: db.DB,
	}
}

func (r *PostgresTableRepository) Create(table *models.Table) error {
	// O restaurant_id já deve estar definido no objeto table antes de chamar este método
	return r.DB.Create(table).Error
}

func (r *PostgresTableRepository) FindByID(restauranteID, id uuid.UUID) (*models.Table, error) {
	var table models.Table
	if err := r.DB.Where("restaurant_id = ? AND id = ?", restauranteID, id).First(&table).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("table not found")
		}
		return nil, err
	}
	return &table, nil
}

func (r *PostgresTableRepository) FindByNumber(restauranteID uuid.UUID, number int) (*models.Table, error) {
	var table models.Table
	if err := r.DB.Where("restaurant_id = ? AND number = ?", restauranteID, number).First(&table).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("table not found")
		}
		return nil, err
	}
	return &table, nil
}

func (r *PostgresTableRepository) Update(table *models.Table) error {
	// Assumindo que o restaurant_id já está definido no objeto table
	// Opcionalmente, você pode adicionar uma verificação adicional de segurança:
	// return r.DB.Where("restaurant_id = ?", table.RestauranteID).Save(table).Error
	return r.DB.Save(table).Error
}

func (r *PostgresTableRepository) Delete(restauranteID, id uuid.UUID) error {
	return r.DB.Where("restaurant_id = ?", restauranteID).Delete(&models.Table{}, id).Error
}

func (r *PostgresTableRepository) List(restauranteID uuid.UUID) ([]models.Table, error) {
	var tables []models.Table
	if err := r.DB.Where("restaurant_id = ?", restauranteID).Find(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

func (r *PostgresTableRepository) UpdateStatus(restauranteID, id uuid.UUID, status models.TableStatus) error {
	return r.DB.Model(&models.Table{}).Where("restaurant_id = ? AND id = ?", restauranteID, id).Update("status", status).Error
}

func (r *PostgresTableRepository) SetCurrentOrder(restauranteID, id uuid.UUID, orderID *uuid.UUID) error {
	return r.DB.Model(&models.Table{}).Where("restaurant_id = ? AND id = ?", restauranteID, id).Update("current_order", orderID).Error
}
