package domain

type User struct {
	ID       string `db:"id,omitempty"`
	GoogleID string `db:"google_id,omitempty"`
	Username string `db:"username,omitempty"`
}

type Session struct {
	ID        string `db:"id,omitempty"`
	CSRFToken string `db:"csrf_token,omitempty"`
	UserID    string `db:"user_id,omitempty"`
}

type Game struct {
	ID     string `db:"id,omitempty"`
	HostID string `db:"host_id,omitempty"`
	Title  string `db:"title,omitempty"`
}

type PDF struct {
	ID      string `db:"id,omitempty"`
	OwnerID string `db:"owner_id,omitempty"`

	Name   string `db:"name,omitempty"`
	Schema string `db:"schema,omitempty"`

	MainPage   []byte `db:"main_page,omitempty"`
	InfoPage   []byte `db:"info_page,omitempty"`
	SpellsPage []byte `db:"spells_page,omitempty"`
}

type GameEvent struct {
	Headers map[string]string `json:"HEADERS"`
	Type    string            `json:"TYPE"`
	Body    map[string]string
}
