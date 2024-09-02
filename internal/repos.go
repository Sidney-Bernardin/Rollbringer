package internal

import (
	"context"

	"github.com/google/uuid"
)

type PubSub interface {
	Publish(ctx context.Context, subject string, event Event) error
	Request(ctx context.Context, subject string, req Event) (Event, error)
	Subscribe(ctx context.Context, subject string, errChan chan<- error, cb func(event Event, subject []string) (Event, *ProblemDetail))
	ChanSubscribe(ctx context.Context, subject string, resChan chan<- Event, errChan chan<- error)
}

type UsersDatabase interface {
	UserInsert(ctx context.Context, user *User) error
	UserGet(ctx context.Context, userID uuid.UUID, view UserView) (*User, error)

	SessionUpsert(ctx context.Context, session *Session) error
	SessionGet(ctx context.Context, sessionID uuid.UUID, view SessionView) (*Session, error)
}

type GamesDatabase interface {
	GameInsert(ctx context.Context, game *Game) error
	GamesCount(ctx context.Context, hostID uuid.UUID) (int, error)
	GamesGetForHost(ctx context.Context, hostID uuid.UUID, view GameView) ([]*Game, error)
	GameGet(ctx context.Context, gameID uuid.UUID, view GameView) (*Game, error)
	GameDelete(ctx context.Context, gameID, hostID uuid.UUID) error

	PDFInsert(ctx context.Context, pdf *PDF, pageCount int) error
	PDFsGetForOwner(ctx context.Context, ownerID uuid.UUID, view PDFView) ([]*PDF, error)
	PDFsGetForGame(ctx context.Context, gameID uuid.UUID, view PDFView) ([]*PDF, error)
	PDFGet(ctx context.Context, pdfID uuid.UUID, view PDFView) (*PDF, error)
	PDFGetPage(ctx context.Context, pdfID uuid.UUID, pageIdx int) (map[string]string, error)
	PDFUpdatePage(ctx context.Context, pdfID uuid.UUID, pageIdx int, fieldName, fieldValue string) error
	PDFDelete(ctx context.Context, pdfID, ownerID uuid.UUID) error

	RollInsert(ctx context.Context, roll *Roll) error
	RollsGetForGame(ctx context.Context, gameID uuid.UUID) ([]*Roll, error)
}

type OAuth interface {
	GenerateCodeVerifier() string
	GetConsentURL(state, codeVerifier string) string
	AuthenticateConsent(ctx context.Context, stateA, stateB, code, codeVerifier string) (*GoogleUserInfo, error)
}
