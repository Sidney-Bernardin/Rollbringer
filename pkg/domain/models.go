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
	Name    string `db:"name,omitempty"`
	Schema  string `db:"schema,omitempty"`
	Content []byte `db:"content,omitempty"`
}

type GameEvent map[string]any
