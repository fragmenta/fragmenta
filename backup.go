package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// RunBackup creates a backup of the chosen database
func RunBackup(args []string) {
	// Remove fragmenta backup from args list
	args = args[2:]

	switch fragmentaConfig(args) {
	case ModeProduction:
		backupDB(ConfigProduction)
	case ModeTest:
		backupDB(ConfigTest)
	default:
		backupDB(ConfigDevelopment)
	}
}

// RunRestore restores the chosen database from a backup
func RunRestore(args []string) {
	// Remove fragmenta backup from args list
	args = args[2:]

	switch fragmentaConfig(args) {
	case ModeProduction:
		restoreDB(ConfigProduction)
	case ModeTest:
		restoreDB(ConfigTest)
	default:
		restoreDB(ConfigDevelopment)
	}

	// Now that we have restored, run a post restore script if it exists
	restore := "./bin/restore"
	_, err := os.Stat(restore)
	if err == nil {
		log.Printf("Running restore script from " + restore)
		mode := fragmentaConfig(args)
		result, err := runCommand(restore, mode)
		if err != nil {
			log.Printf("Error running restore script %s", err)
			return
		}
		log.Printf("%s", result)
	}

}

// restoreDB restores using pg_dump from the latest backup
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

// backupDB backs up the db using the pg_dump binary
// FIXME - look into alternatives using sql or adapters
// we could have a backup/restore function on adapters?
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

// fragmentaConfig returns the config set by args (development by default)
func fragmentaConfig(args []string) string {
	if len(args) > 0 {
		return args[0]
	}

	return ModeDevelopment // default to dev config
}
