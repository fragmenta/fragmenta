package [[ .fragmenta_resource ]]_actions

import (
	"net/http"

	"github.com/fragmenta/model/url"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"
	"github.com/fragmenta/view/helpers"

	"bitbucket.org/kennygrant/frithandco/src/[[ .fragmenta_resources ]]"
)

// Serve a get request at /[[ .fragmenta_resources ]]/1

func HandleShow(context *router.Context) {

	// Setup context for template
	view := view.New(context)

	[[ .fragmenta_resource ]], err := [[ .fragmenta_resources ]].Find(context.ParamInt("id"))
	if err != nil {
		view.RenderError(context.Writer, err)
		return
	}

	// Authorize
	if !context.Authorize([[ .fragmenta_resource ]]) {
		view.RenderStatus(context.Writer, http.StatusUnauthorized)
		return
	}

	// Serve template
	view.AddKey("[[ .fragmenta_resource ]]", [[ .fragmenta_resource ]])
	view.AddKey("admin_links", helpers.LinkTo("Edit [[ .Fragmenta_Resource ]]", url.Update([[ .fragmenta_resource ]])))

	view.Render(context.Writer)

}
