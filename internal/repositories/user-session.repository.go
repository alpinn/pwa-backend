package repositories

import (
	"database/sql"
	"time"

	"pwa-backend/internal/models"
)

type UserSessionRepository struct {
	db *sql.DB
}

func NewUserSessionRepository(db *sql.DB) *UserSessionRepository {
	return &UserSessionRepository{db: db}
}

func (r *UserSessionRepository) Create(session *models.UserSession) error {
	query := `
		INSERT INTO user_sessions 
		(id, user_id, refresh_token, device_id, user_agent, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err := r.db.Exec(
		query,
		session.ID,
		session.UserID,
		session.RefreshToken,
		session.DeviceID,
		session.UserAgent,
		session.ExpiresAt,
		session.CreatedAt,
	)
	
	return err
}

func (r *UserSessionRepository) GetByRefreshToken(refreshToken string) (*models.UserSession, error) {
	var session models.UserSession
	query := `
		SELECT id, user_id, refresh_token, device_id, user_agent, expires_at, revoked_at, created_at
		FROM user_sessions
		WHERE refresh_token = $1 AND revoked_at IS NULL
	`
	
	err := r.db.QueryRow(query, refreshToken).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.DeviceID,
		&session.UserAgent,
		&session.ExpiresAt,
		&session.RevokedAt,
		&session.CreatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	
	return &session, nil
}

func (r *UserSessionRepository) GetByUserID(userID string) ([]*models.UserSession, error) {
	query := `
		SELECT id, user_id, refresh_token, device_id, user_agent, expires_at, revoked_at, created_at
		FROM user_sessions
		WHERE user_id = $1 AND revoked_at IS NULL ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var sessions []*models.UserSession
	for rows.Next() {
		var session models.UserSession
		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.RefreshToken,
			&session.DeviceID,
			&session.UserAgent,
			&session.ExpiresAt,
			&session.RevokedAt,
			&session.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}
	
	return sessions, rows.Err()
}

func (r *UserSessionRepository) RevokeByID(sessionID string) error {
	query := `UPDATE user_sessions SET revoked_at = $1 WHERE id = $2`
	_, err := r.db.Exec(query, time.Now(), sessionID)
	return err
}

func (r *UserSessionRepository) RevokeAllByUserID(userID string) error {
	query := `UPDATE user_sessions SET revoked_at = $1 WHERE user_id = $2 AND revoked_at IS NULL`
	_, err := r.db.Exec(query, time.Now(), userID)
	return err
}

func (r *UserSessionRepository) IsSessionValid(sessionID string) (bool, error) {
	var expiresAt time.Time
	var revokedAt *time.Time
	
	query := `
		SELECT expires_at, revoked_at FROM user_sessions WHERE id = $1
	`
	
	err := r.db.QueryRow(query, sessionID).Scan(&expiresAt, &revokedAt)
	
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	
	// Session is valid if it's not revoked and not expired
	isValid := revokedAt == nil && time.Now().Before(expiresAt)
	return isValid, nil
}
