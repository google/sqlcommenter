package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/google/sqlcommenter/go/core"
	gosql "github.com/google/sqlcommenter/go/database/sql"
	httpnet "github.com/google/sqlcommenter/go/net/http"
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

// middleware is used to intercept incoming HTTP calls and apply general functions upon them.
func middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := core.ContextInject(r.Context(), httpnet.NewHTTPRequestExtractor(r, h))
		log.Printf("HTTP request sent to %s", r.URL.Path)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

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
	r.Use(otelmux.Middleware("sqlcommenter sample-server"))

	r.HandleFunc("/todos", todosController.ActionList).Methods("GET")
	r.HandleFunc("/todos", todosController.ActionInsert).Methods("POST")
	r.HandleFunc("/todos/{id}", todosController.ActionUpdate).Methods("PUT")
	r.HandleFunc("/todos/{id}", todosController.ActionDelete).Methods("DELETE")

	http.ListenAndServe(":8081", middleware(r))
}

// host = “host.docker.internal”

func runForMysql() *gosql.DB {
	connection := "root:password@tcp(mysql:3306)/sqlcommenter_db"
	db := mysqldb.ConnectMySQL(connection)
	todosController := &todos.TodosController{Engine: "mysql", DB: db, SQL: todos.MySQLQueries{}}
	runApp(todosController)
	return db
}

func runForPg() *gosql.DB {
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

	flag.StringVar(&engine, "db_engine", "mysql", "db-engine to run the sample application on")
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
