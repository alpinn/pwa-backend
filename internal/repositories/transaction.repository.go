package repositories

import (
	"database/sql"
	"pwa-backend/internal/models"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *sql.Tx, transaction *models.Transaction) error {
	query := `INSERT INTO transactions (id, user_id, total_amount, status, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, NOW(), NOW())`

	_, err := tx.Exec(query, transaction.ID, transaction.UserID, transaction.TotalAmount, transaction.Status)
	return err
}

func (r *TransactionRepository) CreateItems(tx *sql.Tx, items []models.TransactionItem) error {
	query := `INSERT INTO transaction_items (id, transaction_id, product_id, product_name, quantity, price, subtotal, user_id, created_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())`

	for _, item := range items {
		_, err := tx.Exec(query, item.ID, item.TransactionID, item.ProductID, item.ProductName, item.Quantity, item.Price, item.Subtotal, item.UserID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *TransactionRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}

func (r *TransactionRepository) UpdateStatus(id, status string) error {
	query := `UPDATE transactions SET status = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *TransactionRepository) GetByID(id string) (*models.Transaction, error) {
	var t models.Transaction
	query := `SELECT id, user_id, total_amount, status, created_at, updated_at 
	          FROM transactions WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&t.ID, &t.UserID, &t.TotalAmount, &t.Status, &t.CreatedAt, &t.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TransactionRepository) GetByIDWithItems(id string) (*models.Transaction, error) {
	var t models.Transaction
	query := `SELECT id, user_id, total_amount, status, created_at, updated_at 
	          FROM transactions WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&t.ID, &t.UserID, &t.TotalAmount, &t.Status, &t.CreatedAt, &t.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Fetch items
	itemsQuery := `SELECT id, transaction_id, product_id, product_name, quantity, price, subtotal, user_id, created_at 
	              FROM transaction_items WHERE transaction_id = $1 ORDER BY created_at DESC`
	
	rows, err := r.db.Query(itemsQuery, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.TransactionItem
		err := rows.Scan(
			&item.ID, &item.TransactionID, &item.ProductID, &item.ProductName, 
			&item.Quantity, &item.Price, &item.Subtotal, &item.UserID, &item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		t.Items = append(t.Items, item)
	}

	return &t, nil
}