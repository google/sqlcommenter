package mysqldb

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/sqlcommenter/go/core"
	gosql "github.com/google/sqlcommenter/go/database/sql"
)

func ConnectMySQL(connection string) *gosql.DB {
	db, err := gosql.Open("mysql", connection, core.CommenterOptions{EnableDBDriver: true, EnableRoute: true, EnableAction: true, EnableFramework: true, EnableTraceparent: true})
	if err != nil {
		log.Fatalf("Failed to connect to MySQL(%q), error: %v", connection, err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping the database, error: %v", err)
	}

	return db
}
