package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"rollbringer/pkg/domain"
	"rollbringer/pkg/handlers/play/views"
)

var boardViews = map[string]domain.BoardView{
	"all":       domain.BoardViewAll,
	"list-item": domain.BoardViewListItem,
}

func (h *playHandler) handleBoardsPost(w http.ResponseWriter, r *http.Request) {

	var (
		state = h.State(r)
		ctx   = r.Context()

		roomID, _ = uuid.Parse(r.FormValue("room_id"))
		view, _   = boardViews[r.URL.Query().Get("view")]

		board = &domain.Board{
			RoomID: roomID,
			Name:   r.FormValue("name"),
			Konva:  []byte(r.FormValue("konva")),
		}
	)

	if err := h.playSvc.CreateBoard(ctx, state["session"].(*domain.Session), view, board); err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot create board", nil))
		return
	}

	h.Respond(w, r, http.StatusOK, views.BoardView(board, view))
}

func (h *playHandler) handleBoardGet(w http.ResponseWriter, r *http.Request) {

	var (
		ctx = r.Context()

		boardID, _ = uuid.Parse(chi.URLParam(r, "board_id"))
		view, _    = boardViews[r.URL.Query().Get("view")]
	)

	board, err := h.playSvc.GetBoard(ctx, view, boardID)
	if err != nil {
		h.Err(w, r, domain.Wrap(err, "cannot create board", nil))
		return
	}

	h.Respond(w, r, http.StatusOK, views.BoardView(board, view))
}
