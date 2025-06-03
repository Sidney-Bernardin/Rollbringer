package http

import (
	"net/http"

	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/nats"
	"github.com/Sidney-Bernardin/Rollbringer/web/pages"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func (api *API) handleHomePage(w http.ResponseWriter, r *http.Request) {

	var (
		ctx  = r.Context()
		data pages.HomeData
	)

	data.Session, _ = ctx.Value("session").(*nats.Session)
	if data.Session == nil {
		api.respond(w, r, http.StatusOK, pages.HomePage(&data))
		return
	}

	errs, errsCtx := errgroup.WithContext(ctx)

	errs.Go(func() (err error) {
		data.User, err = api.Service.GetUser(errsCtx, data.Session.UserID)
		return errors.Wrap(err, "cannot get user")
	})

	if err := errs.Wait(); err != nil {
		api.err(w, r, errors.WithStack(err))
		return
	}

	api.respond(w, r, http.StatusOK, pages.HomePage(&data))
}
