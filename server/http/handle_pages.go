package http

import (
	"encoding/json"
	"io"
	"net"
	"net/http"

	"github.com/Sidney-Bernardin/Rollbringer/server"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/cache"
	"github.com/Sidney-Bernardin/Rollbringer/server/repositories/pubsub"
	"github.com/Sidney-Bernardin/Rollbringer/web/pages/home"
	"github.com/Sidney-Bernardin/Rollbringer/web/pages/play"

	"github.com/pkg/errors"
	"golang.org/x/net/websocket"
	"golang.org/x/sync/errgroup"
)

func (api *API) handleHomePage(w http.ResponseWriter, r *http.Request) {

	var (
		ctx  = r.Context()
		data home.HomeData
	)

	data.Session, _ = ctx.Value("session").(*cache.Session)
	if data.Session == nil {
		api.respond(w, r, http.StatusOK, home.HomePage(&data))
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

	api.respond(w, r, http.StatusOK, home.HomePage(&data))
}

func (api *API) handlePlayPage(w http.ResponseWriter, r *http.Request) {

	var (
		ctx  = r.Context()
		data play.PlayData
	)

	data.Session, _ = ctx.Value("session").(*cache.Session)
	if data.Session == nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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

	errs.Go(func() (err error) {

		chatMessages, err := api.Service.CQL.SelectChatMessages(errsCtx, roomID)
		if err != nil {
			return errors.Wrap(err, "cannot get user")
		}

		for _, chatMsg := range chatMessages {

			author, err := api.Service.GetUser(ctx, chatMsg.AuthorID)
			if err != nil {
				return errors.Wrap(err, "cannot get author of chat-message")
			}

			data.ChatMessages = append(data.ChatMessages, &pubsub.ChatMessage{
				ChatMessage: *chatMsg,
				Author:      author,
			})
		}

		return nil
	})

	if err := errs.Wait(); err != nil {
		api.err(w, r, errors.WithStack(err))
		return
	}

	api.respond(w, r, http.StatusOK, play.PlayPage(&data))
}

func (api *API) handlePlayPageWebSocket(conn *websocket.Conn) {

	var (
		r   = conn.Request()
		ctx = r.Context()
	)

	session, _ := ctx.Value("session").(*cache.Session)
	if session == nil {
		return
	}

	roomID, err := server.ParseUUID(r.URL.Query().Get("r"))
	if err != nil {
		api.err(conn, r, errors.Wrap(err, "cannot parse room-ID"))
		return
	}

	user, err := api.Service.GetUser(ctx, session.UserID)
	if err != nil {
		api.err(conn, r, errors.Wrap(err, "cannot get user"))
		return
	}

	go func() {
		err := api.Service.PubSub.SubRoom(ctx, roomID,
			func(event any) {
				switch e := event.(type) {
				case *pubsub.ChatMessage:
					api.respond(conn, r, 0, play.ChatMessage(e))
				}
			})
		api.err(conn, r, errors.Wrap(err, "cannot subscribe to room events"))
		conn.Close()
	}()

	var (
		msg     []byte
		payload any
		action  = func() {}
	)

	for {

		msgType, err := wsReceive(conn, &msg)
		if err != nil {
			if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
				return
			}

			api.err(conn, r, errors.Wrap(err, "cannot receive message"))
			continue
		}

		switch msgType {
		case "send-chat-message":
			payload, action = &sendChatMessageReq{}, func() {
				err := api.Service.SendChatMessage(ctx, roomID, user, payload.(*sendChatMessageReq).Content)
				api.err(conn, r, errors.Wrap(err, "cannot send chat-message"))
			}
		}

		if err := json.Unmarshal(msg, &payload); err != nil {
			api.err(conn, r, &server.UserError{
				Type:    server.UserErrorTypeJSONInvalid,
				Message: err.Error()})
			continue
		}

		action()
	}
}
