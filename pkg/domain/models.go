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

type Event struct {
	Headers map[string]any `json:"HEADERS"`
	Type    string         `json:"TYPE"`

	PDFID     string            `json:"pdf_id,omitempty"`
	SenderID  string            `json:"sender_ID,omitempty"`
	PageNum   int               `json:"page_num,string,omitempty"`
	PDFFields map[string]string `json:"pdf_fields"`

	Modifier string `json:"modifier,omitempty"`
}
