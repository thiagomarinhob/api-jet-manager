// internal/infrastructure/repositories/postgres_finance_repository.go
package repositories

import (
	"errors"
	"time"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/infrastructure/database"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresFinanceRepository struct {
	DB *gorm.DB
}

func NewPostgresFinanceRepository(db *database.PostgresDB) *PostgresFinanceRepository {
	return &PostgresFinanceRepository{
		DB: db.DB,
	}
}

func (r *PostgresFinanceRepository) FindByDateRange(startDate, endDate time.Time) ([]models.FinancialTransaction, error) {
	var transactions []models.FinancialTransaction
	if err := r.DB.Where("date BETWEEN ? AND ?", startDate, endDate).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *PostgresFinanceRepository) FindByOrder(orderID uuid.UUID) ([]models.FinancialTransaction, error) {
	var transactions []models.FinancialTransaction
	if err := r.DB.Where("order_id = ?", orderID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *PostgresFinanceRepository) GetDailySummary(date time.Time) (float64, float64, error) {
	// Define o início e o fim do dia
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	// Calcula o total de receitas
	var income float64
	if err := r.DB.Model(&models.FinancialTransaction{}).
		Where("type = ? AND date BETWEEN ? AND ?", models.TransactionTypeIncome, startOfDay, endOfDay).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&income).Error; err != nil {
		return 0, 0, err
	}

	// Calcula o total de despesas
	var expense float64
	if err := r.DB.Model(&models.FinancialTransaction{}).
		Where("type = ? AND date BETWEEN ? AND ?", models.TransactionTypeExpense, startOfDay, endOfDay).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&expense).Error; err != nil {
		return 0, 0, err
	}

	return income, expense, nil
}

func (r *PostgresFinanceRepository) GetMonthlySummary(year int, month int) (float64, float64, error) {
	// Define o início e o fim do mês
	loc := time.Now().Location()
	startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)

	var endOfMonth time.Time
	if month == 12 {
		endOfMonth = time.Date(year+1, 1, 1, 0, 0, 0, 0, loc)
	} else {
		endOfMonth = time.Date(year, time.Month(month+1), 1, 0, 0, 0, 0, loc)
	}

	// Calcula o total de receitas
	var income float64
	if err := r.DB.Model(&models.FinancialTransaction{}).
		Where("type = ? AND date BETWEEN ? AND ?", models.TransactionTypeIncome, startOfMonth, endOfMonth).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&income).Error; err != nil {
		return 0, 0, err
	}

	// Calcula o total de despesas
	var expense float64
	if err := r.DB.Model(&models.FinancialTransaction{}).
		Where("type = ? AND date BETWEEN ? AND ?", models.TransactionTypeExpense, startOfMonth, endOfMonth).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&expense).Error; err != nil {
		return 0, 0, err
	}

	return income, expense, nil
}

func (r *PostgresFinanceRepository) Create(transaction *models.FinancialTransaction) error {
	return r.DB.Create(transaction).Error
}

func (r *PostgresFinanceRepository) FindByID(id uuid.UUID) (*models.FinancialTransaction, error) {
	var transaction models.FinancialTransaction
	if err := r.DB.Where("id = ?", id).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}
	return &transaction, nil
}

func (r *PostgresFinanceRepository) Update(transaction *models.FinancialTransaction) error {
	return r.DB.Save(transaction).Error
}

func (r *PostgresFinanceRepository) Delete(id uuid.UUID) error {
	return r.DB.Delete(&models.FinancialTransaction{}, id).Error
}

func (r *PostgresFinanceRepository) List() ([]models.FinancialTransaction, error) {
	var transactions []models.FinancialTransaction
	if err := r.DB.Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *PostgresFinanceRepository) FindByType(transactionType models.TransactionType) ([]models.FinancialTransaction, error) {
	var transactions []models.FinancialTransaction
	if err := r.DB.Where("type = ?", transactionType).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
