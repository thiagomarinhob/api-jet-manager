package repositories

import (
	"time"

	"api-jet-manager/internal/domain/models"

	"github.com/google/uuid"
)

type FinanceRepository interface {
	Create(transaction *models.FinancialTransaction) error
	FindByID(restaurantID, id uuid.UUID) (*models.FinancialTransaction, error)
	Update(transaction *models.FinancialTransaction) error
	Delete(restaurantID, id uuid.UUID) error
	List(restaurantID uuid.UUID) ([]models.FinancialTransaction, error)
	FindByType(restaurantID uuid.UUID, transactionType models.TransactionType) ([]models.FinancialTransaction, error)
	FindByDateRange(restaurantID uuid.UUID, startDate, endDate time.Time) ([]models.FinancialTransaction, error)
	FindByOrder(restaurantID, orderID uuid.UUID) ([]models.FinancialTransaction, error)
	GetDailySummary(restaurantID uuid.UUID, date time.Time) (float64, float64, error) // Retorna (receitas, despesas)
	GetMonthlySummary(restaurantID uuid.UUID, year int, month int) (float64, float64, error)
}
