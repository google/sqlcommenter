# SQLCommenter Core  [In development]

SQLcommenter is a plugin/middleware/wrapper to augment application related information/tags with SQL Statements that can be used later to correlate user code with SQL statements.

This package contains configuration options, framework interface and support functions for all the sqlcommenter go modules

## Installation

This is a support package and will be installed indirectly by other go sqlcommenter packages

## Usages

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


