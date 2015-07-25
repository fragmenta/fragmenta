package main

import (
	"log"
	"os"
	"path"
)

/*
// What this should do is the following:

// first arg - site type
// This can be a valid github.com/fragmenta/fragmenta-* repo (so e.g. cms for fragmenta-cms)
// OR any valid git repo path - uses that repo to call git clone and create the new website

Possible examples:

cms - full cms
blog - simple blog - hmm, do we need this? yes would be nice
app - simple app with no pages, tags, only users + auth
api - simple app with no users or anything, so that they can build an api on it.


Start with CMS, and work downwards by stripping out parts.

// Second arg - path - this must be a path within GOPATH - either a full path (in which case remove gopath)
// or a path within GOPATH/src/

// Usage examples:
// fragmenta new cms github.com/kennygrant/mycms
// fragmenta new github.com/kennygrant/app-template github.com/kennygrant/myapp

*/

// Usage: fragmenta new [cms] [path]
//
func runNew(args []string) {
	// Remove fragmenta backup from args list
	args = args[2:]

	siteType := "spartan"
	if len(args) > 0 {
		siteType = args[0]
	}

	projectPath := "."
	if len(args) > 1 {
		projectPath = args[1]
	}

	if isValidProject(projectPath) && len(args) < 2 {
		log.Printf("\nA fragmenta project already exists here\n")
		return
	}

	tmpPath := cloneExamples(projectPath)

	// Copy out the desired site to path - unfortunately golang doesn't
	// provide a simple facility for this - we assume unix
	result, err := runCommand("cp", "-r", path.Join(tmpPath, siteType), projectPath)
	if err != nil {
		log.Printf("Error copying example site %s", err)
		return
	}
	log.Printf("%s", string(result))

	// Remove our tmp clone of examples
	os.RemoveAll(tmpPath)

	// We then need to rewrite refs in any files to refer to this project path
	// within godir ... think about how this would work...
	// alternative is to leave them with relative imports

	// Run the server with this path
	runServer(projectPath)

	// Open a browser at the default port
	result, err = runCommand("open", "http://localhost:3000")
	if err != nil {
		log.Printf("Error opening site %s", err)
		return
	}
	log.Printf("%s", string(result))

}
