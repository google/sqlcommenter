---
title: "database/sql"
date: 2022-12-02T12:56:00+05:30
draft: false
weight: 1
tags: ["go", "sql"]
---

SQLCommenter provides drop-in replacement for [go's `database/sql`](https://pkg.go.dev/database/sql) library. You can start using the [SQLcommenter's `database/sql`](https://github.com/google/sqlcommenter/tree/master/go/database/sql) without changing your business logic code and use one of the framework instrumentation to proapagate tags and traceparent to the SQL queries.

## Installation

```shell
go get -u github.com/google/sqlcommenter/go/database/sql
```

## Usage

Please use the sqlcommenter's go-sql database driver to open a connection and use that connection to execute statements. 
Due to inherent nature of Go, the safer way to pass information from framework to database driver is via `context`. So, it is recommended to use the context based methods of `DB` interface like `QueryContext`, `ExecContext` and `PrepareContext`. 

```go
import (
    "database/sql"

    gosql "github.com/google/sqlcommenter/go/database/sql"
    sqlcommentercore "github.com/google/sqlcommenter/go/core"
    _ "github.com/lib/pq" // or any other database driver
)

var (
  db *sql.DB
  err error
)

db, err = gosql.Open("<driver>", "<connectionString>", sqlcommentercore.CommenterOptions{
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
