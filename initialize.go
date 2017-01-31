package main

import (
	"database/sql"
	"encoding/json"
	"log"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/dstroot/chi_api/database"
	env "github.com/joeshaw/envdecode"
	_ "github.com/joho/godotenv/autoload"
	"github.com/pkg/errors"
)

var (
	cfg Config // global configuration
)

// Config contains the configuration from environment variables
type Config struct {
	Debug bool   `env:"DEBUG,default=true"`
	Port  string `env:"PORT,default=9102"`
	Site  struct {
		Intuit   string `env:"SITE_INTUIT,default=http://localhost:3001"`
		TaxSayer string `env:"SITE_TAXSLAYER,default=http://localhost:3002"`
	}
	SQL struct {
		Host     string `env:"MSSQL_HOST,default=localhost"`
		Port     string `env:"MSSQL_PORT,default=1433"`
		User     string `env:"MSSQL_USER,default=admin"`
		Password string `env:"MSSQL_PASSWORD,default=admin"`
		Database string `env:"MSSQL_DATABASE,default=test"`
	}
	GiactURL           string `env:"GIACT_URL,default=https://api.giact.com/"`
	GiactAuthIntuit    string `env:"GIACT_AUTH_INTUIT,default=Basic..."`
	GiactAuthTaxSlayer string `env:"GIACT_AUTH_TAXSLAYER,default=Basic..."`
}

// setupDatabase connects to our SQL Server
func setupDatabase() (err error) {

	connString := "server=" + cfg.SQL.Host +
		";port=" + cfg.SQL.Port +
		";user id=" + cfg.SQL.User +
		";password=" + cfg.SQL.Password +
		";database=" + cfg.SQL.Database +
		";connection timeout=60" + // in seconds (default is 30)
		";dial timeout=10" + // in seconds (default is 5)
		";keepAlive=10" // in seconds; 0 to disable (default is 0)

	// open connection to SQL Server
	database.DB, err = sql.Open("mssql", connString)
	if err != nil {
		return errors.Wrap(err, "error connecting to database")
	}
	database.DB.SetMaxIdleConns(100)

	if cfg.Debug {
		// The first actual connection to the underlying datastore will be
		// established lazily, when it's needed for the first time. If you want
		// to check right away that the database is available and accessible
		// (for example, check that you can establish a network connection and log
		// in), use database.DB.Ping().
		err = database.DB.Ping()
		if err != nil {
			log.Printf("Connection: %s\n", connString)
			return errors.Wrap(err, "error pinging database")
		}
	}
	return nil
}

// initialize our configuration from environment variables.
func initialize() error {

	// For development, github.com/joho/godotenv/autoload
	// loads env variables from .env file for you.

	// Read configuration from env variables
	err := env.Decode(&cfg)
	if err != nil {
		return errors.Wrap(err, "configuration decode failed")
	}

	// log configuration for debugging
	if cfg.Debug {
		prettyCfg, _ := json.MarshalIndent(cfg, "", "  ")
		log.Printf("Configuration: \n%v", string(prettyCfg))
	}

	err1 := setupDatabase()
	if err1 != nil {
		return errors.Wrap(err1, "database connection failed")
	}

	return nil
}
