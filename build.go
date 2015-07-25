package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

func runBuild(args []string) {
	// Remove fragmenta build from args list
	args = args[2:]

	if len(args) > 0 {
		buildDeployServer()
	} else {
		buildLocalServer()
	}

}

func buildServer(server string, env []string) error {

	// Remove old binary
	_, err := os.Stat(server)
	if err == nil {
		err = os.Remove(server)
		if err != nil {
			log.Printf("Error removing server %s", err)
		}
	}

	// If we have a goimports, run that, if not run go fmt
	_, err = os.Stat(os.ExpandEnv("$GOPATH/bin/goimports"))

	if err == nil {
		// Go imports behaviour differs from go fmt
		srcPath := "./src"
		log.Printf("Running goimports at %s", srcPath)
		result, err := runCommand("goimports", "-w", srcPath)
		if err != nil {
			log.Printf("Error running goimports %s", err)
			return err
		}
		if len(result) > 0 {
			log.Printf(string(result))
		}

	} else {
		srcPath := "./src/..."
		// Run go fmt on any packages with src
		log.Printf("Running go fmt at %s", srcPath)
		result, err := runCommand("go", "fmt", srcPath)
		if err != nil {
			log.Printf("Error running fmt %s", err)
			return err
		}
		if len(result) > 0 {
			log.Printf(string(result))
		}
	}

	// Build new binary
	log.Printf("Building server at %s", server)
	started := time.Now()

	log.Printf("CMD %s %s %s %s %s", "go", "build", "-o", server, appPath("."))

	// NB we set environment here because we may be cross=compiling
	cmd := exec.Command("go", "build", "-o", server, appPath("."))
	cmd.Stderr = os.Stdout

	if env != nil {
		cmd.Env = env
	}

	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error running build %s\n%s", err, string(output))
		return err
	}

	// We should also be rebuilding assets here
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
