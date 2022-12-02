---
title: "gorrila/mux"
date: 2022-12-02T13:33:00+05:30
draft: false
weight: 1
tags: ["go", "gorrila", "mux", "http"]
---

This library provides a middleware that extracts SQLCommenter HTTP request tags from a request being handled by [`gorrila/mux`](https://github.com/gorilla/mux) and attaches them to the request's context. This same context, when used to run queries using [sqlcommenter/go/database/sql](/go/database_sql), allows request tags and traceparent(if using [otelmux](https://github.com/open-telemetry/opentelemetry-go-contrib/tree/main/instrumentation/github.com/gorilla/mux/otelmux)) to be passed into SQL comments.

## Installation

```shell
go get -u github.com/google/sqlcommenter/go/gorrila/mux
```

## Usage

```go
import (
    "net/http"

    sqlcommentermux "github.com/google/sqlcommenter/go/gorrila/mux"
    "github.com/gorilla/mux"
)

func runApp() {
    r := mux.NewRouter()
    r.Use(sqlcommentermux.SQLCommenterMiddleware)

    r.HandleFunc("/", ActionHome).Methods("GET")

    http.ListenAndServe(":8081", r)
}
```

## Usage (with otelmux)

```go
import (
    "context"
    "log"
    "net/http"

    sqlcommentermux "github.com/google/sqlcommenter/go/gorrila/mux"
    "github.com/gorilla/mux"
    "go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
    "go.opentelemetry.io/otel"
    stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
    "go.opentelemetry.io/otel/propagation"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func main() {
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
    r.Use(otelmux.Middleware("sqlcommenter sample-server"), sqlcommentermux.SQLCommenterMiddleware)

    r.HandleFunc("/", ActionHome).Methods("GET")

    http.ListenAndServe(":8081", r)
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
```

## Tags

| Options         | Included by default? | gorrila/mux                                                                                                                                                                     |
| --------------- | -------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Action`        |          Y            | name of the handler function                                                                                                                        |
| `Route`         |          Y            | [routing path](https://pkg.go.dev/github.com/gorilla/mux#Route.GetPathTemplate)                                                                                             |
| `Framework`     |          Y            | [gorrila/mux](https://github.com/gorilla/mux)                                                                                                                                      |
| `Opentelemetry` |          N            | [W3C TraceContext.Traceparent](https://www.w3.org/TR/trace-context/#traceparent-field), [W3C TraceContext.Tracestate](https://www.w3.org/TR/trace-context/#tracestate-field) |