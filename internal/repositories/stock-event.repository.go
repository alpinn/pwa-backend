package repositories

import (
	"database/sql"
	"pwa-backend/internal/models"
)

type StockEventRepository struct {
	db *sql.DB
}

func NewStockEventRepository(db *sql.DB) *StockEventRepository {
	return &StockEventRepository{db: db}
}

func (r *StockEventRepository) Create(tx *sql.Tx, event *models.StockEvent) error {
	query := `
		INSERT INTO stock_events 
		(id, product_id, qty, type, source, transaction_id, user_id, device_id, note, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	
	_, err := tx.Exec(query,
		event.ID,
		event.ProductID,
		event.Qty,
		event.Type,
		event.Source,
		event.TransactionID,
		event.UserID,
		event.DeviceID,
		event.Note,
		event.CreatedAt,
	)
	
	return err
}

func (r *StockEventRepository) GetByProduct(productID string, limit int) ([]models.StockEvent, error) {
	query := `
		SELECT id, product_id, qty, type, source, transaction_id, user_id, device_id, note, created_at
		FROM stock_events
		WHERE product_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`
	
	rows, err := r.db.Query(query, productID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.StockEvent
	for rows.Next() {
		var e models.StockEvent
		err := rows.Scan(
			&e.ID, &e.ProductID, &e.Qty, &e.Type, &e.Source,
			&e.TransactionID, &e.UserID, &e.DeviceID, &e.Note, &e.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (r *StockEventRepository) GetAll(limit, offset int) ([]models.StockEvent, error) {
	query := `
		SELECT id, product_id, qty, type, source, transaction_id, user_id, device_id, note, created_at
		FROM stock_events
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []models.StockEvent
	for rows.Next() {
		var e models.StockEvent
		err := rows.Scan(
			&e.ID, &e.ProductID, &e.Qty, &e.Type, &e.Source,
			&e.TransactionID, &e.UserID, &e.DeviceID, &e.Note, &e.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (r *StockEventRepository) BeginTx() (*sql.Tx, error) {
	return r.db.Begin()
}