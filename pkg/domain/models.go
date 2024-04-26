package domain

import "github.com/google/uuid"

type PlayPage struct {
	LoggedIn bool
	IsHost   bool

	User *User
	Game *Game
}

// =====

type UserView string

const (
	UserViewDefault UserView = "default"
)

type User struct {
	ID uuid.UUID `json:"id,omitempty"`

	GoogleID *string `json:"google_id,omitempty"`
	Username string  `json:"username,omitempty"`

	PDFs        []*PDF  `json:"pdfs,omitempty"`
	HostedGames []*Game `json:"hosted_games,omitempty"`
	JoinedGames []*Game `json:"joined_games,omitempty"`
}

// =====

type SessionView string

const (
	SessionViewDefault SessionView = "default"
)

type Session struct {
	ID uuid.UUID `json:"id,omitempty"`

	UserID    uuid.UUID `json:"user_id,omitempty"`
	CSRFToken string    `json:"csrf_token,omitempty"`
}

// =====

type GameView string

const (
	GameViewDefault  GameView = "default"
	GameViewWithHost GameView = "with_host"
)

type Game struct {
	ID uuid.UUID `json:"id,omitempty"`

	HostID uuid.UUID `json:"host_id,omitempty"`
	Host   *User     `json:"host,omitempty"`

	Name string `json:"name,omitempty"`

	Players []*User `json:"players,omitempty"`
	PDFs    []*PDF  `json:"pdfs,omitempty"`
}

// =====

type PDFView string

const (
	PDFViewDefault   PDFView = "default"
	PDFViewWithOwner PDFView = "with_owner"
	PDFViewWithGame  PDFView = "with_game"
)

type PDF struct {
	ID uuid.UUID `json:"id,omitempty"`

	OwnerID uuid.UUID `json:"owner_id,omitempty"`
	Owner   *User     `json:"owner,omitempty"`

	GameID *uuid.UUID `json:"game_id,omitempty"`
	Game   *Game      `json:"game,omitempty"`

	Name   string              `json:"name,omitempty"`
	Schema string              `json:"schema,omitempty"`
	Pages  []map[string]string `json:"pages,omitempty"`
}
