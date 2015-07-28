# Fragmenta

Fragmenta is a command line tool for creating, managing and deploying golang web applications. 

### Using Fragmenta

* fragmenta help -> display help
* fragmenta version -> display version
* fragmenta new [path/to/app] -> creates a new app at the path supplied
* fragmenta server -> runs server locally
* fragmenta -> also runs server locally
* fragmenta test  -> run tests
* fragmenta backup [development|production|test] -> backup the database to db/backup
* fragmenta restore [development|production|test] -> backup the database from latest file in db/backup
* fragmenta deploy [development|production|test] -> build and deploy using bin/deploy
* fragmenta migrate -> runs new sql migrations in db/migrate
* fragmenta generate resource [name] [fieldname]:[fieldtype]* -> creates resource CRUD actions and views
* fragmenta generate migration [name] -> creates a new named sql migration in db/migrate

*NB At present Fragmenta is in private beta, and the example cms is in not yet complete.*
