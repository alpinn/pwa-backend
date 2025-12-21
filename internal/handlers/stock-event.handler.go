package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"pwa-backend/internal/mathutil"
	"pwa-backend/internal/models"
	"pwa-backend/internal/repositories"
)

type StockEventHandler struct {
	stockEventRepo *repositories.StockEventRepository
	productRepo    *repositories.ProductRepository
}

func NewStockEventHandler(stockEventRepo *repositories.StockEventRepository, productRepo *repositories.ProductRepository) *StockEventHandler {
	return &StockEventHandler{
		stockEventRepo: stockEventRepo,
		productRepo:    productRepo,
	}
}

// CreateStockEvent godoc
// @Summary Create stock event
// @Description Create stock event (restock, adjustment, etc). Stock will be updated automatically.
// @Tags stock_events
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CreateStockEventRequest true "Stock event data"
// @Success 201 {object} models.StockEvent
// @Failure 400 {object} map[string]string
// @Router /stock-events [post]
func (h *StockEventHandler) CreateStockEvent(c *gin.Context) {
	var req models.CreateStockEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")

	product, err := h.productRepo.GetByID(req.ProductID)
	if err != nil || product == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found"})
		return
	}

	qty := req.Qty
	switch req.Type {
	case "sale", "reject":
		qty = -mathutil.Abs(req.Qty)
	case "restock", "opening_stock":
		qty = mathutil.Abs(req.Qty)
	}

	stockEvent := &models.StockEvent{
		ID:        generateID(),
		ProductID: req.ProductID,
		Qty:       qty,
		Type:      req.Type,
		Source:    req.Source,
		UserID:    &userID,
		DeviceID:  req.DeviceID,
		Note:      req.Note,
		CreatedAt: time.Now(),
	}

	tx, err := h.stockEventRepo.BeginTx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	if err := h.stockEventRepo.Create(tx, stockEvent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stock event: " + err.Error()})
		return
	}

	if err := h.productRepo.UpdateStockByQty(tx, req.ProductID, qty); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product stock: " + err.Error()})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit"})
		return
	}

	c.JSON(http.StatusCreated, stockEvent)
}

// GetStockEventsByProduct godoc
// @Summary Get stock events by product
// @Description Get stock event history for a specific product
// @Tags stock_events
// @Produce json
// @Security BearerAuth
// @Param product_id path string true "Product ID"
// @Param limit query int false "Limit" default(50)
// @Success 200 {array} models.StockEvent
// @Router /stock-events/product/{product_id} [get]
func (h *StockEventHandler) GetStockEventsByProduct(c *gin.Context) {
	productID := c.Param("product_id")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	events, err := h.stockEventRepo.GetByProduct(productID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock events"})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetAllStockEvents godoc
// @Summary Get all stock events
// @Description Get all stock events with pagination
// @Tags stock_events
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Limit" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} models.StockEvent
// @Router /stock-events [get]
func (h *StockEventHandler) GetAllStockEvents(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	events, err := h.stockEventRepo.GetAll(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock events"})
		return
	}

	c.JSON(http.StatusOK, events)
}