package games

import (
	"context"
	"rollbringer/internal"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

func (svc *service) CreateGame(ctx context.Context, session *internal.Session, game *internal.Game) error {
	count, err := svc.schema.GamesCount(ctx, session.UserID)
	if err != nil {
		return errors.Wrap(err, "cannot count games")
	}

	if count >= 15 {
		return internal.NewProblemDetail(ctx, internal.PDOpts{
			Type:   internal.PDTypeMaxGames,
			Detail: "You cannot host more than 15 games at a time.",
		})
	}

	game.HostID = session.UserID
	err = svc.schema.GameInsert(ctx, game)
	return errors.Wrap(err, "cannot insert game")
}

func (svc *service) DeleteGame(ctx context.Context, session *internal.Session, gameID uuid.UUID) error {
	err := svc.schema.GameDelete(ctx, gameID, session.UserID)
	return errors.Wrap(err, "cannot delete game")
}
