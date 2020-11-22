package main

import (
	"log"
	"path/filepath"
	"regexp"
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

	log.Printf("Tests Complete\n" + colorizeResults(string(result)))
}

// Color constants used for writing in colored output
const (
	ColorNone  = "\033[0m" // Use this to clear output
	ColorRed   = "\033[31m"
	ColorGreen = "\033[32m"
	ColorAmber = "\033[33m"
	ColorCyan  = "\033[1;36m"
)

// colorizeResults colours lines according to the status of the tests
func colorizeResults(results string) string {

	// First remove empty FAIL lines
	emptyFail := regexp.MustCompile(`(?m)^FAIL$`)
	results = emptyFail.ReplaceAllString(results, "")

	okRE := regexp.MustCompile(`(?m)^ok.*$`)
	results = okRE.ReplaceAllString(results, ColorGreen+"$0"+ColorNone)

	missingRE := regexp.MustCompile(`(?m)^\?.*$`)
	results = missingRE.ReplaceAllString(results, ColorAmber+"$0"+ColorNone)

	failRE := regexp.MustCompile(`(?m)^(--- )?FAIL.*$`)
	results = failRE.ReplaceAllString(results, ColorRed+"$0"+ColorNone+"\n")

	return results
}
