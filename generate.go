package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

const (
	permissions = 0744
)

// These variables are set from user input and then used in generation
var (
	resourceName string
	columns      map[string]string
)

// RunGenerate runs the generate command
// Expects:
// - generate migration
// - generate resource pages name:text summary:text
func RunGenerate(args []string) {
	// Remove fragmenta generate from args list
	args = args[2:]

	if len(args) < 2 {
		fmt.Println("Not enough arguments")
		return
	}
	command := args[0]
	args = args[1:]
	switch command {
	case "migration":
		name := args[0]
		sql := fmt.Sprintf("/* SQL migration %s */", name)
		generateMigration(name, sql)
	case "resource":
		generateResource(args)
	case "join":
		if len(args) < 2 {
			fmt.Println("Error - not enough arguments for join table")
			return
		}
		sort.Strings(args)
		name := fmt.Sprintf("%s-%s", args[0], args[1])
		sql := generateJoinSQL(args)
		generateMigration(name, sql)
	default:
		fmt.Println("Sorry, I didn't recognise that argument, you can use fragmenta generate [migration|resource|join]")
	}
}

// generateResource creates the scaffold for a new REST resource
// args should use snake_case, which is converted to camel case as necessary.
func generateResource(args []string) {

	// Read user input from args
	resourceName = ""
	columns = make(map[string]string, 0)

	var joins []string
	for _, v := range args {

		if len(resourceName) == 0 {
			resourceName = strings.ToLower(v)
		} else {
			parts := strings.Split(v, ":")
			if len(parts) == 2 {
				key := strings.ToLower(parts[0])
				value := strings.ToLower(parts[1])

				if key == "joins" {
					// We have a list of joins, potentially separated by ,
					joins = strings.Split(value, ",")
				} else {
					// Add a normal column
					columns[key] = value
				}

			} else {
				fmt.Printf("Invalid fields at: %s", v)
			}
		}

	}

	// NB we expect to start with a lower case singular
	fmt.Printf("Generating resource with\n - name:%s\n - attributes:%v\n", resourceName, columns)

	joinSQL := ""
	if len(joins) > 0 {
		for _, j := range joins {
			joinSQL += generateJoinSQL([]string{resourceName, j})
		}
	}

	// Copy our files first, if we fail, return the error
	err := generateResourceFiles()
	if err != nil {
		log.Printf("Error generating resource: %s", err)
		return
	}

	// Then db migration
	generateResourceMigration(joinSQL)

	// Then generate routes
	generateResourceRoutes()

}

// Generate the routes required and insert them into the routes.go file
func generateResourceRoutes() {

	// TODO - this routesTemplate should be a file
	routesTemplate := `
    r.Add("/[[.fragmenta_resources]]", [[.fragmenta_resource]]actions.HandleIndex)
    r.Add("/[[.fragmenta_resources]]/create", [[.fragmenta_resource]]actions.HandleCreateShow)
    r.Add("/[[.fragmenta_resources]]/create", [[.fragmenta_resource]]actions.HandleCreate).Post()
    r.Add("/[[.fragmenta_resources]]/{id:[0-9]+}/update", [[.fragmenta_resource]]actions.HandleUpdateShow)
    r.Add("/[[.fragmenta_resources]]/{id:[0-9]+}/update", [[.fragmenta_resource]]actions.HandleUpdate).Post()
    r.Add("/[[.fragmenta_resources]]/{id:[0-9]+}/destroy", [[.fragmenta_resource]]actions.HandleDestroy).Post()
    r.Add("/[[.fragmenta_resources]]/{id:[0-9]+}", [[.fragmenta_resource]]actions.HandleShow)`

	resourceRoutes := reifyString(routesTemplate)

	routesPath := appRoutesFilePath()
	data, err := ioutil.ReadFile(routesPath)
	if err != nil {
		fmt.Printf("#error Error reading routes at:%s :%s", routesPath, err)
		return
	}

	fmt.Println("Generating resource routes at: ", routesPath)

	routes := string(data)

	if strings.Contains(routes, ToPlural(resourceName)+"/actions") {
		fmt.Println("Routes already exist for resource: ", resourceName)
		return
	}

	routesStart := "// Resource Routes\n"
	routes = strings.Replace(routes, routesStart, routesStart+resourceRoutes+"\n", 1)

	resourceImport := reifyString("\t\"[[.fragmenta_app_path]]/[[.fragmenta_resources]]/actions\"\n")
	importStart := "\"github.com/fragmenta/router\"\n\n"
	routes = strings.Replace(routes, importStart, importStart+resourceImport, 1)

	err = ioutil.WriteFile(routesPath, []byte(routes), permissions)
	if err != nil {
		fmt.Println("Error writing routes file: ", routesPath)
		return
	}

	fmt.Println("Generated resource routes")
}

