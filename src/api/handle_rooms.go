package api

import (
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
		)

		// Parse the room-ID.
		roomID, err := uuid.Parse(r.PathValue("room_id"))
		if err != nil {
			svr.err(conn, r, &domain.ExternalError{Type: domain.ExternalErrorTypeInvalidUUID, Msg: err.Error()})
			return
		}

		// Join the room.
		_, err = svr.play.JoinRoom(ctx, roomID, &domain.EventRoomJoinedNewcomer{
			UserID:         session.User.ID.String(),
			Username:       string(session.User.Username),
			ProfilePicture: session.User.ProfilePicture,
		})

		if err != nil {
			svr.err(conn, r, errors.Wrap(err, "cannot join room"))
			return
		}

		// Subscribe to the room's events.
		go func() {
			err := svr.playBroker.SubRoom(ctx, roomID, svr.roomCallback(conn, r))
			svr.err(conn, r, errors.Wrap(err, "cannot subscribe to room events"))
			conn.Close()
		}()

		// Subscribe to the room's chat events.
		go func() {
			err := svr.playBroker.SubChat(ctx, roomID, svr.chatCallback(ctx, conn, r))
			svr.err(conn, r, errors.Wrap(err, "cannot subscribe to chat events"))
			conn.Close()
		}()

		// Subscribe to the user's events.
		go func() {
			err := svr.playBroker.SubUser(ctx, roomID, svr.userCallback(conn, r))
			svr.err(conn, r, errors.Wrap(err, "cannot subscribe to user events"))
			conn.Close()
		}()

		for {

			// Receive the next message.
			event, err := svr.nextMessage(conn, func(operation string) any {
				switch operation {
				case "chat":
					return &views.ReqChat{}
				case "create-board":
					return &play.CreateBoardOpts{}
				default:
					return nil
				}
			})

			switch e := event.(type) {
			case error:
				if !errors.Is(err, io.EOF) && !errors.Is(err, net.ErrClosed) {
					svr.err(conn, r, errors.Wrap(err, "cannot read from WebSocket connection"))
				}
				return

			case *views.ReqChat:

				// Publish the chat message.
				svr.playBroker.Pub(ctx, &play.EventChat{
					RoomID:   roomID.String(),
					AuthorID: session.User.ID.String(),
					Message:  e.Message,
				})

			case *play.CreateBoardOpts:

				// Parse the user-IDs.
				userIDs := make([]uuid.UUID, 0, len(e.Users))
				for _, userID := range e.Users {
					uID, err := uuid.Parse(userID)
					if err != nil {
						svr.err(conn, r, &domain.ExternalError{Type: domain.ExternalErrorTypeInvalidUUID, Msg: err.Error()})
						continue
					}
					userIDs = append(userIDs, uID)
				}

				// Get the board's users.
				users, err := svr.accountsDatabase.GetUsersByUserIDs(ctx, userIDs...)
				if err != nil {
					svr.err(conn, r, errors.Wrap(err, "cannot get users for new board"))
					continue
				}

				// Create the board.
				_, err = svr.play.CreateBoard(ctx, e,
					&domain.EventNewBoardUser{
						UserID:         session.User.ID.String(),
						Username:       string(session.User.Username),
						ProfilePicture: session.User.ProfilePicture,
					},
					src.Map(users, func(_ int, u *accounts.User) *domain.EventNewBoardUser {
						return &domain.EventNewBoardUser{
							UserID:         u.ID.String(),
							Username:       string(u.Username),
							ProfilePicture: u.ProfilePicture,
						}
					}),
				)

				svr.err(conn, r, errors.Wrap(err, "cannot create board"))
			}
		}
	})
}
