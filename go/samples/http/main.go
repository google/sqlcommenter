package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/google/sqlcommenter/go/core"
	gosql "github.com/google/sqlcommenter/go/database/sql"
	httpnet "github.com/google/sqlcommenter/go/net/http"
	"github.com/julienschmidt/httprouter"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"sqlcommenter-http/mysqldb"
	"sqlcommenter-http/pgdb"
	"sqlcommenter-http/todos"
)

func MakeIndexRoute(db *gosql.DB) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
}

// middleware is used to intercept incoming HTTP calls and apply general functions upon them.
func middleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := core.ContextInject(r.Context(), httpnet.NewHTTPRequestExtractor(r, next))
		log.Printf("HTTP request sent to %s from %v", r.URL.Path, next)
		next(w, r.WithContext(ctx), p)
	}
}

func runApp(todosController *todos.TodosController) {
	err := todosController.CreateTodosTableIfNotExists()
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()

	index := MakeIndexRoute(todosController.DB)
	router.GET("/", middleware(index))

	router.GET("/todos", middleware(todosController.ActionList))
	router.POST("/todos", middleware(todosController.ActionInsert))
	router.PUT("/todos/:id", middleware(todosController.ActionUpdate))
	router.DELETE("/todos/:id", middleware(todosController.ActionDelete))

	http.ListenAndServe(":8080", router)
}

func runForMysql() *gosql.DB {
	connection := "root:password@/sqlcommenter_db"
	db := mysqldb.ConnectMySQL(connection)
	todosController := &todos.TodosController{DB: db, SQL: todos.MySQLQueries{}}
	runApp(todosController)
	return db
}

func runForPg() *gosql.DB {
	connection := "postgres://dev:dev@localhost/sqlcommenter_db?sslmode=disable"
	db := pgdb.ConnectPG(connection)
	todosController := &todos.TodosController{DB: db, SQL: todos.PGQueries{}}
	runApp(todosController)
	return db
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
