package models

import "time"

type UserSession struct {
	ID           string     `json:"id"`
	UserID       string     `json:"user_id"`
	RefreshToken string     `json:"-"`
	DeviceID     *string    `json:"device_id"`
	UserAgent    *string    `json:"user_agent"`
	ExpiresAt    time.Time  `json:"expires_at"`
	RevokedAt    *time.Time `json:"revoked_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
