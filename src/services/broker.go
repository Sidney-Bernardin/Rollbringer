package services

import (
	"context"

	"rollbringer/src"
	account_models "rollbringer/src/services/accounts/models"
)

type Broker interface {
	PubChat(event *EventChat) error
	SubRoom(ctx context.Context, roomID src.UUID, cb func(event any)) error
}

type (
	EventChat struct {
		AuthorID       src.UUID                `json:"author_id"`
		Username       account_models.Username `json:"username"`
		ProfilePicture string                  `json:"profile_picture"`

		RoomID  src.UUID `json:"room_id"`
		Message string   `json:"message"`
	}
)