// Generate SQL for a join table migration
func generateJoinSQL(args []string) string {

	if len(args) < 2 {
		return ""
	}

	// Sort the table names
	sort.Strings(args)
	a := args[0]
	b := args[1]

	sql := `
DROP TABLE IF EXISTS [[.join_table]];
CREATE TABLE [[.join_table]] (
[[.a]]_id int NOT NULL,
[[.b]]_id int NOT NULL
);
`

	context := map[string]string{
		"join_table": ToPlural(a) + "_" + ToPlural(b), // e.g. places_tags
		"a":          a,                               // places
		"b":          b,                               // tags
	}

	return renderTemplate(sql, context)

}

// Generate a migration to create this resource table
func generateResourceMigration(joinsSQL string) {

	// We add the following fields to all resourceNames
	sql := `DROP TABLE IF EXISTS [[.fragmenta_resources]];
CREATE TABLE [[.fragmenta_resources]] (
id SERIAL NOT NULL,
created_at timestamp,
updated_at timestamp,
`

	for k, v := range columns {
		sql = sql + fmt.Sprintf("%s %s,\n", k, toSQLType(v))
	}

	sql = sql + ");\n"
	sql = strings.Replace(sql, ",\n)", "\n)", -1)

	sql += "ALTER TABLE [[.fragmenta_resources]] OWNER TO [[.fragmenta_db_user]];\n"

	sql = reifyString(sql)

	sql += joinsSQL

	name := fmt.Sprintf("Create-%s", ToCamel(resourceName))
	generateMigration(name, sql)

}

// Return the path of the routes.go file
func appRoutesFilePath() string {
	// Find the routes.go file, and add the routes at the start of setRoutes()
	// We expect a config option to be set on development
	// otherwise we default to ./src/app/routes.go
	routesPath := ConfigDevelopment["path_routes"]
	if len(routesPath) == 0 {
		routesPath = "src/app/routes.go"
	}

	return routesPath
}

func appGeneratePath() string {
	codePath := ConfigDevelopment["path_generate"]
	if len(codePath) == 0 {
		codePath = "src"
	}
	return codePath
}

// goPath returns the setting of env variable $GOPATH
// or $HOME/go if no $GOPATH is set.
func goPath() string {
	p := os.ExpandEnv("$GOPATH")
	if len(p) > 0 {
		return p
	}

	return filepath.Join(os.ExpandEnv("$HOME"), "go")
}

// fullAppPath returns an absolute path to the app.
func fullAppPath() string {
	// Golang expects all source under GOPATH/src
	return filepath.Join(goPath(), "src", appPath())
}

// appPath returns a relative path for the app from GOPATH/src
// If no path is set in secrets file, we default to
// a relative path starting from $GOPATH/src.
func appPath() string {
	p := ConfigDevelopment["path"]
	if p == "" {
		// If no path set in secrets file
		// default to a relative path starting from GOPATH/src
		pwd, err := filepath.Abs(".")
		if err != nil {
			return p
		}
		// Remove gopath prefix
		base := filepath.Join(goPath(), "src")
		base = base + string(filepath.Separator)
		p = strings.Replace(pwd, base, "", 1)
	}
	return p
}

func appServerName() string {
	return filepath.Base(appPath())
}

func appTemplatesPath() string {
	return filepath.Join(fullAppPath(), "src", "lib", "templates", "fragmenta_resources")
}

func generateResourceFiles() error {

	srcPath := appTemplatesPath()

	// Log an error and return if we can't find src path for templates
	_, err := os.Stat(srcPath)
	if err != nil {
		return fmt.Errorf("no template files at :%s", srcPath)
	}

	// For a destination, use the set path or default to ./src/xxx
	dstPath := filepath.Join(fullAppPath(), appGeneratePath(), ToPlural(resourceName))

	// Log our usage of templates
	log.Printf("Using templates at %s, saving to:%s\n", srcPath, dstPath)

	return copyAndReifyFiles(srcPath, dstPath)
}

