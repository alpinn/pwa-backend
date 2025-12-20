package models

import "time"

type Transaction struct {
	ID            string             `json:"id"`
	UserID        string             `json:"user_id"`
	TotalAmount   float64            `json:"total_amount"`
	Status        string             `json:"status"` // pending, completed, cancelled
	CreatedAt     time.Time          `json:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at"`
	Items         []TransactionItem  `json:"items,omitempty"`
}

type TransactionItem struct {
	ID            string    `json:"id"`
	TransactionID string    `json:"transaction_id"`
	ProductID     string    `json:"product_id"`
	ProductName   string    `json:"product_name"`
	Quantity      int       `json:"quantity"`
	Price         float64   `json:"price"`
	Subtotal      float64   `json:"subtotal"`
	UserID        string    `json:"user_id"`
	CreatedAt     time.Time `json:"created_at"`
}

type CheckoutRequest struct {
	Items []CheckoutItem `json:"items" binding:"required"`
}

type CheckoutItem struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}


type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending completed cancelled"`
}