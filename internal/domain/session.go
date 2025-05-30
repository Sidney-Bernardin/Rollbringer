package domain

import "github.com/google/uuid"

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CSRFToken string    `json:"csrf_token"`
}
