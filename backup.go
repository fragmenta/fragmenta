package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func fragmentaConfig(args []string) string {
	if len(args) > 0 {
		return args[0]
	}

	return "development" // default to dev config
}

func runBackup(args []string) {
	// Remove fragmenta backup from args list
	args = args[2:]

	switch fragmentaConfig(args) {
	case "production":
		backupDB(ConfigProduction)
	case "test":
		backupDB(ConfigTest)
	default:
		backupDB(ConfigDevelopment)
	}
}

func runRestore(args []string) {
	// Remove fragmenta backup from args list
	args = args[2:]

	switch fragmentaConfig(args) {
	case "production":
		restoreDB(ConfigProduction)
	case "test":
		restoreDB(ConfigTest)
	default:
		restoreDB(ConfigDevelopment)
	}

}

// Restore back to our db from latestbackup
func restoreDB(config map[string]string) {
	// Just assume it is psql for now
	db := config["db"]

	if len(db) == 0 {
		log.Printf("Error running restore - no config")
		return
	}

	files, err := filepath.Glob("./db/backup/*.sql.gz")
	if err != nil {
		log.Printf("Error running restore - %s", err)
		return
	}

	if len(files) == 0 {
		log.Printf("Error running restore - no files")
		return
	}

	gz := files[len(files)-1:][0]
	sql := strings.Trim(gz, ".gz")

	// Delete the sql file when we exit
	defer os.Remove(sql)

	log.Printf("Running restore for %s with %s", db, gz)

	// Unzip the file
	result, err := runCommand("gzip", "-d", "-k", gz)
	if err != nil {
		log.Printf("Error running gz %s", err)
		return
	}
	log.Printf("%s", string(result))

	// Create our psql command
	result, err = runCommand("psql", "-d", db, "-f", sql)
	if err != nil {
		log.Printf("Error running psql %s", err)
		return
	}
	log.Printf("%s", string(result))

	log.Printf("Restore complete to db %s with %s", db, sql)
}

func backupDB(config map[string]string) {

	// Just assume it is psql for now
	adapter := "pg_dump"
	db := config["db"]

	if len(db) == 0 {
		log.Println("Error running backup - no config")
		return
	}

	log.Printf("Running backup for %s", db)

	date := time.Now().Format("2006-01-02-15-04")
	dst := fmt.Sprintf("./db/backup/%s.sql", date)

	// Create our psql command c for clean, f for file
	result, err := runCommand(adapter, "-c", "-f", dst, db)
	if err != nil {
		log.Printf("Error running psql %s", err)
		return
	}
	log.Printf("%s", string(result))

	// use compress/gzip instead?
	result, err = runCommand("gzip", dst)
	if err != nil {
		log.Printf("Error running gz %s", err)
		return
	}
	log.Printf(string(result))

	log.Printf("Backup complete of db %s", db)
}
