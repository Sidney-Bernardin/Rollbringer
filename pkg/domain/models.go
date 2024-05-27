package domain

import "github.com/google/uuid"

type PlayPage struct {
	LoggedIn bool
	IsHost   bool

	User *User
	Game *Game
}

// =====

type UserView int

const (
	UserViewAll UserView = iota
)

var UserViews = map[string]UserView{
	"All": UserViewAll,
}

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

var SessionViews = map[string]SessionView{
	"All": SessionViewAll,
}

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
}

// =====

type PDFView int

const (
	PDFViewAll PDFView = iota
	PDFViewAll_OwnerInfo_GameInfo
	PDFViewAll_OwnerInfo
	PDFViewAll_GameInfo
)

var PDFViews = map[string]PDFView{
	"All":                    PDFViewAll,
	"All_OwnerInfo_GameInfo": PDFViewAll_OwnerInfo_GameInfo,
	"All_OwnerInfo":          PDFViewAll_OwnerInfo,
	"All_GameInfo":           PDFViewAll_GameInfo,
}

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
