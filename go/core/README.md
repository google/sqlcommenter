# SQLCommenter Core  [In development]

SQLcommenter is a plugin/middleware/wrapper to augment application related information/tags with SQL Statements that can be used later to correlate user code with SQL statements.

This package contains configuration options, framework interface and support functions for all the sqlcommenter go modules

## Installation

This is a support package and will be installed indirectly by other go sqlcommenter packages

## Usages

### Configuration

Users are given control over what tags they want to append by using `core.CommenterConfig` struct.

```go
type CommenterConfig struct {
    EnableDBDriver    bool
    EnableRoute       bool
    EnableFramework   bool
    EnableController  bool
    EnableAction      bool
    EnableTraceparent bool
    EnableApplication bool
}
```

Users can also set the values for some static tags by using `core.StaticTags` struct.

```go
type StaticTags struct {
    Application string
    DriverName  string
}
```

These two structs together form the `core.CommenterOptions` struct, which is used by [sqlcommenter/go/database/sql](../database/sql/README.md).

```go
type CommenterOptions struct {
    Config CommenterConfig
    Tags   StaticTags
}
```