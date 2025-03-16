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
	return r.DB.Create(table).Error
}

func (r *PostgresTableRepository) FindByID(id uuid.UUID) (*models.Table, error) {
	var table models.Table
	if err := r.DB.Where("id = ?", id).First(&table, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("table not found")
		}
		return nil, err
	}
	return &table, nil
}

func (r *PostgresTableRepository) FindByNumber(number int) (*models.Table, error) {
	var table models.Table
	if err := r.DB.Where("number = ?", number).First(&table).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("table not found")
		}
		return nil, err
	}
	return &table, nil
}

func (r *PostgresTableRepository) Update(table *models.Table) error {
	return r.DB.Save(table).Error
}

func (r *PostgresTableRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.Table{}, id).Error
}

func (r *PostgresTableRepository) List() ([]models.Table, error) {
	var tables []models.Table
	if err := r.DB.Find(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

func (r *PostgresTableRepository) UpdateStatus(id uuid.UUID, status models.TableStatus) error {
	return r.DB.Model(&models.Table{}).Where("id = ?", id).Update("status", status).Error
}

func (r *PostgresTableRepository) SetCurrentOrder(id uuid.UUID, orderID *uuid.UUID) error {
	return r.DB.Model(&models.Table{}).Where("id = ?", id).Update("current_order", orderID).Error
}
