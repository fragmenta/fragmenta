package [[ .fragmenta_resource ]]_actions

import (
	"net/http"

	"github.com/fragmenta/model/url"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"
	"github.com/fragmenta/view/helpers"

	".."
)

// Serve a get request at /[[ .fragmenta_resources ]]/1/update (show form to update)
func HandleUpdateShow(context *router.Context) {
	// Setup context for template
	view := view.New(context)

	[[ .fragmenta_resource ]], err := [[ .fragmenta_resources ]].Find(context.ParamInt("id"))
	if err != nil {
		context.Log.Error("Error finding [[ .fragmenta_resource ]] %s", err)
		view.RenderStatus(context.Writer, http.StatusNotFound)
		return
	}

	// Authorize
	if !context.Authorize([[ .fragmenta_resource ]]) {
		view.RenderStatus(context.Writer, http.StatusUnauthorized)
		return
	}

	view.AddKey("[[ .fragmenta_resource ]]", [[ .fragmenta_resource ]])
	view.AddKey("admin_links", helpers.LinkTo("Destroy [[ .Fragmenta_Resource ]]", url.Destroy([[ .fragmenta_resource ]]), "method=delete"))

	view.Render(context.Writer)
}

// POST or PUT /[[ .fragmenta_resources ]]/1/update
func HandleUpdate(context *router.Context) {
	// Setup context for template
	view := view.New(context)

	// Find the [[ .fragmenta_resource ]]
	[[ .fragmenta_resource ]], err := [[ .fragmenta_resources ]].Find(context.ParamInt("id"))
	if err != nil {
		context.Log.Error("Error finding [[ .fragmenta_resource ]] %s", err)
		view.RenderStatus(context.Writer, http.StatusNotFound)
		return
	}

	// Authorize
	if !context.Authorize([[ .fragmenta_resource ]]) {
		view.RenderStatus(context.Writer, http.StatusUnauthorized)
		return
	}

	// Update the [[ .fragmenta_resource ]]
	params, err := context.Params()
	if err != nil {
		view.RenderError(context.Writer, err)
		return
	}

	err = [[ .fragmenta_resource ]].Update(params.Map())
    if err != nil {
		view.RenderError(context.Writer, err)
		return
	}


	// Redirect to [[ .fragmenta_resource ]]
	router.Redirect(context, url.Show([[ .fragmenta_resource ]]) )
}
