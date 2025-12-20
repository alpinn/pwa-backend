package handlers

import (
	"net/http"
	"pwa-backend/internal/models"
	"pwa-backend/internal/repositories"
	"time"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionRepo *repositories.TransactionRepository
	productRepo     *repositories.ProductRepository
}

func NewTransactionHandler(transactionRepo *repositories.TransactionRepository, productRepo *repositories.ProductRepository) *TransactionHandler {
	return &TransactionHandler{
		transactionRepo: transactionRepo,
		productRepo:     productRepo,
	}
}

// Checkout godoc
// @Summary Checkout transaction
// @Description Create new transaction with stock validation
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.CheckoutRequest true "Checkout items"
// @Success 201 {object} models.Transaction
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transactions/checkout [post]
func (h *TransactionHandler) Checkout(c *gin.Context) {
	var req models.CheckoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetString("user_id")

	tx, err := h.transactionRepo.BeginTx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback()

	var totalAmount float64
	var items []models.TransactionItem
	transactionID := generateID()

	for _, item := range req.Items {
		product, err := h.productRepo.GetByID(item.ProductID)
		if err != nil || product == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Product not found: " + item.ProductID})
			return
		}

		if product.Stock < item.Quantity {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock for product: " + product.Name})
			return
		}

		if err := h.productRepo.DeductStock(tx, product.ID, item.Quantity); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock"})
			return
		}

		subtotal := product.Price * float64(item.Quantity)
		totalAmount += subtotal

		items = append(items, models.TransactionItem{
			ID:            generateID(),
			TransactionID: transactionID,
			ProductID:     product.ID,
			ProductName:   product.Name,
			Quantity:      item.Quantity,
			Price:         product.Price,
			Subtotal:      subtotal,
			UserID:        userID,
		})
	}

	transaction := &models.Transaction{
		ID:          transactionID,
		UserID:      userID,
		TotalAmount: totalAmount,
		Status:      "pending",
	}

	if err := h.transactionRepo.Create(tx, transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction"})
		return
	}

	if err := h.transactionRepo.CreateItems(tx, items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction items"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	transaction.Items = items
	c.JSON(http.StatusCreated, transaction)
}

// GetTransaction godoc
// @Summary Get transaction by ID
// @Description Get transaction with all items
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID"
// @Success 200 {object} models.Transaction
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /transactions/{id} [get]
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	id := c.Param("id")

	transaction, err := h.transactionRepo.GetByIDWithItems(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction"})
		return
	}

	if transaction == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

// UpdateStatus godoc
// @Summary Update transaction status
// @Description Update status of existing transaction
// @Tags transactions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Transaction ID"
// @Param request body models.UpdateStatusRequest true "Status update"
// @Success 200 {object} models.Transaction
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /transactions/{id}/status [put]
func (h *TransactionHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var req models.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transaction, err := h.transactionRepo.GetByIDWithItems(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch transaction"})
		return
	}

	if transaction == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	if err := h.transactionRepo.UpdateStatus(id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	transaction.Status = req.Status
	c.JSON(http.StatusOK, transaction)
}

func generateID() string {
	return time.Now().Format("20060102150405") + randomString(6)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}