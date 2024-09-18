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
	return fmt.Sprintf(`ProblemDetail("%s", "%s")`, pd.Type, pd.Detail)
}

// =====

type GoogleUserInfo struct {
	GoogleID  string
	GivenName string
}

// =====

type UserView string

const (
	UserViewUserAll UserView = "all"

	EventUser Event = "USER"
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

const EventGetUserRequest Event = "GET_USER_REQUEST"

type GetUserRequest struct {
	UserID    uuid.UUID `json:"user_id,omitempty"`
	ViewQuery string    `json:"view_query,omitempty"`
}

// =====

const EventAuthenticateUserRequest Event = "AUTHENTICATE_USER_REQUEST"

type AuthenticateUserRequest struct {
	SessionID uuid.UUID `json:"session_id,omitempty"`
	CSRFToken string    `json:"csrf_token,omitempty"`
}

// =====

type SessionView string

const (
	SessionViewSessionAll SessionView = "all"

	CtxKeySession CtxKey = "session"
	EventSession  Event  = "SESSION"
)

type Session struct {
	ID uuid.UUID `json:"id,omitempty"`

	UserID    uuid.UUID `json:"user_id,omitempty"`
	CSRFToken string    `json:"csrf_token,omitempty"`
}

// =====

const EventGetSessionRequest Event = "GET_SESSION_REQUEST"

type GetSessionRequest struct {
	SessionID uuid.UUID `json:"session_id,omitempty"`
	ViewQuery string    `json:"view_query,omitempty"`
}

// =====

type GameView string

const (
	GameViewGameAll  GameView = "all"
	GameViewHostInfo GameView = "info"

	EventGame  Event = "GAME"
	EventGames Event = "GAMES"
)

type Game struct {
	ID uuid.UUID `json:"id,omitempty"`

	HostID uuid.UUID `json:"host_id,omitempty"`
	Host   *User     `json:"host,omitempty"`

	Name string `json:"name,omitempty"`

	Guests []*User `json:"guests,omitempty"`
	PDFs   []*PDF  `json:"pdfs,omitempty"`
	Rolls  []*Roll `json:"rolls,omitempty"`
}

// =====

const EventGetGameRequest Event = "GET_GAME_REQUEST"

type GetGameRequest struct {
	GameID    uuid.UUID `json:"game_id,omitempty"`
	ViewQuery string    `json:"view_query,omitempty"`
}

// =====

const EventGetGamesByHostRequest Event = "GET_GAMES_BY_HOST_REQUEST"

type GetGamesByHostRequest struct {
	HostID    uuid.UUID `json:"host_id,omitempty"`
	ViewQuery string    `json:"view_query,omitempty"`
}

// =====

const EventGetGamesByGuestRequest Event = "GET_GAMES_BY_GUEST_REQUEST"

type GetGamesByGuestRequest struct {
	GuestID   uuid.UUID `json:"guest_id,omitempty"`
	ViewQuery string    `json:"view_query,omitempty"`
}

// =====

type PDFView string

const (
	PDFViewPDFAll    PDFView = "all"
	PDFViewOwnerInfo PDFView = "all"
	PDFViewGameInfo  PDFView = "all"

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

const EventGetPDFsByOwnerRequest Event = "GET_PDFS_BY_OWNER_REQUEST"

type GetPDFsByOwnerRequest struct {
	OwnerID   uuid.UUID `json:"owner_id,omitempty"`
	ViewQuery string    `json:"view_query,omitempty"`
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

const EventRoll Event = "ROLL"

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
