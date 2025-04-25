package services

import (
	"context"

	"rollbringer/src"
	account_models "rollbringer/src/services/accounts/models"
	play_models "rollbringer/src/services/play/models"
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

	EventBoardShared struct {
		ID    src.UUID              `json:"board_id"`
		Name  play_models.BoardName `json:"name"`
		Users map[src.UUID]struct {
			Username       account_models.Username `json:"username"`
			ProfilePicture string                  `json:"profile_picture"`
		} `json:"users"`
	}
)
