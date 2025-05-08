package api

import (
	"encoding/json"
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

		room, err := svr.play.CreateRoom(r.Context(), session.User.ID, &play.CreateRoomOpts{
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

		go svr.subChat(conn, r, roomID)
		go svr.subRoom(conn, r, roomID)
		go svr.subUser(conn, r, session.User.ID)

		for {

			// Receive next message.
			var msg []byte
			if err := websocket.Message.Receive(conn, &msg); err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
					return
				}

				svr.err(conn, r, errors.Wrap(err, "cannot read from WebSocket connection"))
				return
			}

			var head struct {
				Operation string `json:"operation"`
			}

			// Decode the message's operation.
			if err := json.Unmarshal(msg, &head); err != nil {
				continue
			}

			// Set the message's schema based on it's operation.
			var event any
			switch head.Operation {
			case "chat":
				event = &views.ReqChat{}
			case "create-board":
				event = &play.CreateBoardOpts{}
			}

			// Decode the message.
			if err := json.Unmarshal(msg, event); err != nil {
				continue
			}

			switch e := event.(type) {
			case *views.ReqChat:

				// Publish chat event.
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

				// Get the users.
				users, err := svr.accountsDatabase.GetUsersByUserIDs(ctx, userIDs...)
				if err != nil {
					svr.err(conn, r, errors.Wrap(err, "cannot get users by user-IDs"))
					continue
				}

				// Create a board.
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
