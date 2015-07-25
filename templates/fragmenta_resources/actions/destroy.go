package [[ .fragmenta_resource ]]_actions

import (
	"net/http"

	"github.com/fragmenta/model/url"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

	".."
)

// POST /[[ .fragmenta_resources ]]/1/destroy
func HandleDestroy(context *router.Context) {

	// Set the [[ .fragmenta_resource ]] on the context for checking
	[[ .fragmenta_resource ]], err := [[ .fragmenta_resources ]].Find(context.ParamInt("id"))
	if err != nil {
		view.RenderStatus(context.Writer, http.StatusNotFound)
		return
	}

	// Authorize
	if !context.Authorize([[ .fragmenta_resource ]]) {
		view.RenderStatus(context.Writer, http.StatusUnauthorized)
		return
	}

	// Destroy the [[ .fragmenta_resource ]]
	[[ .fragmenta_resource ]].Destroy()

	// Redirect to [[ .fragmenta_resources ]] root
	router.Redirect(context, url.Index([[ .fragmenta_resource ]]))
}
