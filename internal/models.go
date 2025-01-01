package internal

import (
	"fmt"

	"github.com/google/uuid"
)

type Event string

type EventWrapper[T any] struct {
	Event   Event `json:"event"`
	Payload T
}

// =====

const EventError Event = "ERROR"

type ProblemDetail struct {
	Instance string `json:"instance,omitempty"`
	Type     PDType `json:"type"`
	Detail   string `json:"detail,omitempty"`
	Extra    map[string]any
}

func (pd *ProblemDetail) Error() string {
	return fmt.Sprintf(`%s: %s`, pd.Type, pd.Detail)
}

// =====

const EventPlayPage Event = "PLAY_PAGE"

type PlayPage struct {
	Session *Session `json:"session"`
	Game    *Game    `json:"game"`
}

// =====

type GoogleUserInfo struct {
	GoogleID  string
	GivenName string
	Picture   string
}

// =====

const (
	EventUser  Event = "USER"
	EventUsers Event = "USERS"
)

type User struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Username string    `json:"username,omitempty"`

	GoogleID      *string `json:"google_id,omitempty"`
	GooglePicture *string `json:"google_picture,omitempty"`

	PDFs        []*PDF  `json:"pdfs,omitempty"`
	HostedGames []*Game `json:"hosted_games,omitempty"`
	JoinedGames []*Game `json:"joined_games,omitempty"`
}

// =====

const EventGetUsersForGameRequest Event = "GET_USERS_FOR_GAME_REQUEST"

type GetUsersForGameRequest struct {
	GameID uuid.UUID `json:"game_id,omitempty"`
}

// =====

const EventAuthenticateUserRequest Event = "AUTHENTICATE_USER_REQUEST"

type AuthenticateRequest struct {
	SessionID      uuid.UUID   `json:"session_id,omitempty"`
	SessionView    SessionView `json:"session_view,omitempty"`
	CheckCSRFToken bool        `json:"check_csrf_token,omitempty"`
	CSRFToken      string      `json:"csrf_token,omitempty"`
}

// =====

type SessionView string

const (
	SessionViewPage SessionView = "page"
	CtxKeySession   CtxKey      = "session"
	EventSession    Event       = "SESSION"
)

type Session struct {
	ID uuid.UUID `json:"id,omitempty"`

	UserID uuid.UUID `json:"user_id,omitempty"`
	User   *User     `json:"user,omitempty"`

	CSRFToken string `json:"csrf_token,omitempty"`
}

// =====

type GameView string

const (
	GameViewListItem GameView = "list_item"

	EventGame  Event = "GAME"
	EventGames Event = "GAMES"
)

type Game struct {
	ID uuid.UUID `json:"id,omitempty"`

	HostID uuid.UUID `json:"host_id,omitempty"`
	Host   *User     `json:"host,omitempty"`

	Name string `json:"name,omitempty"`

	Users []*User `json:"users,omitempty"`
	PDFs  []*PDF  `json:"pdfs,omitempty"`
	Rolls []*Roll `json:"rolls,omitempty"`
}

// =====

type PDFView string

const (
	PDFViewListItem PDFView = "list_item"

	EventPDF  Event = "PDF"
	EventPDFs Event = "PDFS"
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

// =====

const EventSubToPDFRequest Event = "SUB_TO_PDF_REQUEST"

type SubToPDFRequest struct {
	PDFID   uuid.UUID `json:"pdf_id,omitempty"`
	PageNum int       `json:"page_num,string,omitempty"`
}

// =====

const EventPDFPage Event = "PDF_PAGE"

type PDFPage struct {
	PDFID   uuid.UUID         `json:"pdf_id,omitempty"`
	PageNum int               `json:"page_num,string,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// =====

const EventUpdatePDFPageRequest Event = "UPDATE_PDF_PAGE_REQUEST"

type UpdatePDFPageRequest struct {
	PDFID      uuid.UUID `json:"pdf_id,omitempty"`
	PageNum    int       `json:"page_num,string,omitempty"`
	FieldName  string    `json:"field_name,omitempty"`
	FieldValue string    `json:"field_value,omitempty"`
}

// =====

type RollView string

const (
	RollViewListItem RollView = "list_item"

	EventRoll Event = "ROLL"
)

type Roll struct {
	ID uuid.UUID `json:"id,omitempty"`

	OwnerID uuid.UUID `json:"owner_id,omitempty"`
	Owner   *User     `json:"owner,omitempty"`

	GameID uuid.UUID `json:"game_id,omitempty"`
	Game   *Game     `json:"game,omitempty"`

	DiceTypes   []int32 `json:"dice_types,omitempty"`
	DiceResults []int32 `json:"dice_results,omitempty"`
	Modifiers   string  `json:"modifiers,omitempty"`
}

// =====

const EventCreateRollRequest Event = "CREATE_ROLL_REQUEST"

type CreateRollRequest struct {
	DiceTypes []int  `json:"dice_types,omitempty"`
	Modifiers string `json:"modifiers,omitempty"`
}
