package repositories

import (
	"database/sql"
	"pwa-backend/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, password, name, role, created_at, updated_at 
	          FROM users WHERE username = $1`
	
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Name, 
		&user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return &user, nil
}

