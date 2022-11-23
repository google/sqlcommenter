# SQLCommenter http-net  [In development]

SQLcommenter is a plugin/middleware/wrapper to augment application related information/tags with SQL Statements that can be used later to correlate user code with SQL statements.

## Installation

```bash
go get -u github.com/google/sqlcommenter/go/net/http
```

## Usage

This is a low-level package that can be used to prepare SQLCommeneterTags out of an http request. The core package can then be used to inject these tags into a context

```go
import (
    sqlcommenterhttp "github.com/google/sqlcommenter/go/net/http"
)

requestTags := sqlcommenterhttp.NewHTTPRequestTags(framework string, route string, action string)
ctx := core.ContextInject(request.Context(), requestTags)
requestWithTags := request.WithContext(ctx)
```

This package can be used to instrument SQLCommenter for various frameworks.