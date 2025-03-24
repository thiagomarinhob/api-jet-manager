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

func (s *FinanceService) GetByID(restaurant_id, id uuid.UUID) (*models.FinancialTransaction, error) {
	return s.financeRepo.FindByID(restaurant_id, id)
}

func (s *FinanceService) Update(transaction *models.FinancialTransaction) error {
	return s.financeRepo.Update(transaction)
}

func (s *FinanceService) Delete(restaurant_id, id uuid.UUID) error {
	return s.financeRepo.Delete(restaurant_id, id)
}

func (s *FinanceService) List(restaurant_id uuid.UUID) ([]models.FinancialTransaction, error) {
	return s.financeRepo.List(restaurant_id)
}

func (s *FinanceService) GetByType(restaurant_id uuid.UUID, transactionType models.TransactionType) ([]models.FinancialTransaction, error) {
	return s.financeRepo.FindByType(restaurant_id, transactionType)
}

func (s *FinanceService) GetByDateRange(restaurant_id uuid.UUID, startDate, endDate time.Time) ([]models.FinancialTransaction, error) {
	return s.financeRepo.FindByDateRange(restaurant_id, startDate, endDate)
}

func (s *FinanceService) GetByOrder(restaurant_id uuid.UUID, orderID uuid.UUID) ([]models.FinancialTransaction, error) {
	return s.financeRepo.FindByOrder(restaurant_id, orderID)
}

func (s *FinanceService) GetDailySummary(restaurant_id uuid.UUID, date time.Time) (float64, float64, error) {
	return s.financeRepo.GetDailySummary(restaurant_id, date)
}

func (s *FinanceService) GetMonthlySummary(restaurant_id uuid.UUID, year, month int) (float64, float64, error) {
	return s.financeRepo.GetMonthlySummary(restaurant_id, year, month)
}
