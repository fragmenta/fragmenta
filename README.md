# Fragmenta

Fragmenta is a command line tool for creating, managing and deploying golang web applications. It comes with a suite of libraries which making developing web apps easier, and aims to allow managing apps without making too many assumptions about which libraries they use or their internal structure. It takes care of generating CRUD actions, handling auth, routing and rendering and leaves you to concentrate the parts of your app which are unique. 

### Using Fragmenta

* fragmenta version -> display version
* fragmenta help -> display help
* fragmenta new [app|cms|blog|URL] path/to/app -> creates a new app from the repository at URL at the path supplied
* fragmenta -> builds and runs a fragmenta app
* fragmenta server -> builds and runs a fragmenta app
* fragmenta test  -> run tests
* fragmenta backup [development|production|test] -> backup the database to db/backup
* fragmenta restore [development|production|test] -> backup the database from latest file in db/backup
* fragmenta deploy [development|production|test] -> build and deploy using bin/deploy
* fragmenta migrate -> runs new sql migrations in db/migrate
* fragmenta generate resource [name] [fieldname]:[fieldtype]* -> creates resource CRUD actions and views
* fragmenta generate migration [name] -> creates a new named sql migration in db/migrate


### App structure

The default apps are laid out with the following structure:

* bin -> server binaries, and optional deploy script
* db -> database backups and migrations
* public -> files for serving publicly, including assets, uploaded files etc
* secrets -> config files, not usually checked in
* server.go -> your app entrypoint (required)
* src -> app source files - structure within this folder is up to you

The pkg layout within the app is up to you - defaults are provided but are not mandatory. Within src the default are arranged in packages by resource - the generator generates a new resource with the following structure:

* pages -> resource name
* * actions -> go actions (handling CRUD etc) for this resource
* * assets -> js,css, images for this resource
* * pages.go -> the resource model file
* * pages_test.go -> tests for this model
* * views -> views for this resource


### Libraries

The following independent packages are available for use with fragmenta apps (or other go web apps). 

* assets - an asset pipeline with minification and concatenation with asset fingerprinting
* auth - utilities for authentication and authorisation
* fragmenta - a command line tool for generating and developing websites
* model - a base model class for optional inclusion in models
* * file - a package for handling file uploads
* * validate - a package for handling field validation
* query - a query builder and result wrapper for mysql, psql and optionally sqlite
* router - a router which allows pattern matching, routes, redirects, filters and provides a handler interface
* server - a simple server based on http.listenandserve
* view - a library for rendering view templates using html/template
* * helpers - helpers for form fields, currencies etc

### Database setup 

At present fragmenta assumes usage of Postresql, while it can be used with other databases, this is the default. 

If you're setting up postgresql for the first time, you my find it simplest to simply create a user for yourself either as a superuser or with the specific privileges for local development as follows:

```sql
sudo su postgres
createuser $USER
psql
ALTER USER $USER WITH SUPERUSER;
ALTER USER $USER WITH CREATEROLE;
ALTER USER $USER WITH CREATEDB;
```
