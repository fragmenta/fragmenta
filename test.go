package main

import (
	"fmt"
	"log"
)

// RunTests runs all tests below the current path or the path specified
func RunTests(args []string) {

	testDir := "./..."

	if len(args) > 0 {
		testDir = fmt.Sprintf("./%s", args[0])
	}

	log.Printf("Running tests at %s", testDir)

	result, err := runCommand("go", "test", testDir)
	if err != nil {
		log.Printf("Error running tests %s", err)
	}

	log.Printf(string(result))
}
