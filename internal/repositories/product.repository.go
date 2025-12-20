package repositories

import (
	"database/sql"
	"pwa-backend/internal/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll() ([]models.Product, error) {
	query := `SELECT id, name, description, price, stock, image_url, created_at, updated_at 
	          FROM products ORDER BY name`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepository) GetByID(id string) (*models.Product, error) {
	var p models.Product
	query := `SELECT id, name, description, price, stock, image_url, created_at, updated_at 
	          FROM products WHERE id = $1`
	
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Price, &p.Stock, &p.ImageURL, &p.CreatedAt, &p.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return &p, nil
}

func (r *ProductRepository) DeductStock(tx *sql.Tx, productID string, quantity int) error {
	query := `UPDATE products SET stock = stock - $1, updated_at = NOW() WHERE id = $2`
	_, err := tx.Exec(query, quantity, productID)
	return err
}