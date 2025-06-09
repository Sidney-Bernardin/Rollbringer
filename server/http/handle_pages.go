package http

import (
	"net/http"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/cache"
	"github.com/Sidney-Bernardin/Rollbringer/web/pages"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

func (api *API) handleHomePage(w http.ResponseWriter, r *http.Request) {

	var (
		ctx  = r.Context()
		data pages.HomeData
	)

	data.Session, _ = ctx.Value("session").(*cache.Session)
	if data.Session == nil {
		api.respond(w, r, http.StatusOK, pages.HomePage(&data))
		return
	}

	errs, errsCtx := errgroup.WithContext(ctx)

	errs.Go(func() (err error) {
		data.User, err = api.Service.GetUser(errsCtx, data.Session.UserID)
		return errors.Wrap(err, "cannot get user")
	})

	errs.Go(func() (err error) {
		data.UserRooms, err = api.Service.GetUserRooms(errsCtx, data.Session.UserID)
		return errors.Wrap(err, "cannot get user rooms")
	})

	if err := errs.Wait(); err != nil {
		api.err(w, r, errors.WithStack(err))
		return
	}

	api.respond(w, r, http.StatusOK, pages.HomePage(&data))
}

func (api *API) handlePlayPage(w http.ResponseWriter, r *http.Request) {

	var (
		ctx  = r.Context()
		data pages.PlayData
	)

	data.Session, _ = ctx.Value("session").(*cache.Session)
	if data.Session == nil {
		api.respond(w, r, http.StatusOK, pages.PlayPage(&data))
		return
	}

	roomID, err := server.ParseUUID(r.URL.Query().Get("r"))
	if err != nil {
		api.err(w, r, errors.Wrap(err, "cannot parse room-ID"))
		return
	}

	errs, errsCtx := errgroup.WithContext(ctx)

	errs.Go(func() (err error) {
		data.User, err = api.Service.GetUser(errsCtx, data.Session.UserID)
		return errors.Wrap(err, "cannot get user")
	})

	errs.Go(func() (err error) {
		data.Room, err = api.Service.GetRoom(ctx, roomID)
		return errors.Wrap(err, "cannot get room")
	})

	if err := errs.Wait(); err != nil {
		api.err(w, r, errors.WithStack(err))
		return
	}

	api.respond(w, r, http.StatusOK, pages.PlayPage(&data))
}
