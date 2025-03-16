// internal/services/finance_service.go
package services

import (
	"time"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/domain/repositories"

	"github.com/google/uuid"
)

type FinanceService struct {
	financeRepo repositories.FinanceRepository
}

func NewFinanceService(financeRepo repositories.FinanceRepository) *FinanceService {
	return &FinanceService{
		financeRepo: financeRepo,
	}
}

func (s *FinanceService) Create(transaction *models.FinancialTransaction) error {
	return s.financeRepo.Create(transaction)
}

func (s *FinanceService) GetByID(id uuid.UUID) (*models.FinancialTransaction, error) {
	return s.financeRepo.FindByID(id)
}

func (s *FinanceService) Update(transaction *models.FinancialTransaction) error {
	return s.financeRepo.Update(transaction)
}

func (s *FinanceService) Delete(id uuid.UUID) error {
	return s.financeRepo.Delete(id)
}

func (s *FinanceService) List() ([]models.FinancialTransaction, error) {
	return s.financeRepo.List()
}

func (s *FinanceService) GetByType(transactionType models.TransactionType) ([]models.FinancialTransaction, error) {
	return s.financeRepo.FindByType(transactionType)
}

func (s *FinanceService) GetByDateRange(startDate, endDate time.Time) ([]models.FinancialTransaction, error) {
	return s.financeRepo.FindByDateRange(startDate, endDate)
}

func (s *FinanceService) GetByOrder(orderID uuid.UUID) ([]models.FinancialTransaction, error) {
	return s.financeRepo.FindByOrder(orderID)
}

func (s *FinanceService) GetDailySummary(date time.Time) (float64, float64, error) {
	return s.financeRepo.GetDailySummary(date)
}

func (s *FinanceService) GetMonthlySummary(year, month int) (float64, float64, error) {
	return s.financeRepo.GetMonthlySummary(year, month)
}
