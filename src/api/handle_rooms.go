package api

import (
	"context"
	"io"
	"net"
	"net/http"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/net/websocket"

	"rollbringer/src"
	"rollbringer/src/api/views"
	"rollbringer/src/domain"
	"rollbringer/src/domain/services/accounts"
	"rollbringer/src/domain/services/play"
)

func (svr *server) handleRoomCreate() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var session, _ = r.Context().Value("session").(*accounts.Session)

		// Create the room.
		room, err := svr.play.CreateRoom(r.Context(), session.User.ID,
			&play.CreateRoomOpts{
				Name: r.FormValue("name"),
			})

		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot create room"))
			return
		}

		svr.respond(w, r, http.StatusOK, views.RoomCard(room, []accounts.User{
			{
				ID:             session.User.ID,
				Username:       session.User.Username,
				ProfilePicture: session.User.ProfilePicture,
			},
		}))
	})
}

func (svr *server) handleRoomWebSocket() websocket.Handler {
	return websocket.Handler(func(conn *websocket.Conn) {

		var (
			r          = conn.Request()
			ctx        = r.Context()
			session, _ = ctx.Value("session").(*accounts.Session)

			board       *play.Board
			boardCtx    context.Context
			boardCancel context.CancelFunc = func() {}
		)

		// Parse the room-ID.
		roomID, err := uuid.Parse(r.PathValue("room_id"))
		if err != nil {
			svr.err(conn, r, &domain.ExternalError{Type: domain.ExternalErrorTypeInvalidUUID, Msg: err.Error()})
			return
		}

		// Join the room.
		_, err = svr.play.JoinRoom(ctx, roomID, &domain.PublicUser{
			UserID:         session.User.ID,
			Username:       string(session.User.Username),
			ProfilePicture: session.User.ProfilePicture,
		})

		if err != nil {
			svr.err(conn, r, errors.Wrap(err, "cannot join room"))
			return
		}

		// Subscribe to the room's events.
		go func() {
			err := svr.broker.SubRoom(ctx, roomID, svr.roomCallback(conn, r))
			svr.err(conn, r, errors.Wrap(err, "cannot subscribe to room events"))
			conn.Close()
		}()

		// Subscribe to the room's chat events.
		go func() {
			err := svr.broker.SubChat(ctx, roomID, svr.chatCallback(ctx, conn, r))
			svr.err(conn, r, errors.Wrap(err, "cannot subscribe to chat events"))
			conn.Close()
		}()

		// Subscribe to the user's events.
		go func() {
			err := svr.broker.SubUser(ctx, session.User.ID, svr.userCallback(conn, r))
			svr.err(conn, r, errors.Wrap(err, "cannot subscribe to user events"))
			conn.Close()
		}()

		for {

			// Receive the next message.
			message, err := svr.nextMessage(conn, func(operation string) any {
				switch operation {
				case "chat":
					return &views.ReqChatMessage{}
				case "create-board":
					return &play.CreateBoardOpts{}
				case "open-board":
					return &views.ReqGetBoard{}
				case "subscribe-to-canvas":
					return &views.ReqSubscribeToCanvas{}
				case "update-canvas-node":
					return &domain.EventUpdateCanvasNode{}
				default:
					return nil
				}
			})

			switch msg := message.(type) {
			case error:
				if !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed) {
					svr.err(conn, r, errors.Wrap(err, "cannot read from WebSocket connection"))
				}
				return

			case *views.ReqChatMessage:

				// Publish a chat-message event.
				svr.broker.Pub(ctx, &domain.EventChatMessage{
					RoomID:   roomID.String(),
					AuthorID: session.User.ID.String(),
					Message:  msg.Message,
				})

			case *play.CreateBoardOpts:

				// Get the board's users.
				users, err := svr.accountsDatabase.GetUsersByUserIDs(ctx, msg.UserIDs...)
				if err != nil {
					svr.err(conn, r, errors.Wrap(err, "cannot get users for new board"))
					continue
				}

				// Create the board.
				_, err = svr.play.CreateBoard(ctx, msg,
					&domain.PublicUser{
						UserID:         session.User.ID,
						Username:       string(session.User.Username),
						ProfilePicture: session.User.ProfilePicture,
					},
					src.Map(users, func(_ int, u *accounts.User) *domain.PublicUser {
						return &domain.PublicUser{
							UserID:         u.ID,
							Username:       string(u.Username),
							ProfilePicture: u.ProfilePicture,
						}
					}),
				)

				svr.err(conn, r, errors.Wrap(err, "cannot create board"))

			case *views.ReqGetBoard:

				// Get the user's board.
				board, err = svr.playDatabase.GetUserBoard(ctx, session.User.ID, msg.BoardID)
				if err != nil {
					svr.err(conn, r, errors.Wrap(err, "cannot get user's board"))
					continue
				}

				svr.respond(conn, r, 0, views.Board(board))

			case *views.ReqSubscribeToCanvas:
				if board == nil {
					continue
				}

				// Subscribe to the board's canvas events.
				go func() {
					boardCancel()
					boardCtx, boardCancel = context.WithCancel(ctx)
					err := svr.broker.SubCanvas(boardCtx, board.ID, svr.canvasCallback(conn, r))
					svr.err(conn, r, errors.Wrap(err, "cannot subscribe to user events"))
					boardCancel()
				}()

			case *domain.EventUpdateCanvasNode:
				if board == nil {
					continue
				}

				msg.BoardID = board.ID
				svr.broker.Pub(ctx, msg)
			}
		}
	})
}