func copyAndReifyFiles(srcPath string, dstPath string) error {
	var err error

	//log.Printf(" %s =>\n", srcPath)

	// Get info on the src
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		log.Fatal("Error statting src path ", srcPath)
		return err
	}

	// If this is a directory, copy every file within the src folder over to dst folder
	if srcInfo.IsDir() {

		err = filepath.Walk(srcPath, func(fileSrc string, info os.FileInfo, err error) error {
			fileDst := dstPath

			// split the srcPath on 'fragmenta_resources'
			// and use everything after that as the dst path
			srcParts := strings.Split(fileSrc, "/fragmenta_resources/")
			if len(srcParts) == 2 {
				fileDst = filepath.Join(dstPath, srcParts[1])
			}

			fileDst = reifyName(fileDst)

			// Do not operate on dot files
			if strings.HasPrefix(filepath.Base(fileSrc), ".") && filepath.Base(fileSrc) != ".keep" {
				return nil
			}

			// If this entry is a dir, just make sure it exists
			if info.IsDir() {
				os.MkdirAll(dstPath, permissions)
				return nil
			}

			// If this entry is a file, recurse and reify the file
			return copyAndReifyFiles(fileSrc, fileDst)

		})

		return nil
	}

	// If the file already exists, we should probably prompt the user as to whether they want to overwrite?
	// We shouldn't overwrite by default as here...

	// Print file destinations without prefix of time on log, to make them stand out
	log.Printf("=> %s\n", dstPath)

	// Read the file
	template, err := ioutil.ReadFile(srcPath)
	if err != nil {
		log.Fatal("Error reading file ", srcPath)
	}

	// Substitutions
	output := reifyString(string(template))

	// Make sure enclosing dir exists
	os.MkdirAll(filepath.Dir(dstPath), permissions)

	// Now write out again at same path
	err = ioutil.WriteFile(dstPath, []byte(output), permissions)
	if err != nil {
		log.Fatal("Error writing file ", dstPath)
	}

	return err

}

// Render a template to a string with a given context
func renderTemplate(tmpl string, context map[string]string) string {

	t := template.New("fields")
	t.Delims("[[", "]]")
	t, err := t.Parse(tmpl)
	if err != nil {
		log.Printf("Error creating fields template")
		return ""
	}

	var rendered bytes.Buffer
	err = t.Execute(&rendered, context)
	if err != nil {
		log.Printf("Error rendering fields template")
		return ""
	}

	return rendered.String()
}

// Generate golang assignments for our struct fields with validation.
func newFields() string {
	tmpl := "\t[[.fragmenta_resource]].[[.field_name]] = resource.Validate[[.validate_type]](cols[\"[[.col_name]]\"])\n"
	fields := ""
	for _, k := range sortedKeys(columns) {
		fieldContext := map[string]string{
			"fragmenta_resource": resourceName,
			"col_name":           k,
			"field_name":         ToCamel(k),
			"validate_type":      toValidateType(columns[k]),
		}

		fields += renderTemplate(tmpl, fieldContext)

	}
	return fields
}

// Generate golang struct fields for our columns
func structFields() string {
	tmpl := "\t[[.field_name]]\t\t[[.field_type]]\n"
	fields := ""
	for _, k := range sortedKeys(columns) {
		fieldContext := map[string]string{
			"fragmenta_resources": ToPlural(resourceName),
			"fragmenta_resource":  resourceName,
			"Fragmenta_Resources": ToCamel(ToPlural(resourceName)),
			"Fragmenta_Resource":  ToCamel(resourceName),
			"field_name":          ToCamel(k),
			"field_type":          toGoType(columns[k]),
		}

		fields += renderTemplate(tmpl, fieldContext)

	}
	return fields
}

// Generate show page fields for our columns
func showFields() string {
	tmpl := "\t<p>[[.field_name]]: {{ .[[.fragmenta_resource]].[[.field_name]] }}</p>\n"
	fields := ""

	for _, k := range sortedKeys(columns) {
		fieldContext := map[string]string{
			"fragmenta_resources": ToPlural(resourceName),
			"fragmenta_resource":  resourceName,
			"Fragmenta_Resources": ToCamel(ToPlural(resourceName)),
			"Fragmenta_Resource":  ToCamel(resourceName),
			"field_name":          ToCamel(k),
		}
		fields += renderTemplate(tmpl, fieldContext)
	}
	return fields
}

// Generate a columns list
func showcolumns() string {
	tmpl := "\"[[.col_name]]\","
	cols := ""

	for _, k := range sortedKeys(columns) {

		context := map[string]string{
			"col_name": k,
		}
		cols += renderTemplate(tmpl, context)
	}

	cols = strings.TrimRight(cols, ",")

	return cols
}

