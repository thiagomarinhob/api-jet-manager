package handlers

import (
	"net/http"
	"strconv"
	"time"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FinanceTransactionRequest struct {
	Type        string     `json:"type" binding:"required"`
	Category    string     `json:"category" binding:"required"`
	Amount      float64    `json:"amount" binding:"required,gt=0"`
	Description string     `json:"description" binding:"required"`
	Date        string     `json:"date" binding:"required"`
	OrderID     *uuid.UUID `json:"order_id"`
}

type FinanceHandler struct {
	financeService *services.FinanceService
}

func NewFinanceHandler(financeService *services.FinanceService) *FinanceHandler {
	return &FinanceHandler{
		financeService: financeService,
	}
}

func (h *FinanceHandler) Create(c *gin.Context) {
	var req FinanceTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Validação do tipo de transação
	var transactionType models.TransactionType
	switch req.Type {
	case string(models.TransactionTypeIncome):
		transactionType = models.TransactionTypeIncome
	case string(models.TransactionTypeExpense):
		transactionType = models.TransactionTypeExpense
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction type"})
		return
	}

	// Validação da categoria
	var category models.TransactionCategory
	switch req.Category {
	case string(models.TransactionCategorySales):
		category = models.TransactionCategorySales
	case string(models.TransactionCategoryOther):
		category = models.TransactionCategoryOther
	case string(models.TransactionCategoryIngredients):
		category = models.TransactionCategoryIngredients
	case string(models.TransactionCategoryUtilities):
		category = models.TransactionCategoryUtilities
	case string(models.TransactionCategorySalaries):
		category = models.TransactionCategorySalaries
	case string(models.TransactionCategoryRent):
		category = models.TransactionCategoryRent
	case string(models.TransactionCategoryEquipment):
		category = models.TransactionCategoryEquipment
	case string(models.TransactionCategoryMaintenance):
		category = models.TransactionCategoryMaintenance
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction category"})
		return
	}

	// Conversão da data
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format (required: YYYY-MM-DD)"})
		return
	}

	transaction := &models.FinancialTransaction{
		Type:        transactionType,
		Category:    category,
		Amount:      req.Amount,
		Description: req.Description,
		OrderID:     req.OrderID,
		UserID:      userID.(uuid.UUID),
		Date:        date,
	}

	if err := h.financeService.Create(transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

func (h *FinanceHandler) GetByID(c *gin.Context) {
	id := c.Param("transaction_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	transactionID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID format"})
		return
	}

	restaurant_id := c.Param("restaurant_id")
	if restaurant_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	restaurant_uuid, errRes := uuid.Parse(restaurant_id)
	if errRes != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	transaction, err := h.financeService.GetByID(restaurant_uuid, transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *FinanceHandler) List(c *gin.Context) {
	transactionType := c.Query("type")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	restaurant_id := c.Param("restaurant_id")
	if restaurant_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	restaurant_uuid, errRes := uuid.Parse(restaurant_id)
	if errRes != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	var transactions []models.FinancialTransaction
	var err error

	// Filtrar por tipo
	if transactionType != "" {
		var tType models.TransactionType
		switch transactionType {
		case string(models.TransactionTypeIncome):
			tType = models.TransactionTypeIncome
		case string(models.TransactionTypeExpense):
			tType = models.TransactionTypeExpense
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction type"})
			return
		}

		transactions, err = h.financeService.GetByType(restaurant_uuid, tType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transactions"})
			return
		}
	} else if startDateStr != "" && endDateStr != "" {
		// Filtrar por intervalo de datas
		startDate, err := time.Parse("2006-01-02", startDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start date format (required: YYYY-MM-DD)"})
			return
		}

		endDate, err := time.Parse("2006-01-02", endDateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end date format (required: YYYY-MM-DD)"})
			return
		}

		transactions, err = h.financeService.GetByDateRange(restaurant_uuid, startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transactions"})
			return
		}
	} else {
		// Listar todas
		transactions, err = h.financeService.List(restaurant_uuid)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch transactions"})
			return
		}
	}

	c.JSON(http.StatusOK, transactions)
}

