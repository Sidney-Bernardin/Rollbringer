package api

import (
	"net/http"
	"rollbringer/src/api/views"
	account_models "rollbringer/src/services/accounts/models"
	play_models "rollbringer/src/services/play/models"

	"github.com/pkg/errors"
)

func (svr *server) handleBoardCreate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var session, _ = r.Context().Value("session").(*account_models.Session)

		board, err := play_models.NewBoard(session.UserID, r.Header.Get("HX-Prompt"))
		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot create board"))
			return
		}

		if err = svr.playDatabase.CreateBoard(r.Context(), board); err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot create board"))
			return
		}

		svr.respond(w, r, http.StatusOK, views.BoardCard(board, []*account_models.User{session.User}))
	})
}
