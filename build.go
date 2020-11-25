package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

// RunBuild builds the server (either for deployment or for local use)
func RunBuild(args []string) {
	// Remove fragmenta build from args list
	args = args[2:]

	if len(args) > 0 {
		buildDeployServer()
	} else {
		buildLocalServer()
	}

}

// buildServer removes the old binary and rebuilds the server
// - this is simply a wrapper around go build, you can instead
// run go build server.go directly if you prefer.
func buildServer(server string, env []string) error {

	// If environment is empty, we are doing a local build
	// localBuild := (len(env) == 0)

	// First run go fmt on any packages below root
	srcPath := "./..."
	log.Printf("Running go fmt at %s", srcPath)
	result, err := runCommand("go", "fmt", srcPath)
	if err != nil {
		log.Printf("Error running fmt %s", err)
		return err
	}
	if len(result) > 0 {
		log.Printf(string(result))
	}

	// Then remove old binary
	_, err = os.Stat(server)
	if err == nil {
		err = os.Remove(server)
		if err != nil {
			log.Printf("Error removing server %s", err)
		}
	}

	// Build a new binary
	started := time.Now()
	log.Printf("Building server at %s", server)

	// Start with build command
	args := []string{"build"}

	// Add output location
	args = append(args, `-o`)
	args = append(args, server)

	// Finally add the path to server.go
	args = append(args, serverCompilePath("."))

	// Call the command
	cmd := exec.Command("go", args...)
	cmd.Stderr = os.Stdout

	if env != nil {
		cmd.Env = env
	}

	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error running build %s\n%s", err, string(output))
		return err
	}

	// Record the output of our build (success or failure)
	if len(output) == 0 {
		log.Printf("Build completed successfully in %s", time.Since(started).String())
	} else {
		log.Printf(string(output))
	}

	return nil

}

// buildLocalServer builds the server at a local path
func buildLocalServer() {
	buildServer(localServerPath("."), nil)
}

// buildDeployServer builds the server for deployment on linux.
func buildDeployServer() {
	env := append(os.Environ(), "GOOS=linux")
	env = append(env, "GOARCH=amd64")

	// When compiling with cgo, we get this error:
	// ./bin/server: error while loading shared libraries: /usr/lib/libSystem.B.dylib: cannot open shared object file: No such file or directory
	env = append(env, "CGO_ENABLED=0")
	buildServer(serverPath("."), env)
}
