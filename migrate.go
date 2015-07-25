package main

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/fragmenta/query"
)

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

	db := config["db"]
	migrationCount := 0

	openDatabase(config)

	files, err := filepath.Glob("./db/migrate/*.sql")
	if err != nil {
		log.Printf("Error running restore %s", err)
		return
	}

	// NB this check assumes the penultimate migration file exists
	migration := readMetadata(config)

	migrate := false
	for _, file := range files {
		filename := path.Base(file)

		if migrate {
			log.Printf("\n%s\nRunning migration %s", fragmentaDivider, filename)
			// Execute this sql file against the database
			// Create our psql command
			result, err := runCommand("psql", "-d", db, "-f", file)
			if err != nil {
				// If at any point we fail, log it and break
				log.Printf("ERROR loading sql migration %s", err)
				log.Printf("This and all future migrations cancelled %s", err)
				break
			}

			// Now store this as a completed migration
			migration = filename
			migrationCount++

			log.Printf("Completed migration %s\n%s\n%s", migration, string(result), fragmentaDivider)
		} else if filename == migration {
			// NB this check assumes the penultimate migration file exists
			migrate = true
		}

	}

	if migrationCount > 0 {
		writeMetadata(config, migration)
		log.Printf("Migrations complete up to migration %s on db %s\n\n", migration, db)
	} else {
		log.Printf("Database %s is up to date at migration %s\n\n", db, migration)
	}

}

// Open our database
func openDatabase(config map[string]string) {
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
		log.Printf("Database ERROR %s", err)
		os.Exit(9)
	}

	log.Printf("%s\n", fragmentaDivider)
	log.Printf("Opened database at %s for user %s", config["db"], config["db_user"])

}

// We should perhaps do this with the db driver instead
func readMetadata(config map[string]string) string {
	migration := ""

	sql := "select migration_version from fragmenta_metadata order by id desc limit 1;"

	rows, err := query.QuerySQL(sql)
	if err != nil {
		log.Printf("Database ERROR %s", err)
		return ""
	}

	// We expect just one row, with one column (count)
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&migration)
		if err != nil {
			log.Printf("Database ERROR %s", err)
			return ""
		}
	}

	return migration
}

// Update the database with a line recording what we have done
func writeMetadata(config map[string]string, migrationVersion string) {

	sql := "Insert into fragmenta_metadata(updated_at,fragmenta_version,migration_version,status) VALUES(NOW(),$1,$2,100);"

	result, err := query.ExecSQL(sql, fragmentaVersion, migrationVersion)
	if err != nil {
		log.Printf("Database ERROR %s %s", err, result)
	}

}
