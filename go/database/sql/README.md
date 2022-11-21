# go-sql-driver  [In development]

SQLcommenter is a plugin/middleware/wrapper to augment application related information/tags with SQL Statements that can be used later to correlate user code with SQL statements.

## Installation

### Install from source

* Clone the source
* In terminal go inside the client folder location where we need to import sqlcommenter go-sql module and enter the below commands

```shell
go mod edit -replace google.com/sqlcommenter=path/to/google/sqlcommenter/go-sql

go mod tiny

go get google.com/sqlcommenter/gosql
```
### Install from github [To be added]

## Usage

Please use the sqlcommenter's default go-sql database driver to execute statements. 
Due to inherent nature of Go, the safer way to pass information from framework to database driver is via `context`. So, it is recommended to use the context based methods of `DB` interface like `QueryContext`, `ExecContext` and `PrepareContext`. 

```go
db, err := gosql.Open("<driver>", "<connectionString>", sqlcommenter.CommenterOptions{<tag>:<bool>})
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
	Application       string // user-provided application-name. optional
	}
```

The driver will try to use the module-name from the project's `go.mod` as the application name if `EnableApplication` is `true` and no `Application` string is provided (works correctly for compiled go applications).


### Framework Supported
* [http/net](.../../../http-net/README.md)


## Options

With Go SqlCommenter, we have configuration to choose which tags to be appended to the comment.

| Options    | Included by default? | go-sql-driver                                            |
| ---------- | -------------------- | -------------------------------------------------------- |
| `DBDriver` |                      | [ go-sql-driver](https://pkg.go.dev/database/sql/driver) |
