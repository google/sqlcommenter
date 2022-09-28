# Sqlcommenter [In development]

SQLcommenter is a plugin/middleware/wrapper to augment application related information/tags with SQL Statements that can be used later to correlate user code with SQL statements.

## Installation

### Install from source

* Clone the source
* In terminal go inside the client folder location where we need to import google-sqlcommenter package and enter the below commands

```shell
go mod edit -replace google.com/sqlcommenter=path/to/google/sqlcommenter/go

go mod tiny
```
### Install from github [To be added]

## Usages

### go-sql-driver
Please use the sqlcommenter's default database driver to execute statements. \
Due to inherent nature of Go, the safer way to pass information from framework to database driver is via `context`. So, it is recommended to use the context based methods of `DB` interface like `QueryContext`, `ExecContext` and `PrepareContext`. 

```go
db, err := sqlcommenter.Open("<driver>", "<connectionString>", sqlcommenter.CommenterOptions{<tag>:<bool>})
```

#### Configuration

Users are given control over what tags they want to append by using `sqlcommenter.CommenterOptions` struct.

```go
type CommenterOptions struct {
	EnableDBDriver    bool
	EnableTraceparent bool  // OpenTelemetry trace information
	EnableRoute       bool  // applicable for web frameworks
	EnableFramework   bool  // applicable for web frameworks
	EnableController  bool  // applicable for web frameworks
	EnableAction      bool  // applicable for web frameworks
	}
```

### net/http
Populate the request context with sqlcommenter.AddHttpRouterTags(r) function in a custom middleware.

#### Note
* We only support the `database/sql` driver and have provided an implementation for that.
* <b>ORM related tags are added to the driver only when the tags are enabled in the commenter's driver's config and also the request context should passed to the querying functions</b>

#### Example
```go
// middleware is used to intercept incoming HTTP calls and populate request context with commenter tags.
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := sqlcommenter.AddHttpRouterTags(r, next)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
```

## Options

With Go SqlCommenter, we have configuration to choose which tags to be appended to the comment.

| Options         | Included by default? | go-sql-orm                                                                                                                                                                   | net/http                                                                                                                                                                     | Notes |
| --------------- | :------------------: | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :---: |
| `DBDriver`      |                      | [ go-sql-driver](https://pkg.go.dev/database/sql/driver)                                                                                                                     |                                                                                                                                                                              |
| `Action`        |                      |                                                                                                                                                                              | [net/http handle](https://pkg.go.dev/net/http#Handle)                                                                                                                        |       |
| `Route`         |                      |                                                                                                                                                                              | [net/http routing path](https://pkg.go.dev/github.com/gorilla/mux#Route.URLPath)                                                                                             |       |
| `Framework`     |                      |                                                                                                                                                                              | [net/http](https://pkg.go.dev/net/http)                                                                                                                                      |       |
| `Opentelemetry` |                      | [W3C TraceContext.Traceparent](https://www.w3.org/TR/trace-context/#traceparent-field), [W3C TraceContext.Tracestate](https://www.w3.org/TR/trace-context/#tracestate-field) | [W3C TraceContext.Traceparent](https://www.w3.org/TR/trace-context/#traceparent-field), [W3C TraceContext.Tracestate](https://www.w3.org/TR/trace-context/#tracestate-field) |       |
