package main

import (
	"github.com/fragmenta/query"
	"log"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// We provide no facility to rollback at the moment, because rollbacks have all sorts of subtle issues and are not often useful IME.

// runMigrate runs a migration
func runMigrate(args []string) {

	// Remove fragmenta backup from args list
	args = args[2:]

	switch fragmentaConfig(args) {
	case "production":
		migrateDB(ConfigProduction)
	case "test":
		migrateDB(ConfigTest)
	default:
		migrateDB(ConfigDevelopment)
	}

}

// Find the last run migration, and run all those after it in order
// We use the fragmenta_metadata table to do this
func migrateDB(config map[string]string) {
	var latestMigration string
	var migrationCount int
	var migrate bool

	// Get a list of migration files
	files, err := filepath.Glob("./db/migrate/*.sql")
	if err != nil {
		log.Printf("Error running restore %s", err)
		return
	}

	// Sort the list alphabetically
	sort.Strings(files)

	// Try opening the db (db may not exist at this stage)
	err = openDatabase(config)
	if err != nil {
		// if no db, migrate first
		migrate = true
	} else {
		latestMigration = readMetadata()
	}

	for _, file := range files {
		filename := path.Base(file)
		if filename == latestMigration {
			migrate = true
		} else if migrate {
			log.Printf("Running migration %s", filename)

			args := []string{"-d", config["db"], "-f", file}
			if strings.Contains(filename, createDatabaseMigrationName) {
				args = []string{"-f", file}
			}

			// Execute this sql file against the database
			result, err := runCommand("psql", args...)
			if err != nil {
				// If at any point we fail, log it and break
				log.Printf("ERROR loading sql migration %s", err)
				log.Printf("This and all future migrations cancelled %s", err)
				break
			}

			migrationCount++
			latestMigration = filename
			log.Printf("Completed migration %s\n%s\n%s", filename, string(result), fragmentaDivider)
		}
	}

	if migrationCount > 0 {
		writeMetadata(config, latestMigration)
		log.Printf("Migrations complete up to migration %s on db %s\n\n", latestMigration, config["db"])
	} else {
		log.Printf("Database %s is up to date at migration %s\n\n", config["db"], latestMigration)
	}

}

// Open our database
func openDatabase(config map[string]string) error {
	// Open the database
	options := map[string]string{
		"adapter":  config["db_adapter"],
		"user":     config["db_user"],
		"password": config["db_pass"],
		"db":       config["db"],
		// "debug"     : "true",
	}

	err := query.OpenDatabase(options)
	if err != nil {
		return err
	}

	log.Printf("%s\n", fragmentaDivider)
	log.Printf("Opened database at %s for user %s", config["db"], config["db_user"])
	return nil
}

// We should perhaps do this with the db driver instead
func readMetadata() string {
	latestMigration := ""

	sql := "select migration_version from fragmenta_metadata order by id desc limit 1;"

	rows, err := query.QuerySQL(sql)
	if err != nil {
		log.Printf("Database ERROR %s", err)
		return ""
	}

	// We expect just one row, with one column (count)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&latestMigration)
		if err != nil {
			log.Printf("Database ERROR %s", err)
			return ""
		}
	}

	return latestMigration
}

// Update the database with a line recording what we have done
func writeMetadata(config map[string]string, migrationVersion string) {

	sql := "Insert into fragmenta_metadata(updated_at,fragmenta_version,migration_version,status) VALUES(NOW(),$1,$2,100);"

	result, err := query.ExecSQL(sql, fragmentaVersion, migrationVersion)
	if err != nil {
		log.Printf("Database ERROR %s %s", err, result)
	}

}
