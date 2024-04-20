package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/pkg/domain"
)

func (svc *Service) CreateGame(ctx context.Context, session *domain.Session, game *domain.Game) error {
	game.HostID = session.UserID

	count, err := svc.DB.GamesCount(ctx, session.UserID)
	if err != nil {
		return errors.Wrap(err, "cannot get games count")
	}

	if count >= 5 {
		return &domain.ProblemDetail{
			Type:   domain.PDTypeMaxGames,
			Detail: "You cannot host more than 5 games at a time.",
		}
	}

	err = svc.DB.InsertGame(ctx, game)
	return errors.Wrap(err, "cannot insert game")
}

func (svc *Service) DeleteGame(ctx context.Context, session *domain.Session, gameID uuid.UUID) error {
	err := svc.DB.DeleteGame(ctx, gameID, session.UserID)
	return errors.Wrap(err, "cannot delete game")
}
