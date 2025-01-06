package database

import (
	"github.com/google/uuid"
	"github.com/lib/pq"

	"rollbringer/internal"
)

type User struct {
	ID       uuid.UUID `db:"id"`
	Username string    `db:"username"`

	GoogleID      *string `db:"google_id"`
	GooglePicture *string `db:"google_picture"`
}

func (user *User) Internalized() *internal.User {
	if user != nil {
		return &internal.User{
			ID:            user.ID,
			Username:      user.Username,
			GoogleID:      user.GoogleID,
			GooglePicture: user.GooglePicture,
		}
	}
	return nil
}

type Session struct {
	ID uuid.UUID `db:"id"`

	UserID uuid.UUID `db:"user_id"`
	User   *User     `db:"user"`

	CSRFToken string `db:"csrf_token"`
}

func (session *Session) Internalized() *internal.Session {
	if session != nil {
		return &internal.Session{
			ID:        session.ID,
			UserID:    session.UserID,
			User:      session.User.Internalized(),
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
			Owner:   pdf.Owner.Internalized(),
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

	DiceTypes   pq.Int32Array `db:"dice_types"`
	DiceResults pq.Int32Array `db:"dice_results"`
	Modifiers   string        `db:"modifiers"`
}

func (roll *Roll) Internalized() *internal.Roll {
	if roll != nil {
		return &internal.Roll{
			ID:          roll.ID,
			OwnerID:     roll.OwnerID,
			Owner:       roll.Owner.Internalized(),
			GameID:      roll.GameID,
			Game:        roll.Game.Internalized(),
			DiceTypes:   roll.DiceTypes,
			DiceResults: roll.DiceResults,
			Modifiers:   roll.Modifiers,
		}
	}
	return nil
}

type ChatMessage struct {
	ID uuid.UUID `db:"id"`

	OwnerID uuid.UUID `db:"owner_id"`
	Owner   *User     `db:"owner"`

	GameID uuid.UUID `db:"game_id"`
	Game   *Game     `db:"game"`

	Message string `db:"message"`
}

func (chatMessage *ChatMessage) Internalized() *internal.ChatMessage {
	if chatMessage != nil {
		return &internal.ChatMessage{
			ID:      chatMessage.ID,
			OwnerID: chatMessage.OwnerID,
			Owner:   chatMessage.Owner.Internalized(),
			GameID:  chatMessage.GameID,
			Game:    chatMessage.Game.Internalized(),
			Message: chatMessage.Message,
		}
	}
	return nil
}
