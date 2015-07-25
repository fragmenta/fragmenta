// Command line tool for fragmenta which can be used to build and run websites
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"sort"
)

const VERSION = "1.0"
const DIVIDER = "\n------\n"
const APP_NAME = "fragmenta-server" // this should be settable in config
const GO = "/usr/local/go/bin/go"

var ConfigDevelopment map[string]string
var ConfigProduction map[string]string
var ConfigTest map[string]string

// NB NEVER use fragmenta, for obvious reasons - we use the process name to kill
// can we kill our launched processes by pid instead?
func serverName() string {
	return "fragmenta-server" // for now, should use configs
}

func localServerPath(projectPath string) string {
	return fmt.Sprintf("%s/bin/%s-local", projectPath, serverName())
}

func serverPath(projectPath string) string {
	return fmt.Sprintf("%s/bin/%s", projectPath, serverName())
}

func appPath(projectPath string) string {
	return projectPath + "/src/app"
}

func configPath(projectPath string) string {
	return projectPath + "/secrets/fragmenta.json"
}

func templatesPath() string {
	return os.ExpandEnv("$GOPATH/src/github.com/fragmenta/fragmenta/templates")
}

// Parse the command line arguments and respond
func main() {

	log.SetFlags(log.Ltime)

	args := os.Args
	command := ""

	if len(args) > 1 {
		command = args[1]
	}

	// We should intelligently read project path depending on the command?
	// Or just assume we act on the current directory?
	projectPath := "."

	// Will we ever act on another path?
	if isValidProject(projectPath) {
		readConfig(projectPath)
	}

	switch command {

	case "new", "n":
		runNew(args)

	case "version", "v":
		showVersion()

	case "help", "h", "wat", "?":
		showHelp(args)

	case "server", "s":
		if requireValidProject(projectPath) {
			runServer(projectPath)
		}

	case "test", "t":
		if requireValidProject(projectPath) {
			runTests(args)
		}

	case "build", "B":
		if requireValidProject(projectPath) {
			runBuild(args)
		}

	case "generate", "g":
		if requireValidProject(projectPath) {
			runGenerate(args)
		}

	case "migrate", "m":
		if requireValidProject(projectPath) {
			runMigrate(args)
		}

	case "backup", "b":
		if requireValidProject(projectPath) {
			runBackup(args)
		}

	case "restore", "r":
		if requireValidProject(projectPath) {
			runRestore(args)
		}

	case "deploy", "d":
		if requireValidProject(projectPath) {
			runDeploy(args)
		}

	default:
		if requireValidProject(projectPath) {
			runServer(projectPath)
		} else {
			showHelp(args)
		}
	}

}

// Show the version of this tool
func showVersion() {
	helpString := DIVIDER
	helpString += fmt.Sprintf("Fragmenta version: %s", VERSION)
	helpString += DIVIDER
	log.Print(helpString)
}

// Show the help for this tool.
func showHelp(args []string) {
	helpString := DIVIDER
	helpString += fmt.Sprintf("Fragmenta version: %s", VERSION)
	helpString += "\n  fragmenta version -> display version"
	helpString += "\n  fragmenta help -> display help"
	helpString += "\n  fragmenta new [path/to/app] -> creates a new app at the path supplied"
	helpString += "\n  fragmenta server -> runs server locally"
	helpString += "\n  fragmenta migrate -> runs new sql migrations in db/migrate"
	helpString += "\n  fragmenta generate resource [name] [fieldname]:[fieldtype]* -> creates resource CRUD actions and views"
	helpString += "\n  fragmenta generate migration [name] -> creates a new named sql migration in db/migrate"
	helpString += "\n  fragmenta test  -> run tests"
	helpString += "\n  fragmenta -> also runs server locally"
	helpString += DIVIDER
	log.Print(helpString)
}

// Run the server
func runServer(projectPath string) {
	showVersion()

	killServer()

	err := buildServer(localServerPath(projectPath), nil)

	if err != nil {
		log.Printf("Error building server: %s", err)
		return
	}

	log.Println("Building server...")
	buildAssets(projectPath)

	// We should at this point check if migrations need to be run
	// if they do, we should run them first
	// this is a bti slow so not doing it automatically... but that can lead to errors
	//	log.Println("Running migrations...")
	// runMigrate([]string{})

	log.Println("Launching server")

	cmd := exec.Command(localServerPath(projectPath))

	stdout, err := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()
	err = cmd.Start()
	if err != nil {
		log.Println(err)
	}
	go io.Copy(os.Stdout, stdout)
	go io.Copy(os.Stderr, stderr)
	cmd.Wait()

}

func killServer() {
	runCommand("killall", "-9", serverName())
}

func buildAssets(path string) {

}

func requireValidProject(projectPath string) bool {
	if isValidProject(projectPath) {
		return true
	} else {
		log.Printf("\nNo fragmenta project found at this path\n")
		return false
	}
}

func isValidProject(projectPath string) bool {

	_, err := os.Stat(appPath(projectPath))
	if err != nil {
		return false
	}

	_, err = os.Stat(configPath(projectPath))
	if err != nil {
		return false
	}

	return true
}

func runCommand(command string, args ...string) ([]byte, error) {

	cmd := exec.Command(command, args...)
	cmd.Stderr = os.Stdout
	//	cmd.Stderr = cmd.Stdout
	output, err := cmd.Output()
	if err != nil {
		return output, err
	} else {
		return output, nil
	}

}

// Read our config file and set up the server accordingly
func readConfig(projectPath string) {

	c := configPath(projectPath)

	// Read the config json file
	file, err := ioutil.ReadFile(c)
	if err != nil {
		log.Printf("Error opening config %s %v", c, err)
		return
	}

	var data map[string]map[string]string
	err = json.Unmarshal(file, &data)
	if err != nil {
		log.Printf("Error parsing config %s %v", c, err)
		return
	}

	ConfigDevelopment = data["development"]
	ConfigProduction = data["production"]
	ConfigTest = data["test"]
}

func cloneExamples(projectPath string) string {
	tmpDir := path.Join(os.TempDir(), ".fragmenta-examples")
	repo := "https://github.com/fragmenta/examples.git"

	// If templates already exists at tmpDir we remove it to avoid potential git conflicts
	// This means we clone each time this function is called...
	_, err := os.Stat(tmpDir)
	if err == nil {
		os.RemoveAll(tmpDir)
	}

	// Clone the examples repo
	result, err := runCommand("git", "clone", "--depth", "1", repo, tmpDir)
	if err != nil {
		log.Printf("Error calling git %s", err)
		return tmpDir
	}
	log.Printf("%s", string(result))

	return tmpDir
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
