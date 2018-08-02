// A command line tool for fragmenta which can be used to build and run websites
// this tool calls subcommands for most of the work, usually one command per file in this pkg
// See docs at http://godoc.org/github.com/fragmenta/fragmenta

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// The version of this tool
	fragmentaVersion = "1.5.7"

	// Used for outputting console messages
	fragmentaDivider = "\n------\n"
)

// Modes used for setting the config used
const (
	ModeProduction  = "production"
	ModeDevelopment = "development"
	ModeTest        = "test"
)

var (
	// ConfigDevelopment holds the development config from fragmenta.json
	ConfigDevelopment map[string]string

	// ConfigProduction holds development config from fragmenta.json
	ConfigProduction map[string]string

	// ConfigTest holds the app test config from fragmenta.json
	ConfigTest map[string]string
)

// main - parse the command line arguments and respond
func main() {

	// Log time as well as date
	log.SetFlags(log.Ltime)

	// Parse commands
	args := os.Args
	command := ""

	if len(args) > 1 {
		command = args[1]
	}

	// We assume the project path is the current directory (for now)
	projectPath, err := filepath.Abs(".")
	if err != nil {
		log.Printf("Error getting path %s", err)
		return
	}

	// If this is a valid project, read the config, else continue
	if isValidProject(projectPath) {
		readConfig(projectPath)
	}

	switch command {

	case "new", "n":
		RunNew(args)

	case "version", "v":
		ShowVersion()

	case "help", "h", "wat", "?":
		ShowHelp(args)

	case "server", "s":
		if requireValidProject(projectPath) {
			RunTests(nil)
			RunServer(projectPath)
		}

	case "test", "t":
		if requireValidProject(projectPath) {
			// Remove fragmenta test from args list
			args = args[2:]
			RunTests(args)
		}

	case "build", "B":
		if requireValidProject(projectPath) {
			RunTests(nil)
			RunBuild(args)
		}

	case "generate", "g":
		if requireValidProject(projectPath) {
			RunGenerate(args)
		}

	case "migrate", "m":
		if requireValidProject(projectPath) {
			RunMigrate(args)
		}

	case "backup", "b":
		if requireValidProject(projectPath) {
			RunBackup(args)
		}

	case "restore", "r":
		if requireValidProject(projectPath) {
			RunRestore(args)
		}

	case "deploy", "d":
		if requireValidProject(projectPath) {
			RunDeploy(args)
		}
	case "":
		// Special case no commands to build and run the server
		if requireValidProject(projectPath) {
			RunTests(nil)
			RunServer(projectPath)
		}
	default:
		// Command not recognised so show the help
		ShowHelp(args)
	}

}

// ShowVersion shows the version of this tool
func ShowVersion() {
	helpString := fragmentaDivider
	helpString += fmt.Sprintf("Fragmenta version: %s", fragmentaVersion)
	helpString += fragmentaDivider
	log.Print(helpString)
}

// ShowHelp shows the help for this tool
func ShowHelp(args []string) {
	helpString := fragmentaDivider
	helpString += fmt.Sprintf("Fragmenta version: %s", fragmentaVersion)
	helpString += "\n  fragmenta version -> display version"
	helpString += "\n  fragmenta help -> display help"
	helpString += "\n  fragmenta new [app|cms|URL] path/to/app -> creates a new app from the repository at URL at the path supplied"
	helpString += "\n  fragmenta -> builds and runs a fragmenta app"
	helpString += "\n  fragmenta server -> builds and runs a fragmenta app"
	helpString += "\n  fragmenta test  -> run tests"
	helpString += "\n  fragmenta migrate -> runs new sql migrations in db/migrate"
	helpString += "\n  fragmenta backup [development|production|test] -> backup the database to db/backup"
	helpString += "\n  fragmenta restore [development|production|test] -> backup the database from latest file in db/backup"
	helpString += "\n  fragmenta deploy [development|production|test] -> build and deploy using bin/deploy"
	helpString += "\n  fragmenta generate resource [name] [fieldname]:[fieldtype]* -> creates resource CRUD actions and views"
	helpString += "\n  fragmenta generate migration [name] -> creates a new named sql migration in db/migrate"

	helpString += fragmentaDivider
	log.Print(helpString)
}

// Ideally all these paths could be configured,
// rather than baking assumptions about project structure into the tool

