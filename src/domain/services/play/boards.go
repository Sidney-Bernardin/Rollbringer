package play

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/src/domain"
)

const (
	ExternalErrorTypeInvalidBoardName domain.ExternalErrorType = "invalid_board_name"
)

type Board struct {
	ID             uuid.UUID
	Name           BoardName
	Canvas         []byte
	UserPermisions map[uuid.UUID][]BoardUserPermision
}

type BoardName string

func ParseBoardName(str string) (BoardName, error) {
	if len(str) == 0 || 30 < len(str) {
		return "", &domain.ExternalError{
			Type: ExternalErrorTypeInvalidBoardName,
			Msg:  "Must be between 1 and 30 characters",
			Details: map[string]any{
				"board_name": str,
			},
		}
	}

	return BoardName(str), nil
}

type BoardUserPermision string

const (
	BoardUserPermisionOwner BoardUserPermision = "OWNER"
	BoardUserPermisionEdit  BoardUserPermision = "EDIT"
)

type CreateBoardOpts struct {
	Name    string      `json:"name"`
	UserIDs []uuid.UUID `json:"users_ids"`
}

func (svc *service) CreateBoard(ctx context.Context, args *CreateBoardOpts, creator *domain.PublicUser, users []domain.PublicUser) (*Board, error) {

	name, err := ParseBoardName(args.Name)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse name")
	}

	board := &Board{
		ID:     uuid.New(),
		Name:   name,
		Canvas: []byte(`{}`),
		UserPermisions: map[uuid.UUID][]BoardUserPermision{
			creator.UserID: {BoardUserPermisionOwner},
		},
	}

	if err := svc.database.CreateBoard(ctx, board); err != nil {
		return nil, errors.Wrap(err, "database cannot create board")
	}

	svc.broker.Pub(ctx, &domain.EventNewBoard{
		BoardID: board.ID,
		Name:    string(board.Name),
		Users:   append(users, *creator),
	})

	return board, nil
}
