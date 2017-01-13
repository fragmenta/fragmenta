package main

import (
	"log"
	"path/filepath"
	"strings"
)

// RunTests runs all tests below the current path or the path specified
// this is a simplistic wrapper around the go test tool
func RunTests(args []string) {

	// Run tests on the src dir, this skips root server tests but also skips vendor tests
	testDir := strings.Join([]string{".", "src", "..."}, string(filepath.Separator))

	if len(args) > 0 {
		testDir = args[0]
	}

	log.Printf("Running tests at %s", testDir)

	result, err := runCommand("go", "test", testDir)
	if err != nil {
		log.Printf("Error running tests %s", err)
	}

	log.Printf(string(result))
}
