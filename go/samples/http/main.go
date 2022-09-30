package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/sqlcommenter/go/core"
	gosql "github.com/google/sqlcommenter/go/database/sql"
	httpnet "github.com/google/sqlcommenter/go/net/http"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func Index(w http.ResponseWriter, r *http.Request) {

	exp, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	bsp := sdktrace.NewSimpleSpanProcessor(exp) // You should use batch span processor in prod
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(bsp),
	)

	ctx, span := tp.Tracer("foo").Start(r.Context(), "parent-span-name")
	defer span.End()

	db, err := gosql.Open("mysql", "root:root@/gotest", core.CommenterOptions{EnableDBDriver: true, EnableRoute: true, EnableAction: true, EnableFramework: true, EnableTraceparent: true})
	if err != nil {
		fmt.Println(err)
	} else {
		db.ExecContext(ctx, "Select 11;")
		db.Exec("Select 2;")
		db.Prepare("Select 10")
		db.PrepareContext(ctx, "Select 10")
	}
	fmt.Fprintf(w, "Hello World!")
}

// middleware is used to intercept incoming HTTP calls and apply general functions upon them.
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := core.ContextInject(r.Context(), httpnet.NewHTTPRequestExtractor(r, next))
		log.Printf("HTTP request sent to %s from %v", r.URL.Path, next)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	mux := http.NewServeMux()
	finalHandler := http.HandlerFunc(Index)
	mux.Handle("/", middleware((finalHandler)))
	log.Fatal(http.ListenAndServe(":8080", mux))
}
