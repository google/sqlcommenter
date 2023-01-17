# go-sql-driver

SQLcommenter is a plugin/middleware/wrapper to augment application related information/tags with SQL Statements that can be used later to correlate user code with SQL statements.

## Installation

```shell
go get -u github.com/google/sqlcommenter/go/database/sql
```

## Usage

Please use the sqlcommenter's default go-sql database driver to execute statements. 
Due to inherent nature of Go, the safer way to pass information from framework to database driver is via `context`. So, it is recommended to use the context based methods of `DB` interface like `QueryContext`, `ExecContext` and `PrepareContext`. 

```go
import (
    gosql "github.com/google/sqlcommenter/go/database/sql"
    sqlcommentercore "github.com/google/sqlcommenter/go/core"
    _ "github.com/lib/pq" // or any other database driver
)

db, err := gosql.Open("<driver>", "<connectionString>", sqlcommentercore.CommenterOptions{
    Config: sqlcommentercore.CommenterConfig{<flag>:bool}
    Tags  : sqlcommentercore.StaticTags{<tag>: string} // optional
})
```

### Configuration

Users are given control over what tags they want to append by using `core.CommenterOptions` struct.

```go
type CommenterOptions struct {
    EnableDBDriver    bool
    EnableTraceparent bool   // OpenTelemetry trace information
    EnableRoute       bool   // applicable for web frameworks
    EnableFramework   bool   // applicable for web frameworks
    EnableController  bool   // applicable for web frameworks
    EnableAction      bool   // applicable for web frameworks
    EnableApplication bool   // applicable for web frameworks
}
```

Users can also provide static tags they want to use by using `core.StaticTags` struct. These tags are optional. If not provided, the `Application` tag will get auto-populated from the user-project's module-name in `go.mod`. `DriverName` is set to the driver-name provided in `gosql.Open(driverName, ...)` call.

```go
type StaticTags struct {
    Application string
    DriverName  string
}
```

The driver will try to use the module-name from the project's `go.mod` as the application name if `EnableApplication` is `true` and no `Application` string is provided (works correctly for compiled go applications).


### Framework Supported
* [http/net](../../net/http/README.md) - basic support
* [gorrila/mux](../../gorrila//mux/README.md) - proper support


## Options

With Go SqlCommenter, we have configuration to choose which tags to be appended to the comment.

| Options    | Included by default? | go-sql-driver                                            |
| ---------- | -------------------- | -------------------------------------------------------- |
| `DBDriver` |                      | [ go-sql-driver](https://pkg.go.dev/database/sql/driver) |
