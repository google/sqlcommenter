# http-net  [In development]

SQLcommenter is a plugin/middleware/wrapper to augment application related information/tags with SQL Statements that can be used later to correlate user code with SQL statements.

## Installation

### Install from source

* Clone the source
* In terminal go inside the client folder location where we need to import sqlcommenter http/net module and enter the below commands

```shell
go mod edit -replace google.com/sqlcommenter=path/to/google/sqlcommenter/http-net

go mod tiny

go get google.com/sqlcommenter/http-net
```
### Install from github [To be added]

## Usage

Please use the sqlcommenter's default go-sql database driver to execute statements. \
Due to inherent nature of Go, the safer way to pass information from framework to database driver is via `context`. So, it is recommended to use the context based methods of `DB` interface like `QueryContext`, `ExecContext` and `PrepareContext`. 

```go
db, err := gosql.Open("<driver>", "<connectionString>", sqlcommenter.CommenterOptions{<tag>:<bool>})
```

### Configuration

Users are given control over what tags they want to append by using `core.CommenterOptions` struct.

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


#### Note
* We only support the `database/sql` driver and have provided an implementation for that.
* <b>ORM related tags are added to the driver only when the tags are enabled in the commenter's driver's config and also the request context should passed to the querying functions</b>
* <b>The middleware implementing this sqlcommenter http-net module should be added at the last</b>
  
#### Example
```go
// middleware is used to intercept incoming HTTP calls and populate request context with commenter tags.
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	ctx := httpnet.NewHttpNet(r, next).AddTags(r.Context())
	next.ServeHTTP(w, r.WithContext(ctx))
	})
}
```

## Options

With Go SqlCommenter, we have configuration to choose which tags to be appended to the comment.

| Options         | Included by default? | net/http                                                                                                                                                                     |
| --------------- | -------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Action`        |                      | [net/http handle](https://pkg.go.dev/net/http#Handle)                                                                                                                        |
| `Route`         |                      | [net/http routing path](https://pkg.go.dev/github.com/gorilla/mux#Route.URLPath)                                                                                             |
| `Framework`     |                      | [net/http](https://pkg.go.dev/net/http)                                                                                                                                      |
| `Opentelemetry` |                      | [W3C TraceContext.Traceparent](https://www.w3.org/TR/trace-context/#traceparent-field), [W3C TraceContext.Tracestate](https://www.w3.org/TR/trace-context/#tracestate-field) |
