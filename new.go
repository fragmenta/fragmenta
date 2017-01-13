package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	createDatabaseMigrationName = "Create-Database"
	createTablesMigrationName   = "Create-Tables"
)

// RunNew creates a new fragmenta project given the argument
// Usage: fragmenta new [app|cms|api| valid repo path e.g. github.com/fragmenta/fragmenta-cms]
func RunNew(args []string) {

	// Remove fragmenta backup from args list
	args = args[2:]

	// We expect two args left:
	if len(args) < 2 {
		log.Printf("Both a project path and a project type or URL are required to create a new site\n")
		return
	}

	repo := args[0]
	projectPath, err := filepath.Abs(args[1])
	if err != nil {
		log.Printf("Error expanding file path\n")
		return
	}

	if !strings.HasPrefix(projectPath, filepath.Join(os.Getenv("GOPATH"), "src")) {
		log.Printf("You must create your project in $GOPATH/src\n")
		return
	}

	if fileExists(projectPath) {
		log.Printf("A folder already exists at path %s\n", projectPath)
		return
	}

	switch repo {
	case "app":
		repo = "github.com/fragmenta/fragmenta-app"
	case "cms":
		repo = "github.com/fragmenta/fragmenta-cms"
	case "blog":
		repo = "github.com/fragmenta/fragmenta-blog"
	default:
		// TODO clean repo if it contains https or .git...
	}

	// Log fetching our files
	log.Printf("Fetching from url: %s\n", repo)

	// Go get the project url, to make sure it is up to date, should use -u
	_, err = runCommand("go", "get", repo)
	if err != nil {
		log.Printf("Error calling go get %s", err)
		return
	}

	// Copy the pristine new site over
	goProjectPath := filepath.Join(os.Getenv("GOPATH"), "src", repo)
	err = copyNewSite(goProjectPath, projectPath)
	if err != nil {
		log.Printf("Error copying project %s", err)
		return
	}

	// Generate config files
	err = generateConfig(projectPath)
	if err != nil {
		log.Printf("Error generating config %s", err)
		return
	}

	// Generate a migration AND run it
	err = generateCreateSQL(projectPath)
	if err != nil {
		log.Printf("Error generating migrations %s", err)
		return
	}

	// Output instructions to let them change setup first if they wish
	showNewSiteHelp(projectPath)

}

func copyNewSite(goProjectPath, projectPath string) error {

	// Check that the folders up to the path exist, if not create them
	err := os.MkdirAll(filepath.Dir(projectPath), permissions)
	if err != nil {
		log.Printf("The project path could not be created: %s", err)
		return err
	}

	// Now recursively copy over the files from the original repo to project path
	log.Printf("Creating files at: %s", projectPath)
	_, err = copyPath(goProjectPath, projectPath)
	if err != nil {
		return err
	}

	// Delete the .git folder at that path
	gitPath := filepath.Join(projectPath, ".git")
	log.Printf("Removing all at:%s", gitPath)
	err = os.RemoveAll(gitPath)
	if err != nil {
		return err
	}

	// Run git init to get a new git repo here
	log.Printf("Initialising new git repo at:%s", projectPath)
	_, err = runCommand("git", "init", projectPath)
	if err != nil {
		return err
	}

	// Now reifyNewSite
	log.Printf("Updating import paths to: %s", projectPathRelative(projectPath))
	return reifyNewSite(goProjectPath, projectPath)
}

// Copy a path to another one - at present this is unix only
// Unfortunately there is no simple facility for this in golang stdlib,
// so we use unix command (sorry windows!)
// FIXME - do not rely on unix commands and do this properly
func copyPath(src, dst string) ([]byte, error) {
	// Replace this with an os independent version using filepath.Walk
	return runCommand("cp", "-r", src, dst)
}

// reifyNewSite changes import refs within go files to the correct format
func reifyNewSite(goProjectPath, projectPath string) error {
	files, err := collectFiles(projectPath, []string{".go"})
	if err != nil {
		return err
	}

	// For each go file within project, make sure the refs are to the new site,
	// not to the template site
	relGoProjectPath := projectPathRelative(goProjectPath)
	relProjectPath := projectPathRelative(projectPath)
	for _, f := range files {
		// Load the file, if it contains refs to goprojectpath, replace them with relative project path imports
		data, err := ioutil.ReadFile(f)
		if err != nil {
			return err
		}

		// Substitutions - consider reifying instead if it is any more complex
		fileString := string(data)
		if strings.Contains(fileString, relGoProjectPath) {
			fileString = strings.Replace(fileString, relGoProjectPath, relProjectPath, -1)
		}

		err = ioutil.WriteFile(f, []byte(fileString), permissions)
		if err != nil {
			return err
		}

	}

	return nil
}

// the user should be prompted to:

