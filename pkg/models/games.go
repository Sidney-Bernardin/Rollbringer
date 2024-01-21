package models

import "github.com/google/uuid"

type Game struct {
	ID     uuid.UUID
	HostID uuid.UUID
	Title  string
}

type GameEvent struct {
	Type string `json:"type"`
}