// serverName returns the path of the cross-compiled target server binary
// this does not end in .exe as we assume a target of linux
func serverName() string {
	return "fragmenta-server"
}

// localServerName returns a server name for the local server binary (prefixed with local)
func localServerName() string {
	if isWindows() {
		return serverName() + "-local.exe"
	}
	return serverName() + "-local"
}

// localServerPath returns the local server binary for running on the dev machine locally
func localServerPath(projectPath string) string {
	return filepath.Join(projectPath, "bin", localServerName())
}

// serverPath returns the cross-compiled server binary
func serverPath(projectPath string) string {
	return filepath.Join(projectPath, "bin", serverName())
}

// serverCompilePath returns the server entrypoint
func serverCompilePath(projectPath string) string {
	return filepath.Join(projectPath, "server.go")
}

// srcPath returns the path for Go code within the project
func srcPath(projectPath string) string {
	return filepath.Join(projectPath, "src")
}

// publicPath returns the path for the public directory of the web application
func publicPath(projectPath string) string {
	return filepath.Join(projectPath, "public")
}

// configPath returns the path for the fragment config file (required)
func configPath(projectPath string) string {
	return filepath.Join(secretsPath(projectPath), "fragmenta.json")
}

// secretsPath returns the path for secrets
func secretsPath(projectPath string) string {
	return filepath.Join(projectPath, "secrets")
}

// templatesPath returns the path for templates
func templatesPath() string {
	path := filepath.Join(goPath(), "src", "github.com", "fragmenta", "fragmenta", "templates")
	return os.ExpandEnv(path)
}

// dbMigratePath returns a path to store database migrations
func dbMigratePath(projectPath string) string {
	return filepath.Join(projectPath, "db", "migrate")
}

// dbBackupPath returns a path to store database backups
func dbBackupPath(projectPath string) string {
	return filepath.Join(projectPath, "db", "backup")
}

// projectPathRelative returns the relative path
func projectPathRelative(projectPath string) string {
	goSrc := filepath.Join(goPath(), "src")
	return strings.Replace(projectPath, goSrc, "", 1)
}

// goPath returns the setting of env variable $GOPATH
// or $HOME/go if no $GOPATH is set.
func goPath() string {
	// Get the first entry in gopath
	paths := filepath.SplitList(os.ExpandEnv("$GOPATH"))
	if len(paths) > 0 && paths[0] != "" {
		return paths[0]
	}
	return filepath.Join(homePath(), "go")
}

// homePath returns the user's home directory
func homePath() string {
	if isWindows() {
		return os.ExpandEnv("$userprofile")
	}
	return os.ExpandEnv("$HOME")
}

// RunServer runs the server
func RunServer(projectPath string) {
	ShowVersion()

	log.Println("Building server...")
	err := buildServer(localServerPath(projectPath), nil)

	if err != nil {
		log.Printf("Error building server: %s", err)
		return
	}

	log.Println("Launching server...")
	cmd := exec.Command(localServerPath(projectPath))
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	err = cmd.Start()
	if err != nil {
		log.Println(err)
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	cmd.Wait()

}

// runCommand runs a command with exec.Command
func runCommand(command string, args ...string) ([]byte, error) {

	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err
	}

	return output, nil
}

// requireValidProject returns true if we have a valid project at projectPath
func requireValidProject(projectPath string) bool {
	if isValidProject(projectPath) {
		return true
	}

	log.Printf("No fragmenta project found at this path\n")
	return false
}

// isValidProject returns true if this is a valid fragmenta project (checks for server.go file and config file)
func isValidProject(projectPath string) bool {

	// Make sure we have server.go at root of this dir
	_, err := os.Stat(serverCompilePath(projectPath))
	if err != nil {
		return false
	}

	return true
}

// fileExists returns true if this file exists
func fileExists(p string) bool {
	_, err := os.Stat(p)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// readConfig reads our config file and set up the server accordingly
func readConfig(projectPath string) error {
	configPath := configPath(projectPath)

	// Read the config json file
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("Error opening config at %s\n%s", configPath, err)
		return err
	}

	var data map[string]map[string]string
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Printf("Error parsing config %s %v", configPath, err)
		return err
	}

	ConfigDevelopment = data["development"]
	ConfigProduction = data["production"]
	ConfigTest = data["test"]

	return nil
}

// isWindows returns true if the Go architecture target (GOOS) is windows
func isWindows() bool {
	if runtime.GOOS == "windows" {
		return true
	}
	return false
}
