package models

import "time"

type StockEvent struct {
	ID            string    `json:"id"`
	ProductID     string    `json:"product_id"`
	Qty           int       `json:"qty"` // Positive or negative
	Type          string    `json:"type"` // sale, restock, reject, adjustment, opening_stock
	Source        string    `json:"source"` // pos, dashboard, online
	TransactionID *string   `json:"transaction_id,omitempty"`
	UserID        *string   `json:"user_id,omitempty"`
	DeviceID      *string   `json:"device_id,omitempty"`
	Note          string    `json:"note"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateStockEventRequest struct {
	ProductID string  `json:"product_id" binding:"required"`
	Qty       int     `json:"qty" binding:"required"`
	Type      string  `json:"type" binding:"required,oneof=sale restock reject adjustment opening_stock"`
	Source    string  `json:"source" binding:"required,oneof=pos dashboard online"`
	DeviceID  *string `json:"device_id"`
	Note      string  `json:"note"`
}