// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.857
package pages

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"fmt"
	"github.com/google/uuid"
	"rollbringer/src/api/views"
	"rollbringer/src/domain/services/accounts"
	"rollbringer/src/domain/services/play"
)

type PlayData struct {
	Session *accounts.Session

	Room      *play.Room
	RoomUsers []*accounts.User

	Boards      []*play.Board
	BoardsUsers map[uuid.UUID][]*accounts.User
}

func Play(page *PlayData) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(string(page.Room.Name))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/pages/play.templ`, Line: 26, Col: 34}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, " | Rollbringer</title><link rel=\"icon\" type=\"image/x-icon\" href=\"/static/favicon.png\"><link rel=\"stylesheet\" type=\"text/css\" href=\"/static/styles/play.css\"><script src=\"/static/play.js\" defer></script></head><body x-data hx-ext=\"ws\" ws-connect=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(fmt.Sprintf("/rooms/%s/ws", page.Room.ID))
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/pages/play.templ`, Line: 34, Col: 57}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if page.Session != nil {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, " hx-headers=\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			var templ_7745c5c3_Var4 string
			templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(templ.JSONString(map[string]any{"CSRF-Token": page.Session.CSRFToken}))
			if templ_7745c5c3_Err != nil {
				return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/pages/play.templ`, Line: 36, Col: 87}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "\"")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "><div class=\"layout\"><div class=\"nav-bar\"><a href=\"/\">Home</a></div><div class=\"document-viewer\"></div><div class=\"boards\"></div><div class=\"materials\"><div class=\"boards\"><div class=\"tool-bar\"><button @click=\"$refs.createBoard.showModal()\"><iconify-icon icon=\"material-symbols:add\"></iconify-icon> Create Board</button></div><div class=\"grid\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, board := range page.Boards {
			templ_7745c5c3_Err = views.BoardCard(board, page.BoardsUsers[board.ID]).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "</div></div></div><div class=\"social\"><div class=\"users\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, user := range page.RoomUsers {
			templ_7745c5c3_Err = views.UserBubble(user).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "</div><div class=\"chat\"><div class=\"messages\"><div class=\"inner\"><div class=\"anchor\"></div></div></div></div><form ws-send><input type=\"hidden\" name=\"operation\" value=\"chat\"> <input type=\"text\" name=\"message\" placeholder=\"Hello, World!\"> <button><iconify-icon icon=\"material-symbols:send\"></iconify-icon> Send</button></form></div><span class=\"gutter g1\"></span> <span class=\"gutter g2\"></span> <span class=\"gutter g3\"></span> <span class=\"gutter g4\"></span><dialog class=\"create-board\" x-ref=\"createBoard\"><header><h1>Create Board</h1><button @click=\"$refs.createBoard.close()\"><iconify-icon icon=\"material-symbols:close\"></iconify-icon> Close</button></header><hr><form ws-send><input type=\"hidden\" name=\"operation\" value=\"create-board\"> <input type=\"text\" name=\"name\" placeholder=\"Name\"> <button>Create</button></form></dialog></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