func (h *FinanceHandler) GetSummary(c *gin.Context) {
	period := c.Query("period")
	dateStr := c.Query("date")
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	restaurant_id := c.Param("restaurant_id")
	if restaurant_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	restaurant_uuid, errRes := uuid.Parse(restaurant_id)
	if errRes != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	switch period {
	case "daily":
		if dateStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "date parameter is required for daily summary"})
			return
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format (required: YYYY-MM-DD)"})
			return
		}

		income, expense, err := h.financeService.GetDailySummary(restaurant_uuid, date)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate summary"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"date":    dateStr,
			"income":  income,
			"expense": expense,
			"balance": income - expense,
		})

	case "monthly":
		if yearStr == "" || monthStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "year and month parameters are required for monthly summary"})
			return
		}

		year, err := strconv.Atoi(yearStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year format"})
			return
		}

		month, err := strconv.Atoi(monthStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid month format"})
			return
		}

		if month < 1 || month > 12 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "month must be between 1 and 12"})
			return
		}

		income, expense, err := h.financeService.GetMonthlySummary(restaurant_uuid, year, month)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to calculate summary"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"year":    year,
			"month":   month,
			"income":  income,
			"expense": expense,
			"balance": income - expense,
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period parameter (supported: daily, monthly)"})
	}
}

func (h *FinanceHandler) Update(c *gin.Context) {
	id := c.Param("transaction_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	restaurant_id := c.Param("restaurant_id")
	if restaurant_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	restaurant_uuid, errRes := uuid.Parse(restaurant_id)
	if errRes != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	var req FinanceTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactionID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID format"})
		return
	}

	transaction, err := h.financeService.GetByID(restaurant_uuid, transactionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transaction not found"})
		return
	}

	// Validação do tipo de transação
	var transactionType models.TransactionType
	switch req.Type {
	case string(models.TransactionTypeIncome):
		transactionType = models.TransactionTypeIncome
	case string(models.TransactionTypeExpense):
		transactionType = models.TransactionTypeExpense
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction type"})
		return
	}

	// Validação da categoria
	var category models.TransactionCategory
	switch req.Category {
	case string(models.TransactionCategorySales):
		category = models.TransactionCategorySales
	case string(models.TransactionCategoryOther):
		category = models.TransactionCategoryOther
	case string(models.TransactionCategoryIngredients):
		category = models.TransactionCategoryIngredients
	case string(models.TransactionCategoryUtilities):
		category = models.TransactionCategoryUtilities
	case string(models.TransactionCategorySalaries):
		category = models.TransactionCategorySalaries
	case string(models.TransactionCategoryRent):
		category = models.TransactionCategoryRent
	case string(models.TransactionCategoryEquipment):
		category = models.TransactionCategoryEquipment
	case string(models.TransactionCategoryMaintenance):
		category = models.TransactionCategoryMaintenance
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction category"})
		return
	}

	// Conversão da data
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format (required: YYYY-MM-DD)"})
		return
	}

	transaction.Type = transactionType
	transaction.Category = category
	transaction.Amount = req.Amount
	transaction.Description = req.Description
	transaction.Date = date
	transaction.OrderID = req.OrderID

	if err := h.financeService.Update(transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h *FinanceHandler) Delete(c *gin.Context) {
	id := c.Param("transaction_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID"})
		return
	}

	transactionID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid transaction ID format"})
		return
	}

	restaurant_id := c.Param("restaurant_id")
	if restaurant_id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	restaurant_uuid, errRes := uuid.Parse(restaurant_id)
	if errRes != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	if err := h.financeService.Delete(restaurant_uuid, transactionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "transaction deleted successfully"})
}
