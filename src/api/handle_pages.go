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
				svr.err(w, r, errors.Wrap(err, "cannot get user's rooms"))
				return
			}

			// Get the room's user-IDs.
			roomsUserIDs := make([]uuid.UUID, 0, len(page.Rooms))
			for _, room := range page.Rooms {
				roomsUserIDs = append(roomsUserIDs, slices.Collect(maps.Keys(room.UserPermisions))...)
			}

			// Get the room's users.
			roomsUsers, err := svr.accountsDatabase.GetUsersByUserIDs(ctx, roomsUserIDs...)
			if err != nil {
				svr.err(w, r, errors.Wrap(err, "cannot get room's users"))
				return
			}

			// Group the room's users by room-ID.
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
		page.Room, err = svr.play.JoinRoom(ctx, roomID, &domain.PublicUser{
			UserID:         session.User.ID,
			Username:       string(session.User.Username),
			ProfilePicture: session.User.ProfilePicture,
		})

		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot join room"))
			return
		}

		// Get the room's users.
		roomUserIDs := slices.Collect(maps.Keys(page.Room.UserPermisions))
		page.RoomUsers, err = svr.accountsDatabase.GetUsersByUserIDs(ctx, roomUserIDs...)
		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot get room's users"))
			return
		}

		// Get the user's boards.
		page.Boards, err = svr.playDatabase.GetBoardsByUserID(ctx, session.UserID)
		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot get user's boards"))
			return
		}

		// Get the board's user-IDs.
		boardsUserIDs := make([]uuid.UUID, 0, len(page.Boards))
		for _, board := range page.Boards {
			boardsUserIDs = append(boardsUserIDs, slices.Collect(maps.Keys(board.UserPermisions))...)
		}

		// Get the board's users.
		boardsUsers, err := svr.accountsDatabase.GetUsersByUserIDs(ctx, boardsUserIDs...)
		if err != nil {
			svr.err(w, r, errors.Wrap(err, "cannot get board's users"))
			return
		}

		// Group the board's users by board-ID.
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
