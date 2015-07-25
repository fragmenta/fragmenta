package main

import (
	"fmt"
	"log"
)

func runTests(args []string) {
	// Remove fragmenta test from args list
	args = args[2:]

	test_dir := "./src/..."

	if len(args) > 0 {
		test_dir = fmt.Sprintf("./src/%s", args[0])
	}

	log.Printf("Running tests at %s", test_dir)

	result, err := runCommand("go", "test", "-v", test_dir)
	if err != nil {
		log.Printf("Error running tests %s", err)
	}

	log.Printf(string(result))
}
