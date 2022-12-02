---
title: "http-tags"
date: 2022-12-02T13:42:00+05:30
draft: false
weight: 1
tags: ["go", "http", "context", "request"]
---

This is a low-level package that can be used to prepare SQLCommenterTags out of an http request. The core package can then be used to inject these tags into a context.

## Installation

```bash
go get -u github.com/google/sqlcommenter/go/net/http
```

## Usage

```go
import (
    sqlcommenterhttp "github.com/google/sqlcommenter/go/net/http"
    "github.com/google/sqlcommenter/go/core"
)

requestTags := sqlcommenterhttp.NewHTTPRequestTags(framework string, route string, action string)
ctx := core.ContextInject(request.Context(), requestTags)
requestWithTags := request.WithContext(ctx)
```

This package can be used to instrument SQLCommenter for various frameworks.
