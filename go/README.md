# Sqlcommenter [In development]

Sqlcommenter is a plugin/middleware/wrapper to augment SQL statements from go libraries with comments that can be used later to correlate user code with SQL statements.

## Installation

### Install from source

* Clone the source
* In terminal go inside the client folder location where we need to import google-sqlcommenter package and enter the below commands

```shell
go mod edit -replace google.com/sqlcommenterGo=path/to/google/sqlcommenter/go

go mod tiny
```
### Install from github [To be added]

## Usages

### go-sql-driver
Use the sqlcommenter's db driver to  perform query

```go
db, err := sqlcommenterGo.Open("<driver>", "<connectionString>", sqlcommenter.CommenterOptions{<tag>:<bool>})

```

#### Configuration

Tags to be appended can be configured by creating a variable of type sqlcommenter.CommenterOptions

```go
type CommenterOptions struct {
	DBDriver   bool
    Traceparent bool
	Route      bool //applicable for web frameworks
	Framework  bool //applicable for web frameworks
	Controller bool //applicable for web frameworks
	Action     bool //applicable for web frameworks
}
```
### go-http-router
Populate the request context with sqlcommenter.AddHttpRouterTags(r) function in a custom middleware.

#### Note
* <b>It needs to be used with drivers such as go-sql-orm </b>
* <b>ORM related tags are added to the driver only when the tags are enabled in the commenter's driver's config and also the request context should passed to the querying functions</b>

#### Example
```go
// middleware is used to intercept incoming HTTP calls and populate request context with commenter tags.
func middleware(n httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		ctx := sqlcommenter.AddHttpRouterTags(r)
		log.Printf("HTTP request sent to %s from %s", r.URL.Path, r.RemoteAddr)

		// call registered handler
		n(w, r.WithContext(ctx), ps)
	}
}
```



## Options

With Go SqlCommenter, we have configuration to choose which tags to be appended to the comment.

| Options         | Included by default? | go-sql-orm                                                                                                                                                                   | http-router                                                                                                                                                                  | Notes |
| --------------- | :------------------: | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :---: |
| `DBDriver`      |                      | [ go-sql-driver](https://pkg.go.dev/database/sql/driver)                                                                                                                     |                                                                                                                                                                              |
| `Action`        |                      |                                                                                                                                                                              | [http-router handle](https://pkg.go.dev/github.com/julienschmidt/httprouter#Handle)                                                                                          |       |
| `Route`         |                      |                                                                                                                                                                              | [http-router routing path](https://pkg.go.dev/github.com/julienschmidt/httprouter#Router)                                                                                    |       |
| `Framework`     |                      |                                                                                                                                                                              | [github.com/julienschmidt/httprouter](https://pkg.go.dev/github.com/julienschmidt/httprouter)                                                                                |       |
| `Opentelemetry` |                      | [W3C TraceContext.Traceparent](https://www.w3.org/TR/trace-context/#traceparent-field), [W3C TraceContext.Tracestate](https://www.w3.org/TR/trace-context/#tracestate-field) | [W3C TraceContext.Traceparent](https://www.w3.org/TR/trace-context/#traceparent-field), [W3C TraceContext.Tracestate](https://www.w3.org/TR/trace-context/#tracestate-field) |       |
