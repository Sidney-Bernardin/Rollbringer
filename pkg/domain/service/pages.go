package service

import (
	"context"
	"fmt"
	"rollbringer/pkg/domain"
)

func (svc *Service) Play(ctx context.Context, gameID string, incomingChan, outgoingChan chan domain.GameEvent) {

	// Get the game.
	game, err := svc.GetGame(ctx, gameID)
	if err != nil && err != domain.ErrGameNotFound {
		svc.logger.Error().Stack().Err(err).Msg("Cannot get game")
		return
	}

	if game != nil {
		go svc.ps.SubToGame(ctx, game.ID, incomingChan)
	}

	for {
		select {

		case <-ctx.Done():
			return

		case event := <-incomingChan:
			fmt.Println(event)

			switch event["type"] {
			case "pdf_update":
			}
		}
	}
}
