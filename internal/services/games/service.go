package games

import (
	"context"
	"log/slog"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/services"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type Service interface {
	services.BaseServicer

	CreateGame(ctx context.Context, session *internal.Session, game *internal.Game) error
	DeleteGame(ctx context.Context, session *internal.Session, gameID uuid.UUID) error

	CreatePDF(ctx context.Context, session *internal.Session, pdf *internal.PDF) error
	SubToPDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int, outgoing chan<- any) error
	GetPDF(ctx context.Context, pdfID uuid.UUID, views string) (*internal.PDF, error)
	GetPDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int) (map[string]string, error)
	UpdatePDFPage(ctx context.Context, pdfID uuid.UUID, pageNum int, fieldName, fieldValue string) error
	DeletePDF(ctx context.Context, session *internal.Session, pdfID uuid.UUID) error
}

type service struct {
	*services.BaseService

	db internal.GamesSchema

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
			PS:     ps,
		},
		db:     db,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (svc *service) Shutdown() error {
	svc.PS.Close()
	err := svc.db.Close()
	return errors.Wrap(err, "cannot close database")
}

func (svc *service) CreateGame(ctx context.Context, session *internal.Session, game *internal.Game) error {
	count, err := svc.db.GamesCount(ctx, session.UserID)
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
	err = svc.db.GameInsert(ctx, game)
	return errors.Wrap(err, "cannot insert game")
}

func (svc *service) getGame(ctx context.Context, gameID uuid.UUID, views string) (*internal.Game, error) {
	svc.db.Transaction(ctx, func(schema *internal.UsersSchema) error {

	})

	parsedViews, err := internal.ParseViews[internal.GameView](ctx, views)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse game view")
	}

	game, err := svc.db.GameGet(ctx, gameID, parsedViews)
	return game, errors.Wrap(err, "cannot get game")
}

func (svc *service) DeleteGame(ctx context.Context, session *internal.Session, gameID uuid.UUID) error {
	err := svc.db.GameDelete(ctx, gameID, session.UserID)
	return errors.Wrap(err, "cannot delete game")
}

func (svc *service) CreatePDF(ctx context.Context, session *internal.Session, pdf *internal.PDF) error {
	pdf.OwnerID = session.UserID
	pdf.Fields = make([]map[string]string, len(internal.PDFSchemaPageNames[pdf.Schema]))

	err := svc.db.PDFInsert(ctx, pdf, len(internal.PDFSchemaPageNames[pdf.Schema]))
	return errors.Wrap(err, "cannot insert PDF")
}

func (svc *service) GetPDF(ctx context.Context, pdfID uuid.UUID, views string) (*internal.PDF, error) {
	parsedViews, err := internal.ParseViews[internal.PDFView](ctx, views)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse PDF view")
	}

	pdf, err := svc.db.PDFGet(ctx, pdfID, parsedViews)
	return pdf, errors.Wrap(err, "cannot get PDF")
}

func (svc *service) DeletePDF(ctx context.Context, session *internal.Session, pdfID uuid.UUID) error {
	err := svc.db.PDFDelete(ctx, pdfID, session.UserID)
	return errors.Wrap(err, "cannot delete PDF")
}

func (svc *service) roll(ctx context.Context, diceNamesStr string) (*internal.Roll, error) {
	roll := &internal.Roll{
		DiceNames:   []int32{},
		DiceResults: []int32{},
	}

	for _, dieNameStr := range strings.Split(diceNamesStr, "d")[1:] {
		dName, err := strconv.ParseInt(dieNameStr, 10, 32)
		if err != nil {
			return nil, &internal.ProblemDetail{
				Instance: ctx.Value(internal.CtxKeyInstance).(string),
				Type:     internal.PDTypeInvalidDie,
				Detail:   "Die names must resemble 32-bit integers.",
				Extra: map[string]any{
					"die_name": dieNameStr,
				},
			}
		}

		roll.DiceNames = append(roll.DiceNames, int32(dName))
	}

	for _, dieName := range roll.DiceNames {
		roll.DiceResults = append(roll.DiceResults, svc.random.Int31n(dieName)+1)
	}

	return roll, nil
}
