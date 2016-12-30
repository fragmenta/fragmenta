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
func buildServer(server string, env []string) error {

	// Remove old binary for server
	_, err := os.Stat(server)
	if err == nil {
		err = os.Remove(server)
		if err != nil {
			log.Printf("Error removing server %s", err)
		}
	}

	// Run go fmt on any packages below root
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

	// Build new binary for server
	log.Printf("Building server at %s", server)
	started := time.Now()

	args := []string{"build", "-o", server, serverCompilePath(".")}

	// If a local build with no environment settings, use go build -i
	if len(env) == 0 {
		args = []string{"build", "-i", "-o", server, serverCompilePath(".")}
	}

	// Call go build with -i to install artefacts for local builds, and -o to output to ./bin
	// log.Printf("  %s", args)
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

func buildLocalServer() {
	buildServer(localServerPath("."), nil)
}

func buildDeployServer() {
	env := append(os.Environ(), "GOOS=linux")
	env = append(env, "GOARCH=amd64")

	// When compiling with cgo, we get this error:
	// ./bin/server: error while loading shared libraries: /usr/lib/libSystem.B.dylib: cannot open shared object file: No such file or directory
	env = append(env, "CGO_ENABLED=0")
	buildServer(serverPath("."), env)
}
