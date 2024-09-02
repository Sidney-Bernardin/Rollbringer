package internal

import (
	"github.com/google/uuid"
)

type GoogleUserInfo struct {
	GoogleID  string
	GivenName string
}

// =====

type UserView int

const (
	UserViewAll UserView = iota
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

type SessionView int

const (
	SessionViewAll SessionView = iota
)

type Session struct {
	ID uuid.UUID `json:"id,omitempty"`

	UserID    uuid.UUID `json:"user_id,omitempty"`
	CSRFToken string    `json:"csrf_token,omitempty"`
}

// =====

type GameView int

const (
	GameViewAll GameView = iota
	GameViewAll_HostInfo
)

var GameViews = map[string]GameView{
	"All":          GameViewAll,
	"All_HostInfo": GameViewAll_HostInfo,
}

type Game struct {
	ID uuid.UUID `json:"id,omitempty"`

	HostID uuid.UUID `json:"host_id,omitempty"`
	Host   *User     `json:"host,omitempty"`

	Name string `json:"name,omitempty"`

	Players []*User `json:"players,omitempty"`
	PDFs    []*PDF  `json:"pdfs,omitempty"`
	Rolls   []*Roll `json:"rolls,omitempty"`
}

// =====

type PDFView int

const (
	PDFViewAll PDFView = iota
	PDFViewAll_OwnerInfo_GameInfo
	PDFViewAll_OwnerInfo
	PDFViewAll_GameInfo
)

type PDF struct {
	ID uuid.UUID `json:"id,omitempty"`

	OwnerID uuid.UUID `json:"owner_id,omitempty"`
	Owner   *User     `json:"owner,omitempty"`

	GameID *uuid.UUID `json:"game_id,omitempty"`
	Game   *Game      `json:"game,omitempty"`

	Name   string              `json:"name,omitempty"`
	Schema string              `json:"schema,omitempty"`
	Fields []map[string]string `json:"fields,omitempty"`
}

// =====

type RollView int

const (
	RollViewAll RollView = iota
)

type Roll struct {
	ID uuid.UUID `json:"id,omitempty"`

	OwnerID uuid.UUID `json:"owner_id,omitempty"`
	Owner   *User     `json:"owner,omitempty"`

	GameID uuid.UUID `json:"game_id,omitempty"`
	Game   *Game     `json:"game,omitempty"`

	DiceNames   []int32 `json:"dice_names,omitempty"`
	DiceResults []int32 `json:"dice_results,omitempty"`
}
