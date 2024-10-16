package games

import (
	"context"
	"log/slog"
	"math/rand"
	"time"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/services"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

type Service interface {
	services.BaseServicer

	CreateGame(ctx context.Context, session *internal.Session, game *internal.Game) error
	DeleteGame(ctx context.Context, session *internal.Session, gameID uuid.UUID) error
	SubToGame(ctx context.Context, gameID uuid.UUID, resChan chan<- any) error

	CreatePDF(ctx context.Context, session *internal.Session, pdf *internal.PDF) error
	GetPDF(ctx context.Context, pdfID uuid.UUID, view internal.PDFView) (*internal.PDF, error)
	GetPDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int) (map[string]string, error)
	UpdatePDF(ctx context.Context, session *internal.Session, pdf *internal.PDF) error
	UpdatePDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int, fieldName, fieldValue string) error
	DeletePDF(ctx context.Context, session *internal.Session, pdfID uuid.UUID) error
	SubToPDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int, resChan chan<- any) error

	CreateRoll(ctx context.Context, session *internal.Session, gameID uuid.UUID, dice []int, modifiers string) (*internal.Roll, error)
}

type service struct {
	*services.BaseService

	schema internal.GamesSchema

	random *rand.Rand
}

func NewService(cfg *config.Config, logger *slog.Logger, ps internal.PubSub, db internal.GamesSchema) Service {
	return &service{
		BaseService: &services.BaseService{
			Config: cfg,
			Logger: logger,
			PubSub: ps,
		},
		schema: db,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (svc *service) Shutdown() error {
	svc.PubSub.Close()
	err := svc.schema.Close()
	return errors.Wrap(err, "cannot close database")
}

func (svc *service) playPage(ctx context.Context, page *internal.PlayPage) (err error) {
	errs, errsCtx := errgroup.WithContext(ctx)

	errs.Go(func() error {
		page.Session.User.PDFs, err = svc.schema.PDFsGetByOwner(errsCtx, page.Session.UserID, internal.PDFViewListItem)
		return errors.Wrap(err, "cannot get PDFs by owner")
	})

	errs.Go(func() error {
		page.Session.User.HostedGames, err = svc.schema.GamesGetByHost(errsCtx, page.Session.UserID, internal.GameViewListItem)
		return errors.Wrap(err, "cannot get games by host")
	})

	errs.Go(func() error {
		page.Session.User.JoinedGames, err = svc.schema.GamesGetByUser(errsCtx, page.Session.UserID, internal.GameViewListItem)
		return errors.Wrap(err, "cannot get games by user")
	})

	if page.Game != nil {
		if page.Game, err = svc.schema.GameGet(errsCtx, page.Game.ID, ""); err != nil {
			return errors.Wrap(err, "cannot get game")
		}

		errs.Go(func() error {
			page.Game.PDFs, err = svc.schema.PDFsGetByGame(errsCtx, page.Game.ID, internal.PDFViewListItem)
			return errors.Wrap(err, "cannot get PDFs by game")
		})

		errs.Go(func() error {
			page.Game.Rolls, err = svc.schema.RollsGetByGame(errsCtx, page.Game.ID)
			return errors.Wrap(err, "cannot get rolls by game")
		})
	}

	return errs.Wait()
}
