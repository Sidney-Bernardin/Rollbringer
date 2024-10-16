// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package play

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"rollbringer/internal"
	"rollbringer/internal/views/games"
)

func navBar(page *internal.PlayPage) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div x-data class=\"nav-bar\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if page.Game != nil {
			templ_7745c5c3_Err = rollCalculator().Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(" <div class=\"rolls\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, roll := range page.Game.Rolls {
				templ_7745c5c3_Err = games.Roll(roll).Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<button class=\"games-btn icon-btn\" @click=\"$dispatch(&#39;show-games-modal&#39;)\"><iconify-icon icon=\"ic:baseline-meeting-room\"></iconify-icon></button></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

func rollCalculator() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
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
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<form class=\"roll-calculator\" x-data=\"{ selectedDiceTypes: [] }\" :hx-vals=\"`{\n            &#34;EVENT&#34;: &#34;CREATE_ROLL_REQUEST&#34;,\n            &#34;dice_types&#34;: [${selectedDiceTypes}]\n        }`\" ws-send><button class=\"submit-btn icon-btn\"><iconify-icon icon=\"fa6-solid:dice-d20\"></iconify-icon></button><div class=\"expander\"><div class=\"inner\"><select @change=\"selectedDiceTypes.push(Number($el.value)); $el.value=&#39;&#39;\"><option value=\"\">Select Dice</option> <option value=\"4\">D4</option> <option value=\"6\">D6</option> <option value=\"8\">D8</option> <option value=\"10\">d10</option> <option value=\"12\">d12</option> <option value=\"20\">d20</option></select><div class=\"dice-preview\"><template x-for=\"dieType in selectedDiceTypes\"><button class=\"icon-btn\" @click.prevent=\"$el.remove()\"><iconify-icon :icon=\"`mdi:dice-d${dieType}`\"></iconify-icon></button></template></div><input type=\"text\" name=\"modifier\" placeholder=\"Modifier\"></div></div></form>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
