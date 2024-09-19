package users

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"rollbringer/internal"
	"rollbringer/internal/config"
	"rollbringer/internal/repositories/oauth"
	"rollbringer/internal/services"
)

type Service interface {
	services.BaseServicer

	StartLogin() (consentURL, state, codeVerifier string)
	FinishLogin(ctx context.Context, stateA, stateB, code, codeVerifier string) (*internal.Session, error)
}

type service struct {
	*services.BaseService

	schema internal.UsersSchema
	oauth  internal.OAuth
}

func NewService(
	cfg *config.Config,
	logger *slog.Logger,
	ps internal.PubSub,
	schema internal.UsersSchema,
	oauth internal.OAuth,
) Service {
	svc := &service{
		BaseService: &services.BaseService{
			Config: cfg,
			Logger: logger,
			PubSub: ps,
		},
		schema: schema,
		oauth:  oauth,
	}

	return svc
}

func (svc *service) Shutdown() error {
	svc.PubSub.Close()
	err := svc.schema.Close()
	return errors.Wrap(err, "cannot close database")
}

func (svc *service) StartLogin() (consentURL, state, codeVerifier string) {
	state = oauth.NewOauthState()
	codeVerifier = svc.oauth.GenerateCodeVerifier()
	return svc.oauth.GetConsentURL(state, codeVerifier), state, codeVerifier
}

func (svc *service) FinishLogin(ctx context.Context, stateA, stateB, code, codeVerifier string) (*internal.Session, error) {
	userInfo, err := svc.oauth.AuthenticateConsent(ctx, stateA, stateB, code, codeVerifier)
	if err != nil {
		return nil, errors.Wrap(err, "cannot authenticate consent")
	}

	user := &internal.User{
		GoogleID: &userInfo.GoogleID,
		Username: userInfo.GivenName,
	}

	if err := svc.schema.UserInsert(ctx, user); err != nil {
		return nil, errors.Wrap(err, "cannot insert user")
	}

	session := &internal.Session{
		UserID:    user.ID,
		CSRFToken: uuid.New().String(),
	}

	if err := svc.schema.SessionUpsert(ctx, session); err != nil {
		return nil, errors.Wrap(err, "cannot insert session")
	}

	return session, nil
}

func (svc *service) getUser(ctx context.Context, userID uuid.UUID, viewQuery string) (*internal.User, error) {
	views, err := internal.ParseViewQuery[internal.UserView](ctx, viewQuery)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse user view query")
	}

	user, err := svc.schema.UserGet(ctx, userID, views)
	if err != nil {
		return nil, errors.Wrap(err, "cannot get user")
	}

	errs, errsCtx := errgroup.WithContext(ctx)

	if pdfsView, ok := views["pdfs"]; ok {
		errs.Go(func() error {
			err := svc.PubSub.Request(errsCtx, "pdfs", &user.PDFs, &internal.EventWrapper[any]{
				Event: internal.EventGetPDFsByOwnerRequest,
				Payload: &internal.GetPDFsByOwnerRequest{
					OwnerID:   userID,
					ViewQuery: fmt.Sprintf("pdfs-%s", pdfsView),
				},
			})
			return errors.Wrap(err, "cannot get PDFs by owner")
		})
	}

	if gamesView, ok := views["games"]; ok {
		errs.Go(func() error {
			err := svc.PubSub.Request(errsCtx, "games", &user.HostedGames, &internal.EventWrapper[any]{
				Event: internal.EventGetGamesByHostRequest,
				Payload: &internal.GetGamesByHostRequest{
					HostID:    userID,
					ViewQuery: fmt.Sprintf("games-%s", gamesView),
				},
			})
			return errors.Wrap(err, "cannot get games by host")
		})

		errs.Go(func() error {
			err := svc.PubSub.Request(errsCtx, "games", &user.JoinedGames, &internal.EventWrapper[any]{
				Event: internal.EventGetGamesByGuestRequest,
				Payload: &internal.GetGamesByGuestRequest{
					GuestID:   userID,
					ViewQuery: fmt.Sprintf("games-%s", gamesView),
				},
			})
			return errors.Wrap(err, "cannot get games by guest")
		})
	}

	if err := errs.Wait(); err != nil {
		return nil, err
	}

	return user, nil
}

func (svc *service) getSession(ctx context.Context, sessionID uuid.UUID, viewQuery string) (*internal.Session, error) {
	views, err := internal.ParseViewQuery[internal.SessionView](ctx, viewQuery)
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse user view query")
	}

	session, err := svc.schema.SessionGet(ctx, sessionID, views)
	return session, errors.Wrap(err, "cannot get session")
}

func (svc *service) authenticate(ctx context.Context, sessionID uuid.UUID, csrfToken string) (*internal.Session, error) {
	session, err := svc.schema.SessionGet(ctx, sessionID, map[string]internal.SessionView{"session": "all"})
	if err != nil {
		if internal.IsDetailed(err, internal.PDTypeSessionNotFound) {
			return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
				Type: internal.PDTypeUnauthorized,
			})
		}

		return nil, errors.Wrap(err, "cannot get session")
	}

	if session.CSRFToken != csrfToken {
		return nil, internal.NewProblemDetail(ctx, internal.PDOpts{
			Type: internal.PDTypeUnauthorized,
		})
	}

	return session, nil
}
