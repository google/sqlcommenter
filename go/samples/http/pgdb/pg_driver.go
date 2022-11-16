package pgdb

import (
	"log"

	"github.com/google/sqlcommenter/go/core"
	gosql "github.com/google/sqlcommenter/go/database/sql"
	_ "github.com/lib/pq"
)

func ConnectPG(connection string) *gosql.DB {
	db, err := gosql.Open("postgres", connection, core.CommenterOptions{EnableDBDriver: true, EnableRoute: true, EnableAction: true, EnableFramework: true, EnableTraceparent: true})
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
