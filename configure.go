package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/dstroot/chi_api/database"
)

// Config contains the configuration from environment variables
type Config struct {
	Username string `env:"USERNAME,default=admin"`
	Password string `env:"PASSWORD,default=admin"`

	Debug      bool          `env:"DEBUG,default=true"`
	Port       string        `env:"PORT,default=8000"`
	Timeout    time.Duration `env:"TIMEOUT_HOURS,default=72h"`
	MaxRefresh time.Duration `env:"MAX_REFRESH_DAYS,default=30d"`
	Key        string        `env:"JWT_KEY,default=secretkey"`
	Realm      string        `env:"REALM,default=myrealm"`
	NoSSL      bool          `env:"NO_SSL,default=true"`
	RateLimit  int64         `env:"RATE_LIMIT,default=25"`
	ConnLimit  int           `env:"CONN_LIMIT,default=20"`

	SQL struct {
		Host     string `env:"MSSQL_HOST,default=localhost"`
		Port     string `env:"MSSQL_PORT,default=1433"`
		User     string `env:"MSSQL_USER,default=admin"`
		Password string `env:"MSSQL_PASSWORD,default=admin"`
		Database string `env:"MSSQL_DATABASE,default=test"`
	}
}

// check is a helper to streamline our error checks.
// Use it *only* when the program should halt. (unrecoverable)
func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
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
		return
	}
	database.DB.SetMaxIdleConns(100)

	if cfg.Debug {
		// The first actual connection to the underlying datastore will be
		// established lazily, when it's needed for the first time. If you want
		// to check right away that the database is available and accessible
		// (for example, check that you can establish a network connection and log
		// in), use db.Ping().
		err = database.DB.Ping()
		if err != nil {
			return
		}
		log.Printf("Connection: %s\n", connString)
		log.Printf("Database Connected!\n")
	}
	return
}
