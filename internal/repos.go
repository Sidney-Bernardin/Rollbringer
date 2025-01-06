package internal

import (
	"context"

	"github.com/google/uuid"
)

type PubSub interface {
	Close()
	Publish(ctx context.Context, subject string, data *EventWrapper[any]) error
	Request(ctx context.Context, subject string, res any, req *EventWrapper[any]) error
	Subscribe(ctx context.Context, subject string, cb func(*EventWrapper[[]byte]) (*EventWrapper[any], error)) error
}

type Database interface {
	Close() error
}

type UsersSchema interface {
	Database

	UserInsert(ctx context.Context, user *User) error
	UsersGetByGame(ctx context.Context, gameID uuid.UUID) ([]*User, error)

	SessionUpsert(ctx context.Context, session *Session) error
	SessionGet(ctx context.Context, sessionID uuid.UUID, view SessionView) (*Session, error)
}

type GamesSchema interface {
	Database

	GameInsert(ctx context.Context, game *Game) error
	GamesCount(ctx context.Context, hostID uuid.UUID) (int, error)
	GameGet(ctx context.Context, gameID uuid.UUID, view GameView) (*Game, error)
	GamesGetByHost(ctx context.Context, hostID uuid.UUID, view GameView) ([]*Game, error)
	GamesGetByUser(ctx context.Context, userID uuid.UUID, view GameView) ([]*Game, error)
	GameDelete(ctx context.Context, gameID, hostID uuid.UUID) error

	PDFInsert(ctx context.Context, pdf *PDF) error
	PDFGet(ctx context.Context, pdfID uuid.UUID, view PDFView) (*PDF, error)
	PDFGetPage(ctx context.Context, pdfID uuid.UUID, pageNum int) (map[string]string, error)
	PDFsGetByOwner(ctx context.Context, ownerID uuid.UUID, view PDFView) ([]*PDF, error)
	PDFsGetByGame(ctx context.Context, gameID uuid.UUID, view PDFView) ([]*PDF, error)
	PDFUpdate(ctx context.Context, session *Session, pdf *PDF) error
	PDFUpdatePage(ctx context.Context, pdfID uuid.UUID, pageNum int, fieldName, fieldValue string) error
	PDFDelete(ctx context.Context, pdfID, ownerID uuid.UUID) error

	RollInsert(ctx context.Context, roll *Roll) error
	RollsGetByGame(ctx context.Context, gameID uuid.UUID) ([]*Roll, error)

	ChatMessageInsert(ctx context.Context, chatMsg *ChatMessage) error
	ChatMessagesGetByGame(ctx context.Context, gameID uuid.UUID) ([]*ChatMessage, error)
}

type OAuth interface {
	GenerateCodeVerifier() string
	GetConsentURL(state, codeVerifier string) string
	AuthenticateConsent(ctx context.Context, stateA, stateB, code, codeVerifier string) (*GoogleUserInfo, error)
}
