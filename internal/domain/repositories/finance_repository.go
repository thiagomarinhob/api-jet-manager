package repositories

import (
	"time"

	"api-jet-manager/internal/domain/models"

	"github.com/google/uuid"
)

type FinanceRepository interface {
	Create(transaction *models.FinancialTransaction) error
	FindByID(id uuid.UUID) (*models.FinancialTransaction, error)
	Update(transaction *models.FinancialTransaction) error
	Delete(id uuid.UUID) error
	List() ([]models.FinancialTransaction, error)
	FindByType(transactionType models.TransactionType) ([]models.FinancialTransaction, error)
	FindByDateRange(startDate, endDate time.Time) ([]models.FinancialTransaction, error)
	FindByOrder(orderID uuid.UUID) ([]models.FinancialTransaction, error)
	GetDailySummary(date time.Time) (float64, float64, error) // Retorna (receitas, despesas)
	GetMonthlySummary(year int, month int) (float64, float64, error)
}
