package internal

import (
	"fmt"

	"github.com/google/uuid"
)

type Event string

type EventWrapper[T any] struct {
	Event   Event
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
	return fmt.Sprintf("%s: %s", pd.Type, pd.Detail)
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
	UserID uuid.UUID `json:"user_id,omitempty"`
	View   UserView  `json:"view,omitempty"`
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

	EventSession Event = "SESSION"
)

type Session struct {
	ID uuid.UUID `json:"id,omitempty"`

	UserID    uuid.UUID `json:"user_id,omitempty"`
	CSRFToken string    `json:"csrf_token,omitempty"`
}

// =====

const EventGetSessionRequest Event = "GET_SESSION_REQUEST"

type GetSessionRequest struct {
	SessionID   uuid.UUID   `json:"session_id,omitempty"`
	SessionView SessionView `json:"view,omitempty"`
}

// =====

type GameView string

const (
	GameViewGameAll  GameView = "all"
	GameViewHostInfo GameView = "info"

	EventGame Event = "GAME"
)

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

const EventGetGameRequest Event = "GET_GAME_REQUEST"

type GetGameRequest struct {
	GameID uuid.UUID `json:"game_id,omitempty"`
	View   GameView  `json:"view"`
}

// =====

type PDFView string

const (
	PDFViewPDFAll    PDFView = "all"
	PDFViewOwnerInfo PDFView = "info"
	PDFViewGameInfo  PDFView = "game"

	EventPDF Event = "PDF"
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

const EventPDFFields Event = "PDF_FIELDS"

type PDFFields struct {
	PDFID   uuid.UUID         `json:"pdf_id,omitempty"`
	PageNum int               `json:"page_num,string,omitempty"`
	Fields  map[string]string `json:"fields,omitempty"`
}

// =====

const EventUpdatePDFFieldRequest Event = "UPDATE_PDF_FIELD_REQUEST"

type UpdatePDFFieldRequest struct {
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

	DiceNames   []int32 `json:"dice_names,omitempty"`
	DiceResults []int32 `json:"dice_results,omitempty"`
}

// =====

const EventCreateRollRequest Event = "CREATE_ROLL_REQUEST"

type CreateRollRequest struct {
	Dice     string `json:"dice,omitempty"`
	Modifier string `json:"modifier,omitempty"`
}
