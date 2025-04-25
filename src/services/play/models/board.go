package models

import (
	"rollbringer/src"

	"github.com/pkg/errors"
)

const (
	ExternalErrorTypeInvalidBoardName src.ExternalErrorType = "invalid_room_name"
)

type Board struct {
	ID             src.UUID                              `json:"id"`
	Name           BoardName                             `json:"name"`
	Canvas         []byte                                `json:"canvas"`
	UserPermisions map[src.UUID][]src.BoardUserPermision `json:"user_permisions"`
}

func NewBoard(creatorID src.UUID, name string) (*Board, error) {
	boardName, err := ParseBoardName(name)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse board-name")
	}

	return &Board{
		ID:     src.NewUUID(),
		Name:   boardName,
		Canvas: []byte(`{}`),
		UserPermisions: map[src.UUID][]src.BoardUserPermision{
			creatorID: {src.BoardUserPermisionOwner},
		},
	}, nil
}

type BoardName string

func ParseBoardName(str string) (BoardName, error) {
	if len(str) == 0 || 30 < len(str) {
		return "", &src.ExternalError{
			Type: ExternalErrorTypeInvalidBoardName,
			Msg:  "Must be between 1 and 30 characters",
			Details: map[string]any{
				"board_name": str,
			},
		}
	}

	return BoardName(str), nil
}