// Generate form fields for our columns
func formFields() string {

	fields := ""
	tmpl := `    {{ [[.method]] "[[.field_name]]" "[[.column_name]]" .[[.fragmenta_resource]].[[.field_name]] }}
`
	for _, k := range sortedKeys(columns) {

		// We add status as a special case menu
		if k == "status" {
			fields += fmt.Sprintf(`{{ select "Status" "status" .%s.Status .%s.StatusOptions }}`, resourceName, resourceName)
		} else {
			fieldContext := map[string]string{
				"fragmenta_resources": ToPlural(resourceName),
				"fragmenta_resource":  resourceName,
				"Fragmenta_Resources": ToCamel(ToPlural(resourceName)),
				"Fragmenta_Resource":  ToCamel(resourceName),
				"method":              "field",
				"column_name":         k,
				"field_name":          ToCamel(k),
				"resource_name":       ToCamel(k),
				"field_type":          toInputType(columns[k]),
			}

			fields += renderTemplate(tmpl, fieldContext)
		}

	}
	return fields
}

// Make this file name concrete by substituting values
func reifyName(name string) string {
	name = strings.Replace(name, ".go.tmpl", ".go", -1)   // go files
	name = strings.Replace(name, ".got.tmpl", ".got", -1) // template files
	name = strings.Replace(name, "fragmenta_resource", resourceName, -1)
	name = strings.Replace(name, "fragmenta_resources", ToPlural(resourceName), -1)
	return name
}

// Make this template string concrete by filling in values
func reifyString(tmpl string) string {
	context := map[string]string{
		"fragmenta_app_path":    filepath.Join(appPath(), appGeneratePath()),
		"fragmenta_resources":   ToPlural(resourceName),
		"fragmenta_resource":    resourceName,
		"Fragmenta_Resources":   ToCamel(ToPlural(resourceName)),
		"Fragmenta_Resource":    ToCamel(resourceName),
		"fragmenta_fields":      structFields(),
		"fragmenta_form_fields": formFields(),
		"fragmenta_show_fields": showFields(),
		"fragmenta_new_fields":  newFields(),
		"fragmenta_columns":     showcolumns(),
		"fragmenta_db":          ConfigDevelopment["db"],
		"fragmenta_db_user":     ConfigDevelopment["db_user"],
		"fragmenta_app_name":    appServerName(),
	}

	return renderTemplate(tmpl, context)
}

// Convert a user-defined type to a go type
func toValidateType(fieldType string) string {

	switch fieldType {
	case "text", "string", "char(255)":
		return "String"
	case "int", "integer", "bigint":
		return "Int"
	case "time", "datetime", "timestamp", "date":
		return "Time"
	case "float":
		return "Float"
	case "double":
		return "Float"
	}

	return fieldType
}

// Convert a user-defined type to a go type
func toGoType(fieldType string) string {

	switch fieldType {
	case "text", "string", "char(255)":
		return "string"
	case "int", "integer", "bigint":
		return "int64"
	case "time", "datetime", "timestamp", "date":
		return "time.Time"
	case "float":
		return "float"
	case "double":
		return "float64"
	}

	return fieldType
}

// Convert a user-defined type to an sql type
// this may vary with the database
func toSQLType(fieldType string) string {
	switch fieldType {
	case "text", "string", "char(255)":
		return "text"
	case "int", "int64", "integer", "bigint":
		return "integer"
	case "timestamp", "time", "datetime", "date":
		return "timestamp"
	case "float":
		return "real"
	case "double":
		return "double precision"
	default:
		return fieldType
	}

}

// Convert a user-defined type to an input type
func toInputType(fieldType string) string {
	switch fieldType {
	case "text", "string", "char(255)":
		return "textfield"
	case "int", "int64", "integer", "bigint", "float", "double":
		return "number"
	case "timestamp", "time", "datetime", "date":
		return "date"
	default:
		return fieldType
	}
}

// sortedKeys returns the string keys of a map[string] sorted
func sortedKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// ------------------------- MIGRATIONS  --------------

// Generate a migration file in db/migrate
func generateMigration(name string, content string) {
	path := migrationPath(".", name)

	fmt.Println("Generating migration: ", name)

	err := ioutil.WriteFile(path, []byte(content), 0744)
	if err != nil {
		fmt.Println("Error writing migration file: ", path)
		return
	}

	fmt.Println("Generated migration at: ", path)

}

// Generate a suitable path for a migration from the current date/time down to nanosecond
func migrationPath(path string, name string) string {
	now := time.Now()
	layout := "2006-01-02-150405"
	return fmt.Sprintf("%s/db/migrate/%s-%s.sql", path, now.Format(layout), name)
}
