package main

import "log"

// fragmenta new creates a new fragmenta project given the argument
// Usage: fragmenta new [app|cms|api| valid repo path e.g. github.com/fragmenta/fragmenta-cms]
func runNew(args []string) {

	// we can then run go get on that repo path, and ideally cd and start up the server immediately afterward
	// the keywords should find a repo at a defined path under github.com/fragmenta/ and use that

	// Remove fragmenta backup from args list
	args = args[2:]

	// Until template projects are up, just return here
	log.Printf("\nNo template projects found\n")
	return
	/*


		  // Check we are not already in a valid project?
			projectPath := "."

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
	*/
}
