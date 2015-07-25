package [[ .fragmenta_resource ]]_actions

import (
	"net/http"
    "strings"

	"github.com/fragmenta/model/url"
	"github.com/fragmenta/router"
	"github.com/fragmenta/view"
	"github.com/fragmenta/view/helpers"

	"bitbucket.org/kennygrant/frithandco/src/[[ .fragmenta_resources ]]"
)

// Serve a get request at /[[ .fragmenta_resources ]]
//
//
func HandleIndex(context *router.Context) {

	// Authorize
	if !context.Authorize() {
		view.RenderStatus(context.Writer, http.StatusUnauthorized)
		return
	}

	// Setup context for template
	view := view.New(context)

	// Build a query
    q := [[ .fragmenta_resources ]].Query()
    
    // Show only published (defaults to showing all)
    // q.Apply(status.WherePublished)

    
	// Order by required order, or default to id asc
	switch context.Param("order") {

	case "1":
		q.Order("created desc")

	case "2":
		q.Order("updated desc")

	case "3":
		q.Order("name asc")

	default:
		q.Order("id asc")
        
	}

    
	// Filter if necessary - this assumes name and summary cols
	filter := context.Param("filter")
	if len(filter) > 0 {
		filter = strings.Replace(filter, "&", "", -1)
		filter = strings.Replace(filter, " ", "", -1)
		filter = strings.Replace(filter, " ", " & ", -1)
		q.Where("( to_tsvector(name) || to_tsvector(summary) @@ to_tsquery(?) )", filter)
	}
    

	// Fetch the [[ .fragmenta_resources ]]
    var results []*[[ .fragmenta_resources ]].[[ .Fragmenta_Resource ]]
	err := q.Fetch(&results)
	if err != nil {
		context.Log("#error indexing [[ .fragmenta_resources ]] %s", err)
		view.RenderError(context.Writer, err)
		return
	}
    
    

	// Serve template
    view.AddKey("filter", filter)
	view.AddKey("[[ .fragmenta_resources ]]", results)
	// Can we add these programatically?
	view.AddKey("admin_links", helpers.LinkTo("Create [[ .fragmenta_resource ]]", url.Create([[ .fragmenta_resources ]].New())))

	view.Render(context.Writer)

}
