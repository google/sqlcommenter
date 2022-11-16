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
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return db
}
