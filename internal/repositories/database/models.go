package databases

import (
	"github.com/google/uuid"
	"github.com/lib/pq"

	"rollbringer/internal"
)

type User struct {
	ID uuid.UUID `db:"id"`

	GoogleID *string `db:"google_id"`
	Username string  `db:"username"`
}

func (user *User) Internalized() *internal.User {
	if user != nil {
		return &internal.User{
			ID:       user.ID,
			GoogleID: user.GoogleID,
			Username: user.Username,
		}
	}
	return nil
}

type Session struct {
	ID uuid.UUID `db:"id"`

	UserID    uuid.UUID `db:"user_id"`
	CSRFToken string    `db:"csrf_token"`
}

func (session *Session) Internalized() *internal.Session {
	if session != nil {
		return &internal.Session{
			ID:        session.ID,
			UserID:    session.UserID,
			CSRFToken: session.CSRFToken,
		}
	}
	return nil
}

type Game struct {
	ID uuid.UUID `db:"id"`

	HostID uuid.UUID `db:"host_id"`
	Host   *User     `db:"host"`

	Name string `db:"name"`
}

func (game *Game) Internalized() *internal.Game {
	if game != nil {
		return &internal.Game{
			ID:     game.ID,
			HostID: game.HostID,
			Name:   game.Name,
		}
	}
	return nil
}

type PDF struct {
	ID uuid.UUID `db:"id"`

	OwnerID uuid.UUID `db:"owner_id"`
	Owner   *User     `db:"owner"`

	GameID *uuid.UUID `db:"game_id"`
	Game   *Game      `db:"game"`

	Name   string          `db:"name"`
	Schema string          `db:"schema"`
	Pages  pq.GenericArray `db:"pages"`
}

func (pdf *PDF) Internalized() *internal.PDF {
	if pdf != nil {
		return &internal.PDF{
			ID:      pdf.ID,
			OwnerID: pdf.OwnerID,
			GameID:  pdf.GameID,
			Game:    pdf.Game.Internalized(),
			Name:    pdf.Name,
			Schema:  pdf.Schema,
		}
	}
	return nil
}

type Roll struct {
	ID uuid.UUID `db:"id"`

	OwnerID uuid.UUID `db:"owner_id"`
	Owner   *User     `db:"owner"`

	GameID uuid.UUID `db:"game_id"`
	Game   *Game     `db:"game"`

	DiceNames   pq.Int32Array `db:"dice_names"`
	DiceResults pq.Int32Array `db:"dice_results"`
}

func (roll *Roll) Internalized() *internal.Roll {
	if roll != nil {
		return &internal.Roll{
			ID:          roll.ID,
			OwnerID:     roll.OwnerID,
			GameID:      roll.GameID,
			Game:        roll.Game.Internalized(),
			DiceNames:   roll.DiceNames,
			DiceResults: roll.DiceResults,
		}
	}
	return nil
}
