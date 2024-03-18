package domain

type PlayPage struct {
	LoggedIn bool

	User *User
	Game *Game
}

type User struct {
	ID       string
	GoogleID string

	Username string
}

type Session struct {
	ID        string
	CSRFToken string
	UserID    string
}

type Game struct {
	ID     string
	HostID string

	Title string
	PDFs  []string
}

type PDF struct {
	ID      string
	OwnerID string
	GameID  string

	Name   string
	Schema string
	Pages  []map[string]string
}
