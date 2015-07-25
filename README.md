# Fragmenta

Fragmenta is a command line tool for creating, managing and deploying golang web applications. 

### Using Fragmenta

* fragmenta new [path/to/app] -> creates a new app at the path supplied
* fragmenta server -> runs server locally
* fragmenta migrate -> runs new sql migrations in db/migrate
* fragmenta generate resource [name] [fieldname]:[fieldtype]* -> creates resource CRUD actions and views
* fragmenta generate migration [name] -> creates a new sql migration in db/migrate
* fragmenta test  -> run tests

You can create a new application, then run it in two steps. *At present Fragmenta is in beta, and the example cms is in not yet complete.*