func showNewSiteHelp(projectPath string) {
	helpString := fragmentaDivider
	helpString += "Congratulations, we've made a new website at " + projectPathRelative(projectPath)
	helpString += "\n  if you wish you can edit the database config at secrets/fragmenta.json and sql at db/migrate"
	helpString += "\n  To get started, run the following commands:"
	helpString += "\n  cd " + projectPath
	helpString += "\n  fragmenta migrate"
	helpString += "\n  fragmenta"
	helpString += fragmentaDivider + "\n"
	fmt.Print(helpString) // fmt to avoid time output
}

// generateCreateSQL generates an SQL migration file to create the database user and database referred to in config
func generateCreateSQL(projectPath string) error {

	// Set up a Create-Database migration, which comes first
	name := filepath.Base(projectPath)
	d := ConfigDevelopment["db"]
	u := ConfigDevelopment["db_user"]
	p := ConfigDevelopment["db_pass"]
	sql := fmt.Sprintf("/* Setup database for %s */\nCREATE USER \"%s\" WITH PASSWORD '%s';\nCREATE DATABASE \"%s\" WITH OWNER \"%s\";", name, u, p, d, u)

	// Generate a migration to create db with today's date
	file := migrationPath(projectPath, createDatabaseMigrationName)
	err := ioutil.WriteFile(file, []byte(sql), 0744)
	if err != nil {
		return err
	}

	// If we have a Create-Tables file, copy it out to a new migration with today's date
	createTablesPath := filepath.Join(projectPath, "db", "migrate", createTablesMigrationName+".sql.tmpl")
	if fileExists(createTablesPath) {
		sql, err := ioutil.ReadFile(createTablesPath)
		if err != nil {
			return err
		}

		// Now vivify the template, for now we just replace one key
		sqlString := reifyString(string(sql))

		file = migrationPath(projectPath, createTablesMigrationName)
		err = ioutil.WriteFile(file, []byte(sqlString), 0744)
		if err != nil {
			return err
		}
		// Remove the old file
		os.Remove(createTablesPath)

	} else {
		log.Printf("Error: No Tables found at:%s", createTablesPath)
	}

	return nil
}

func projectPathRelative(projectPath string) string {
	goSrc := os.Getenv("GOPATH") + "/src/"
	return strings.Replace(projectPath, goSrc, "", 1)
}

func generateConfig(projectPath string) error {
	configPath := configPath(projectPath)
	prefix := filepath.Base(projectPath)
	log.Printf("Generating new config at %s", configPath)

	ConfigProduction = map[string]string{}
	ConfigDevelopment = map[string]string{}
	ConfigTest = map[string]string{
		"port":            "3000",
		"log":             "log/test.log",
		"db_adapter":      "postgres",
		"db":              prefix + "_test",
		"db_user":         prefix + "_server",
		"db_pass":         randomKey(8),
		"assets_compiled": "no",
		"path":            projectPathRelative(projectPath),
		"hmac_key":        randomKey(32),
		"secret_key":      randomKey(32),
	}

	// Should we ask for db prefix when setting up?
	// hmm, in fact can we do this setup here at all!!
	for k, v := range ConfigTest {
		ConfigDevelopment[k] = v
		ConfigProduction[k] = v
	}
	ConfigDevelopment["db"] = prefix + "_development"
	ConfigDevelopment["log"] = "log/development.log"
	ConfigDevelopment["hmac_key"] = randomKey(32)
	ConfigDevelopment["secret_key"] = randomKey(32)

	ConfigProduction["db"] = prefix + "_production"
	ConfigProduction["log"] = "log/production.log"
	ConfigProduction["port"] = "80"
	ConfigProduction["assets_compiled"] = "yes"
	ConfigProduction["hmac_key"] = randomKey(32)
	ConfigProduction["secret_key"] = randomKey(32)

	configs := map[string]map[string]string{
		ModeProduction:  ConfigProduction,
		ModeDevelopment: ConfigDevelopment,
		ModeTest:        ConfigTest,
	}

	configJSON, err := json.MarshalIndent(configs, "", "\t")
	if err != nil {
		log.Printf("Error parsing config %s %v", configPath, err)
		return err
	}

	// Write the config json file
	err = ioutil.WriteFile(configPath, configJSON, permissions)
	if err != nil {
		log.Printf("Error writing config %s %v", configPath, err)
		return err
	}

	return nil
}

// Generate a random 32 byte key encoded in base64
func randomKey(l int64) string {
	k := make([]byte, l)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return ""
	}
	return hex.EncodeToString(k)
}

// Collect the files with these extensions under src
func collectFiles(dir string, extensions []string) ([]string, error) {

	files := []string{}

	err := filepath.Walk(dir, func(file string, info os.FileInfo, err error) error {
		// If we have an err pass it up
		if err != nil {
			return err
		}

		// Deal with files only
		if !info.IsDir() {
			// Check for go files
			name := filepath.Base(file)
			if !strings.HasPrefix(name, ".") && strings.HasSuffix(name, ".go") {
				files = append(files, file)
			}
		}

		return nil
	})

	if err != nil {
		return files, err
	}

	return files, nil

}
