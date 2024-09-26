package games

import (
	"context"
	"fmt"
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
	SubToGame(ctx context.Context, gameID uuid.UUID, resChan chan<- any) error
	DeleteGame(ctx context.Context, session *internal.Session, gameID uuid.UUID) error

	CreatePDF(ctx context.Context, session *internal.Session, pdf *internal.PDF) error
	SubToPDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int, resChan chan<- any) error
	GetPDF(ctx context.Context, pdfID uuid.UUID, viewQuery string) (*internal.PDF, error)
	GetPDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int) (map[string]string, error)
	UpdatePDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int, fieldName, fieldValue string) error
	DeletePDF(ctx context.Context, session *internal.Session, pdfID uuid.UUID) error

	CreateRoll(ctx context.Context, session *internal.Session, gameID uuid.UUID, dice []int, modifiers string) error
}

type service struct {
	*services.BaseService

	schema internal.GamesSchema

	random *rand.Rand
}

func NewService(
	cfg *config.Config,
	logger *slog.Logger,
	ps internal.PubSub,
	db internal.GamesSchema,
) Service {
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

func (svc *service) CreateGame(ctx context.Context, session *internal.Session, game *internal.Game) error {
	count, err := svc.schema.GamesCount(ctx, session.UserID)
	if err != nil {
		return errors.Wrap(err, "cannot count games")
	}

	if count >= 5 {
		return internal.NewProblemDetail(ctx, internal.PDOpts{
			Type:   internal.PDTypeMaxGames,
			Detail: "You cannot host more than 5 games at a time.",
		})
	}

	game.HostID = session.UserID
	err = svc.schema.GameInsert(ctx, game)
	return errors.Wrap(err, "cannot insert game")
}

func (svc *service) getGame(ctx context.Context, gameID uuid.UUID, viewQuery string) (*internal.Game, error) {
	views, err := internal.ParseViewQuery[internal.GameView](ctx, viewQuery)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse game view query")
	}

	game, err := svc.schema.GameGet(ctx, gameID, views)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get game")
	}

	errs, errsCtx := errgroup.WithContext(ctx)

	if usersView, ok := views["users"]; ok {
		errs.Go(func() error {
			err := svc.PubSub.Request(errsCtx, "users", &game.Users, &internal.EventWrapper[any]{
				Event: internal.EventGetUsersByGameRequest,
				Payload: internal.GetUsersByGameRequest{
					GameID:    gameID,
					ViewQuery: fmt.Sprintf("users-%s", usersView),
				},
			})
			return errors.Wrap(err, "cannot get users by game")
		})
	}

	if pdfsView, ok := views["pdfs"]; ok {
		errs.Go(func() (err error) {
			game.PDFs, err = svc.schema.PDFsGetByGame(errsCtx, gameID, map[string]internal.PDFView{"pdfs": internal.PDFView(pdfsView)})
			return errors.Wrap(err, "cannot get PDFs by game")
		})
	}

	if _, ok := views["rolls"]; ok {
		errs.Go(func() (err error) {
			game.Rolls, err = svc.schema.RollsGetByGame(errsCtx, gameID)
			return errors.Wrap(err, "cannot get rolls by game")
		})
	}

	if err := errs.Wait(); err != nil {
		return nil, err
	}

	return game, nil
}

func (svc *service) getGamesByHost(ctx context.Context, hostID uuid.UUID, viewQuery string) ([]*internal.Game, error) {
	views, err := internal.ParseViewQuery[internal.GameView](ctx, viewQuery)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse game view query")
	}

	games, err := svc.schema.GamesGetByHost(ctx, hostID, views)
	return games, errors.Wrap(err, "cannot get games by host")
}

func (svc *service) getGamesByUser(ctx context.Context, userID uuid.UUID, viewQuery string) ([]*internal.Game, error) {
	views, err := internal.ParseViewQuery[internal.GameView](ctx, viewQuery)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse game view query")
	}

	games, err := svc.schema.GamesGetByUser(ctx, userID, views)
	return games, errors.Wrap(err, "cannot get games by user")
}

func (svc *service) DeleteGame(ctx context.Context, session *internal.Session, gameID uuid.UUID) error {
	err := svc.schema.GameDelete(ctx, gameID, session.UserID)
	return errors.Wrap(err, "cannot delete game")
}

func (svc *service) CreatePDF(ctx context.Context, session *internal.Session, pdf *internal.PDF) error {
	pdf.OwnerID = session.UserID
	pdf.Pages = make([]map[string]string, len(internal.PDFSchemaPageNames[pdf.Schema]))

	err := svc.schema.PDFInsert(ctx, pdf)
	return errors.Wrap(err, "cannot insert PDF")
}

func (svc *service) GetPDF(ctx context.Context, pdfID uuid.UUID, viewQuery string) (*internal.PDF, error) {
	views, err := internal.ParseViewQuery[internal.PDFView](ctx, viewQuery)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse PDF view query")
	}

	pdf, err := svc.schema.PDFGet(ctx, pdfID, views)
	return pdf, errors.Wrap(err, "cannot get PDF")
}

func (svc *service) getPDFsByOwner(ctx context.Context, ownerID uuid.UUID, viewQuery string) ([]*internal.PDF, error) {
	views, err := internal.ParseViewQuery[internal.PDFView](ctx, viewQuery)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse PDF view query")
	}

	pdfs, err := svc.schema.PDFsGetByOwner(ctx, ownerID, views)
	return pdfs, errors.Wrap(err, "cannot get PDF by owner")
}

func (svc *service) GetPDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int) (map[string]string, error) {
	pageFields, err := svc.schema.PDFGetPage(ctx, pdfID, pageNum)
	return pageFields, errors.Wrap(err, "cannot get PDF page")
}

func (svc *service) UpdatePDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int, fieldName, fieldValue string) error {
	err := svc.schema.PDFUpdatePage(ctx, pdfID, pageNum, fieldName, fieldValue)
	return errors.Wrap(err, "cannot update PDF page")
}

func (svc *service) DeletePDF(ctx context.Context, session *internal.Session, pdfID uuid.UUID) error {
	err := svc.schema.PDFDelete(ctx, pdfID, session.UserID)
	return errors.Wrap(err, "cannot delete PDF")
}

func (svc *service) CreateRoll(ctx context.Context, session *internal.Session, gameID uuid.UUID, diceTypes []int, modifiers string) error {
	roll := &internal.Roll{
		OwnerID:     session.UserID,
		GameID:      gameID,
		DiceTypes:   []int32{},
		DiceResults: []int32{},
		Modifiers:   modifiers,
	}

	for _, dType := range diceTypes {
		dType32 := int32(dType)
		roll.DiceTypes = append(roll.DiceTypes, dType32)
		roll.DiceResults = append(roll.DiceResults, svc.random.Int31n(dType32)+1)
	}

	err := svc.schema.RollInsert(ctx, roll)
	return errors.Wrap(err, "cannot insert roll")
}
