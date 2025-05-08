package api

import (
	"maps"
	"net/http"
	"slices"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"rollbringer/src/api/views/pages"
	"rollbringer/src/domain"
	"rollbringer/src/domain/services/accounts"
)

func (svr *server) handlePageHome() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			ctx        = r.Context()
			session, _ = ctx.Value("session").(*accounts.Session)
			page       = &pages.HomeData{Session: session}
		)

		if session != nil {
			var err error

			// Get the user's rooms.
			page.Rooms, err = svr.playDatabase.GetRoomsByUserID(ctx, session.UserID)
			if err != nil {
				svr.err(w, r, errors.Wrap(err, "cannot get rooms by user-ID"))
				return
			}

			// Get the user-IDs for each room.
			roomsUserIDs := make([]uuid.UUID, 0, len(page.Rooms))
			for _, room := range page.Rooms {
				roomsUserIDs = append(roomsUserIDs, slices.Collect(maps.Keys(room.UserPermisions))...)
			}

			// Get the users for each room.
			roomsUsers, err := svr.accountsDatabase.GetUsersByUserIDs(ctx, roomsUserIDs...)
			if err != nil {
				svr.err(w, r, errors.Wrap(err, "cannot get users by room-IDs"))
				return
			}

			page.RoomsUsers = map[uuid.UUID][]accounts.User{}
			for _, room := range page.Rooms {
				for _, user := range roomsUsers {
					if _, ok := room.UserPermisions[user.ID]; ok {
						page.RoomsUsers[room.ID] = append(page.RoomsUsers[room.ID], *user)
					}
				}
			}
		}

		svr.respond(w, r, http.StatusOK, pages.Home(page))
	})
}

func (svr *server) handlePagePlay() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			ctx        = r.Context()
			session, _ = ctx.Value("session").(*accounts.Session)
			page       = &pages.PlayData{Session: session}
		)

		// Parse the room-ID.
		roomID, err := uuid.Parse(r.URL.Query().Get("r"))
		if err != nil {
			svr.err(w, r, &domain.ExternalError{Type: domain.ExternalErrorTypeInvalidUUID, Msg: err.Error()})
			return
		}

		// Join the room.
		page.Room, err = svr.play.JoinRoom(ctx, roomID, &domain.EventRoomJoinedNewcomer{
			UserID:         session.User.ID.String(),
			Username:       string(session.User.Username),
			ProfilePicture: session.User.ProfilePicture,
		})

		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot join room"))
			return
		}

		// Get the room's users.
		page.RoomUsers, err = svr.accountsDatabase.GetUsersByUserIDs(ctx, slices.Collect(maps.Keys(page.Room.UserPermisions))...)
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

		// Get the user-IDs for each board.
		boardsUserIDs := make([]uuid.UUID, 0, len(page.Boards))
		for _, board := range page.Boards {
			boardsUserIDs = append(boardsUserIDs, slices.Collect(maps.Keys(board.UserPermisions))...)
		}

		// Get the users for each board.
		boardsUsers, err := svr.accountsDatabase.GetUsersByUserIDs(ctx, boardsUserIDs...)
		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot get users by board-IDs"))
			return
		}

		page.BoardsUsers = map[uuid.UUID][]*accounts.User{}
		for _, board := range page.Boards {
			for _, user := range boardsUsers {
				if _, ok := board.UserPermisions[user.ID]; ok {
					page.BoardsUsers[board.ID] = append(page.BoardsUsers[board.ID], user)
				}
			}
		}

		svr.respond(w, r, http.StatusOK, pages.Play(page))
	})
}
