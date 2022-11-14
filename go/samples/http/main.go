package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/google/sqlcommenter/go/core"
	gosql "github.com/google/sqlcommenter/go/database/sql"
	httpnet "github.com/google/sqlcommenter/go/net/http"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"sqlcommenter-http/mysqldb"
	"sqlcommenter-http/pgdb"
)

var db *gosql.DB

func Index(w http.ResponseWriter, r *http.Request) {
	exp, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	bsp := sdktrace.NewSimpleSpanProcessor(exp) // You should use batch span processor in prod
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(bsp),
	)

	ctx, span := tp.Tracer("foo").Start(r.Context(), "parent-span-name")
	defer span.End()

	db.ExecContext(ctx, "Select 1")
	db.Exec("Select 2")

	stmt1, err := db.Prepare("Select 3")
	if err != nil {
		log.Fatal(err)
	}
	stmt1.QueryRow()

	stmt2, err := db.PrepareContext(ctx, "Select 4")
	if err != nil {
		log.Fatal(err)
	}
	stmt2.QueryRow()

	db.QueryContext(ctx, "Select 5")

	fmt.Fprintf(w, "Hello World!\r\n")
}

// middleware is used to intercept incoming HTTP calls and apply general functions upon them.
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := core.ContextInject(r.Context(), httpnet.NewHTTPRequestExtractor(r, next))
		log.Printf("HTTP request sent to %s from %v", r.URL.Path, next)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func runForMysql() *gosql.DB {
	connection := "root:password@/sqlcommenter_db"
	db = mysqldb.ConnectMySQL(connection)

	mux := http.NewServeMux()
	finalHandler := http.HandlerFunc(Index)
	mux.Handle("/", middleware((finalHandler)))
	log.Fatal(http.ListenAndServe(":8080", mux))
	return db
	return db
}

func runForPg() *gosql.DB {
	connection := "postgres://dev:dev@localhost/sqlcommenter_db?sslmode=disable"
	db = pgdb.ConnectPG(connection)

	mux := http.NewServeMux()
	finalHandler := http.HandlerFunc(Index)
	mux.Handle("/", middleware((finalHandler)))
	log.Fatal(http.ListenAndServe(":8080", mux))
	return db
}

func main() {
	var engine string

	flag.StringVar(&engine, "db_engine", "mysql", "db-engine to run the sample on")
	flag.Parse()

	if engine != "mysql" && engine != "pg" {
		log.Fatalf("invalid engine: %s", engine)
	}

	var db *gosql.DB

	switch engine {
	case "mysql":
		db = runForMysql()
	case "pg":
		db = runForPg()
	}

	db.Close()
}
