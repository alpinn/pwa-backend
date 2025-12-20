package handlers

import (
	"net/http"
	"pwa-backend/internal/repositories"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productRepo *repositories.ProductRepository
}

func NewProductHandler(productRepo *repositories.ProductRepository) *ProductHandler {
	return &ProductHandler{productRepo: productRepo}
}

// GetProducts godoc
// @Summary Get all products
// @Description Get list of all products (read-only)
// @Tags products
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Product
// @Failure 500 {object} map[string]string
// @Router /products [get]
func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.productRepo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch products"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get single product by ID (read-only)
// @Tags products
// @Produce json
// @Security BearerAuth
// @Param id path string true "Product ID"
// @Success 200 {object} models.Product
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	product, err := h.productRepo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product"})
		return
	}

	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}