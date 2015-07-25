package main

import (
	"log"
	"os"
)

// If it exists, simply run the binary in bin/deploy
// (this might be a shell script which runs ansible for example)
func runDeploy(args []string) {

	// Build deploy server
	buildDeployServer()

	deploy := "./bin/deploy"

	_, err := os.Stat(deploy)
	if err != nil {
		log.Printf("Could not find deploy script at %s", deploy)
		return
	}
    
    // Default to development
    mode := "development"
	if len(args) == 3 {
	    mode = args[2]
    }

	log.Printf("Running deploy from " + deploy)
	result, err := runCommand(deploy,mode)
	if err != nil {
		log.Printf("Error running deploy", err)
		return
	}

	log.Printf(string(result))
}
