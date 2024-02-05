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

	Name   string   `db:"name,omitempty"`
	Schema string   `db:"schema,omitempty"`
	Pages  []string `db:"pages,omitempty"`
}

type GameEvent struct {
	Headers map[string]any `json:"HEADERS"`
	Type    string         `json:"TYPE"`

	SenderID string `json:"sender_id"`
	PDFID    string `json:"pdf_id"`
	PageNum  int    `json:"page_num,string"`

	PDFFields map[string]string `json:"-"`
}
