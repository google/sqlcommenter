package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"net/http"

	gosqlmux "github.com/google/sqlcommenter/go/gorrila/mux"
	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"go.opentelemetry.io/otel"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"sqlcommenter-http/mysqldb"
	"sqlcommenter-http/pgdb"
	"sqlcommenter-http/todos"
)

func runApp(todosController *todos.TodosController) {
	err := todosController.CreateTodosTableIfNotExists()
	if err != nil {
		log.Fatal(err)
	}

	tp, err := initTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	r := mux.NewRouter()
	r.Use(otelmux.Middleware("sqlcommenter sample-server"), gosqlmux.SQLCommenterMiddleware)

	r.HandleFunc("/todos", todosController.ActionList).Methods("GET")
	r.HandleFunc("/todos", todosController.ActionInsert).Methods("POST")
	r.HandleFunc("/todos/{id}", todosController.ActionUpdate).Methods("PUT")
	r.HandleFunc("/todos/{id}", todosController.ActionDelete).Methods("DELETE")

	http.ListenAndServe(":8081", r)
}

func runForMysql() *sql.DB {
	connection := "root:password@tcp(mysql:3306)/sqlcommenter_db"
	db := mysqldb.ConnectMySQL(connection)
	todosController := &todos.TodosController{Engine: "mysql", DB: db, SQL: todos.MySQLQueries{}}
	runApp(todosController)
	return db
}

func runForPg() *sql.DB {
	connection := "host=postgres user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"
	db := pgdb.ConnectPG(connection)
	todosController := &todos.TodosController{Engine: "pg", DB: db, SQL: todos.PGQueries{}}
	runApp(todosController)
	return db
}

func initTracer() (*sdktrace.TracerProvider, error) {
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		return nil, err
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

func main() {
	var engine string

	flag.StringVar(&engine, "db_engine", "pg", "db-engine to run the sample application on")
	flag.Parse()

	if engine != "mysql" && engine != "pg" {
		log.Fatalf("invalid engine: %s", engine)
	}

	var db *sql.DB

	switch engine {
	case "mysql":
		db = runForMysql()
	case "pg":
		db = runForPg()
	}

	db.Close()
}
