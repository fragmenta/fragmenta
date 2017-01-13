package main

import (
	"log"
	"path/filepath"
)

// RunTests runs all tests below the current path or the path specified
// this is a simplistic wrapper around the go test tool
func RunTests(args []string) {

	// Send two paths to go - root and src
	// we do this to ignore the vendor dir
	testDirs := filepath.Join(". .", "src", "...")

	if len(args) > 0 {
		testDirs = args[0]
	}

	log.Printf("Running tests at %s", testDirs)

	result, err := runCommand("go", "test", testDirs)
	if err != nil {
		log.Printf("Error running tests %s", err)
	}

	log.Printf(string(result))
}
