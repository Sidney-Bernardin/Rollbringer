package domain

import "github.com/google/uuid"

type PlayPage struct {
	LoggedIn bool

	User      *User
	UserGames []*Game
	UserPDFs  []*PDF

	Game     *Game
	GamePDFs []*PDF
}

type User struct {
	ID uuid.UUID `json:"id,omitempty"`

	GoogleID string `json:"google_id,omitempty"`
	Username string `json:"username,omitempty"`
}

type Session struct {
	ID uuid.UUID `json:"id,omitempty"`

	UserID uuid.UUID `json:"user_id,omitempty"`
	User   *User     `json:"user,omitempty"`

	CSRFToken string `json:"csrf_token,omitempty"`
}

type Game struct {
	ID uuid.UUID `json:"id,omitempty"`

	HostID uuid.UUID `json:"host_id,omitempty"`
	Host   *User     `json:"host,omitempty"`

	Title string `json:"title,omitempty"`
}

type PDF struct {
	ID uuid.UUID `json:"id,omitempty"`

	OwnerID uuid.UUID `json:"owner_id,omitempty"`
	Owner   *User     `json:"owner,omitempty"`

	GameID uuid.UUID `json:"game_id,omitempty"`
	Game   *Game     `json:"game,omitempty"`

	Name   string              `json:"name,omitempty"`
	Schema string              `json:"schema,omitempty"`
	Fields []map[string]string `json:"pages,omitempty"`
}
