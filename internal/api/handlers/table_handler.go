package handlers

import (
	"net/http"

	"api-jet-manager/internal/domain/models"
	"api-jet-manager/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TableRequest struct {
	Number   int `json:"number" binding:"required,min=1"`
	Capacity int `json:"capacity" binding:"required,min=1"`
}

type TableHandler struct {
	tableService *services.TableService
}

func NewTableHandler(tableService *services.TableService) *TableHandler {
	return &TableHandler{
		tableService: tableService,
	}
}

func (h *TableHandler) Create(c *gin.Context) {
	var req TableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	table := &models.Table{
		Number:   req.Number,
		Capacity: req.Capacity,
		Status:   models.TableStatusFree,
	}

	if err := h.tableService.Create(table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, table)
}

func (h *TableHandler) GetByID(c *gin.Context) {
	id := c.Param("table_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	table, err := h.tableService.GetByID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
		return
	}

	c.JSON(http.StatusOK, table)
}

func (h *TableHandler) List(c *gin.Context) {
	tables, err := h.tableService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tables"})
		return
	}

	c.JSON(http.StatusOK, tables)
}

func (h *TableHandler) Update(c *gin.Context) {
	id := c.Param("table_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	var req TableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	table, err := h.tableService.GetByID(uuid)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "table not found"})
		return
	}

	table.Number = req.Number
	table.Capacity = req.Capacity

	if err := h.tableService.Update(table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, table)
}

func (h *TableHandler) Delete(c *gin.Context) {
	id := c.Param("table_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	if err := h.tableService.Delete(uuid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "table deleted successfully"})
}

func (h *TableHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("table_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var status models.TableStatus
	switch req.Status {
	case string(models.TableStatusFree):
		status = models.TableStatusFree
	case string(models.TableStatusOccupied):
		status = models.TableStatusOccupied
	case string(models.TableStatusReserved):
		status = models.TableStatusReserved
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid status"})
		return
	}

	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid table ID"})
		return
	}

	if err := h.tableService.UpdateStatus(uuid, status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "table status updated successfully"})
}
