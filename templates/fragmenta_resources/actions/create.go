package [[ .fragmenta_resource ]]_actions

import (
	"net/http"

	"github.com/fragmenta/model/url"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"

    "bitbucket.org/kennygrant/frithandco/src/[[ .fragmenta_resources ]]"
)

// GET [[ .fragmenta_resources ]]/create
func HandleCreateShow(context *router.Context) {

	// Authorize
	if !context.Authorize() {
		view.RenderStatus(context.Writer, http.StatusUnauthorized)
		return
	}

	// Setup
	view := view.New(context)
	[[ .fragmenta_resource ]] := [[ .fragmenta_resources ]].New()
	view.AddKey("[[ .fragmenta_resource ]]", [[ .fragmenta_resource ]])

	// Serve
	view.Render(context.Writer)
}

// POST [[ .fragmenta_resources ]]/create
func HandleCreate(context *router.Context) {

	// Authorize
	if !context.Authorize() {
		view.RenderStatus(context.Writer, http.StatusUnauthorized)
		return
	}

	// Setup context
	params, err := context.Params()
	if err != nil {
		view.New(context).RenderError(context.Writer, err)
		return
	}

	id, err := [[ .fragmenta_resources ]].Create(params.Map())
	if err != nil {
		context.Log("#error creating [[ .fragmenta_resource ]],%s", err)
    	view.New(context).RenderError(context.Writer, err)
    	return
	}
    
    // Log creation
    context.Log("#info Created [[ .fragmenta_resource ]] id,%d", id)
	
	// Redirect to the new [[ .fragmenta_resource ]]
	m, err := [[ .fragmenta_resources ]].Find(id)
	if err != nil {
		context.Log("#error creating [[ .fragmenta_resource ]],%s", err)
	}

	router.Redirect(context, url.Index(m))
}


