package main

import (
	"log"
	"os"

	// We depend on assets in order to build - how to tell user doesn't want assets?
	"github.com/fragmenta/assets"
)

// If it exists, simply run the binary in bin/deploy
// (this might be a shell script which runs ansible for example)
func runDeploy(args []string) {

	// Build our app assets and update secrets/assets.json
	buildAssets()

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
	result, err := runCommand(deploy, mode)
	if err != nil {
		log.Printf("Error running deploy %s", err)
		return
	}

	log.Printf(string(result))
}

// Compile the app assets before a deploy, so that they're available for production use
func buildAssets() {
	log.Printf("Compiling assets...")
	err := assets.New(true).Compile("src", "public")
	if err != nil {
		log.Fatalf("#error compiling assets %s", err)
	}
}
