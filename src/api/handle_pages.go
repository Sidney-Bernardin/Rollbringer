package api

import (
	"io"
	"net"
	"net/http"

	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/src"
	"rollbringer/src/api/views"
	"rollbringer/src/api/views/pages"
	"rollbringer/src/services"
	account_models "rollbringer/src/services/accounts/models"
	play_models "rollbringer/src/services/play/models"
)

func (svr *server) handlePageHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			ctx        = r.Context()
			session, _ = ctx.Value("session").(*account_models.Session)

			err  error
			page = &pages.HomeData{
				Session:   session,
				Rooms:     []*play_models.Room{},
				RoomUsers: map[src.UUID][]*account_models.User{},
			}
		)

		if session == nil {
			svr.respond(w, r, http.StatusOK, pages.Home(page))
			return
		}

		// Get the user's rooms.
		page.Rooms, err = svr.playDatabase.GetRoomsByUserID(ctx, session.UserID)
		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot get rooms by user-ID"))
			return
		}

		roomIDs := make([]src.UUID, 0, len(page.Rooms))
		for _, room := range page.Rooms {
			roomIDs = append(roomIDs, room.ID)
		}

		// Get users for each room.
		page.RoomUsers, err = svr.accountsDatabase.GetUsersByRoomIDs(ctx, roomIDs...)
		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot get users by room-IDs"))
			return
		}

		svr.respond(w, r, http.StatusOK, pages.Home(page))
	})
}

func (svr *server) handlePagePlay() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			ctx        = r.Context()
			session, _ = ctx.Value("session").(*account_models.Session)
			page       = &pages.PlayData{Session: session}
		)

		roomID, err := src.ParseUUID(r.URL.Query().Get("r"))
		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot parse room-ID"))
			return
		}

		// Join the room.
		page.Room, err = svr.play.JoinRoom(ctx, session.UserID, roomID)
		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot get room by room-ID"))
			return
		}

		// Get the room's users.
		page.RoomUsers, err = svr.accountsDatabase.GetUsersByRoomID(ctx, roomID)
		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot get users by room-ID"))
			return
		}

		// Get the user's boards.
		page.Boards, err = svr.playDatabase.GetBoardsByUserID(ctx, session.UserID)
		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot get boards by user-ID"))
			return
		}

		boardIDs := make([]src.UUID, 0, len(page.Boards))
		for _, board := range page.Boards {
			boardIDs = append(boardIDs, board.ID)
		}

		// Get users for each board.
		page.BoardUsers, err = svr.accountsDatabase.GetUsersByBoardIDs(ctx, boardIDs...)
		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot get users by board-IDs"))
			return
		}

		svr.respond(w, r, http.StatusOK, pages.Play(page))
	})
}

func (svr *server) handlePagePlayWebSocket() websocket.Handler {
	var events = map[string]any{
		"chat":         &services.EventChat{},
		"create_board": &views.CreateBoardRequest{},
	}

	return websocket.Handler(func(conn *websocket.Conn) {

		var (
			r          = conn.Request()
			ctx        = r.Context()
			session, _ = ctx.Value("session").(*account_models.Session)
		)

		roomID, err := src.ParseUUID(r.URL.Query().Get("r"))
		if err != nil {
			svr.err(conn, r, errors.Wrap(err, "cannot parse room-ID"))
			return
		}

		go func() {
			err := svr.broker.SubRoom(ctx, roomID, func(event any) {
				switch e := event.(type) {
				case *services.EventChat:
					svr.respond(conn, r, 0, views.ChatMessage(e))
				}
			})

			svr.err(conn, r, errors.Wrap(err, "cannot subscribe to room"))
			conn.Close()
		}()

		for {

			var msg []byte
			if err := websocket.Message.Receive(conn, &msg); err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
					return
				}

				svr.err(conn, r, errors.Wrap(err, "cannot read from WebSocket connection"))
				return
			}

			switch event := decodeEvent(msg, events).(type) {
			case *services.EventChat:
				event.RoomID = roomID
				event.Username = session.User.Username
				event.ProfilePicture = session.User.ProfilePicture

				err = svr.broker.PubChat(event)
				svr.err(conn, r, errors.Wrap(err, "cannot publish chat event"))

			case *views.CreateBoardRequest:
				board, err := play_models.NewBoard(session.UserID, event.Name)
				if err != nil {
					svr.err(conn, r, errors.Wrap(err, "cannot create board"))
					return
				}

				if err = svr.playDatabase.CreateBoard(r.Context(), board); err != nil {
					svr.err(conn, r, errors.Wrap(err, "cannot create board"))
					return
				}
			}
		}
	})
}
