package constants

// Environment variables used to configure the application.
const (
	// DATA (used as DATADIR) specifies the data storage directory.
	DATA = "DataDir"
	// DBHOST (used as DB_HOST) specifies the database connection host.
	DBHOST = "DB_Host"
	// DBNAME (used as DB_NAME) specifies the database name.
	DBNAME = "DB_Name"
	// DBPASSWORD (used as DB_PASSWORD) specifies the database connection password.
	DBPASSWORD = "DB_Password"
	// DBPORT (used as DB_PORT) specifies the database connection port.
	DBPORT = "DB_Port"
	// DBUSER (used as DB_USER) specifies the database connection user.
	DBUSER = "DB_User"
	// VERBOSITY (used as VERBOSITY) specifies the logging verbosity.
	VERBOSITY = "Verbosity"

	// CONFIGPATH denotes the expected location of the Prismriver config file.
	CONFIGPATH = "/etc/prismriver/prismriver.yml"
)
