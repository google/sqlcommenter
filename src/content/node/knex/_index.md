---
title: "Knex.js"
date: 2019-05-29T08:40:11-06:00
draft: false
weight: 1
tags: ["knex", "knex.js", "query-builder", "node", "node.js", "express", "express.js"]
---

![](/images/knex-logo.png)

- [Introduction](#introduction)
- [Requirements](#requirements)
- [Installation](#installation)
    - [Manually](#manually)
    - [Package manager](#package-manager)
- [Usage](#usage)
    - [Plain knex wrapper](#plain-knex-wrapper)
    - [Express middleware](#express-middleware)
- [Fields](#fields)
    - [Options](#options)
        - [`include` config](#include-config)
        - [`options` config](#options-config)
        - [Options examples](#options-examples)
- [End to end examples](#end-to-examples)
    - [Source code](#source-code)
    - [Results](#results)
- [References](#references)

#### Introduction

This package is in the form of `Knex.Client.prototype.query` wrapper whose purpose is to augment a SQL statement right before execution, with
information about the controller and user code to help correlate them with SQL statements emitted by Knex.js.

Besides plain knex.js wrapping, we also provide a wrapper for the following frameworks:

{{<card-vendor href="#express-middleware" src="/images/express_js-logo.png">}}

### Requirements

Name|Resource
---|---
Knex.js|https://knexjs.org/
Node.js|https://nodejs.org/

### Installation

We can add integration into our applications in the following ways:

#### Manually

Please read [installing sqlcommenter-nodejs from source](/node/#install-from-source)

#### Package manager
Add to your package.json the dependency
{{<highlight json>}}
{
    "@google-cloud/sqlcommenter-knex": "*"
}
{{</highlight>}}
and then run `npm install` to get the latest version or 
{{<highlight shell>}}
npm install @google-cloud/sqlcommenter-knex --save
{{</highlight>}}

### Usage
#### Plain knex wrapper
{{<highlight javascript>}}
const {wrapMainKnex} = require('@google-cloud/sqlcommenter-knex');
const Knex = require('knex');
wrapMainKnex(Knex);

// Now you can create the knex client.
const knex = Knex(options);
{{</highlight>}}

#### Express middleware
This wrapper/middleware can be used as is or better with [express.js](https://expressjs.com/)
{{<highlight javascript>}}
const {wrapMainKnexAsMiddleware} = require('@google-cloud/sqlcommenter-knex');
const Knex = require('knex');
const app = require('express')();

// This is the important step where we set the middleware.
app.use(wrapMainKnexAsMiddleware(Knex));

// Now you can create the knex client.
const knex = Knex(options);
{{</highlight>}}

### Fields

In the database server logs, the comment's fields are:

* comma separated key-value pairs e.g. `route='%5Epolls/%24'`
* values are SQL escaped i.e. `key='value'`
* URL-quoted except for the equals(`=`) sign e.g `route='%5Epolls/%24'`. so should be URL-unquoted

| Field             | Format                        | Description                                                                                          | Example                                                                 |
|-------------------|-------------------------------|------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------|
| `db_driver`       | `<database_driver>:<version>` | URL quoted name and version of the database driver                                                   | `db_driver='knex%3A0.16.5'`                                             |
| `route`           | `<the route used>`            | URL quoted route used to match the express.js controller                                             | `route='%5E%2Fpolls%2F`                                                 |
| `traceparent`     | `<traceparent header>`        | URL quoted [W3C `traceparent` header](https://www.w3.org/TR/trace-context/#traceparent-header)       | `traceparent='00-3e2914ebce6af09508dd1ff1128493a8-81d09ab4d8cde7cf-01'` |
| `tracestate`      | `<tracestate header>`         | URL quoted [W3C `tracestate` header](https://www.w3.org/TR/trace-context/#tracestate-header)         | `tracestate='rojo%253D00f067aa0ba902b7%2Ccongo%253Dt61rcWkgMzE'`        |

#### Options

When creating the middleware, one can optionally configure the injected
comments by passing in the `include` and `options` objects:

```javascript
wrapMainKnexAsMiddleware(Knex, include={...}, options={...});
```

##### `include` config
A map of values to be optionally included in the SQL comments.

| Field       | On by default |
|-------------|---------------|
| db_driver   | {{<uncheck>}} |
| route       | {{<check>}}   |
| traceparent | {{<uncheck>}} |
| tracestate  | {{<uncheck>}} |

##### `options` config
A configuration object specifying where to collect trace data from. Accepted
fields are: TraceProvider: Should be either `OpenCensus` or `OpenTelemetry`,
indicating which library to collect trace context from.

| Field         | Possible values                 |
|---------------|---------------------------------|
| TraceProvider | `OpenCensus` or `OpenTelemetry` |

##### Options examples

{{<tabs "trace attributes" route db_driver "all set">}}

{{<highlight javascript>}}
wrapMainKnexAsMiddleware(
    Knex,
    include={ traceparent: true, tracestate: true },
    options={ TraceProvider: 'OpenTelemetry' }
);
{{</highlight>}}

{{<highlight javascript>}}
wrapMainKnexAsMiddleware(Knex, include={route: true});
{{</highlight>}}

{{<highlight javascript>}}
wrapMainKnexAsMiddleware(Knex, include={db_driver: true});
{{</highlight>}}

{{<highlight javascript>}}
// Manually set all the variables.
wrapMainKnexAsMiddleware(
    Knex,
    include={
        db_driver: true,
        route: true,
        traceparent: true,
        tracestate: true,
    },
    options={ TraceProvider: 'OpenTelemetry' }
);
{{</highlight>}}

{{</tabs>}}

### End to end examples

Check out a full express + opentelemetry example
[here](https://github.com/google/sqlcommenter/tree/master/nodejs/sqlcommenter-nodejs/samples/express-opentelemetry).

#### Source code

{{<tabs "With OpenCensus" "With OpenTelemetry" "With Route" "With DB Driver" "With All Options Set">}}

{{<highlight javascript>}}
// In file app.js.
const tracing = require('@opencensus/nodejs');
const {B3Format} = require('@opencensus/propagation-b3');
const {ZipkinTraceExporter} = require('@opencensus/exporter-zipkin');
const Knex = require('knex'); // Knex to be wrapped say v0.0.1
const {wrapMainKnexAsMiddleware} = require('@google-cloud/sqlcommenter-knex');
const express = require('express');

const exporter = new ZipkinTraceExporter({
    url: process.env.ZIPKIN_TRACE_URL || 'localhost://9411/api/v2/spans',
    serviceName: 'trace-542'
});

const b3 = new B3Format();
const traceOptions = {
    samplingRate: 1, // Always sample
    propagation: b3,
    exporter: exporter
};

// start tracing
tracing.start(traceOptions);

const knexOptions = {
    client: 'postgresql',
    connection: {
        host: '127.0.0.1',
        password: '$postgres$',
        database: 'quickstart_nodejs'
    }
};
const knex = Knex(knexOptions); // knex instance

const app = express();
const port = process.env.APP_PORT || 3000;

// Use the knex+express middleware with trace attributes
app.use(wrapMainKnexAsMiddleware(Knex, {
    traceparent: true,
    tracestate: true,
    route: false
}));

app.get('/', (req, res) => res.send('Hello, sqlcommenter-nodejs!!'));
app.get('^/polls/:param', function(req, res) {
    knex.raw('SELECT * from polls_question').then(function(polls) {
        const blob = JSON.stringify(polls);
        res.send(blob);
    }).catch(function(err) {
        console.log(err);
        res.send(500);
    });
});
app.listen(port, () => console.log(`Application listening on ${port}`));
{{</highlight>}}

{{<highlight javascript>}}
// In file app.js.
const { NodeTracerProvider } = require("@opentelemetry/node");
const { BatchSpanProcessor } = require("@opentelemetry/tracing");
const {
  TraceExporter,
} = require("@google-cloud/opentelemetry-cloud-trace-exporter");

const tracerProvider = new NodeTracerProvider();
// Export to Google Cloud Trace
tracerProvider.addSpanProcessor(
  new BatchSpanProcessor(new TraceExporter({ logger }), {
    bufferSize: 500,
    bufferTimeout: 5 * 1000,
  })
);
tracerProvider.register();

// OpenTelemetry initialization should happen before importing any libraries
// that it instruments
const express = require("express");
const Knex = require("knex");
const { wrapMainKnexAsMiddleware } = require("@google-cloud/sqlcommenter-knex");

const knexOptions = {
    client: 'postgresql',
    connection: {
        host: '127.0.0.1',
        password: '$postgres$',
        database: 'quickstart_nodejs'
    }
};
const knex = Knex(knexOptions); // knex instance

const app = express();
const port = process.env.APP_PORT || 3000;

// SQLCommenter express middleware injects the route into the traces
app.use(
  wrapMainKnexAsMiddleware(
    Knex,
    {
      traceparent: true,
      tracestate: true,

      // Optional
      db_driver: false,
      route: false,
    },
    { TraceProvider: "OpenTelemetry" }
  )
);

app.get('/', (req, res) => res.send('Hello, sqlcommenter-nodejs!!'));
app.get('^/polls/:param', function(req, res) {
    knex.raw('SELECT * from polls_question').then(function(polls) {
        const blob = JSON.stringify(polls);
        res.send(blob);
    }).catch(function(err) {
        console.log(err);
        res.send(500);
    });
});
app.listen(port, () => console.log(`Application listening on ${port}`));
{{</highlight>}}

{{<highlight javascript>}}
// In file app.js.
const Knex = require('knex'); // Knex to be wrapped say v0.0.1
const {wrapMainKnexAsMiddleware} = require('@google-cloud/sqlcommenter-knex');
const express = require('express');

const options = {
    client: 'postgresql',
    connection: {
        host: '127.0.0.1',
        password: '$postgres$',
        database: 'quickstart_nodejs'
    }
};
const knex = Knex(options); // knex instance

const app = express();
const port = process.env.APP_PORT || 3000;

// Use the knex+express middleware with route
app.use(wrapMainKnexAsMiddleware(Knex, {route: true}));

app.get('/', (req, res) => res.send('Hello, sqlcommenter-nodejs!!'));
app.get('^/polls/:param', function(req, res) {
    knex.raw('SELECT * from polls_question').then(function(polls) {
        const blob = JSON.stringify(polls);
        res.send(blob);
    }).catch(function(err) {
        console.log(err);
        res.send(500);
    });
});
app.listen(port, () => console.log(`Application listening on ${port}`));
{{</highlight>}}


{{<highlight javascript>}}
// In file app.js
const Knex = require('knex'); // Knex to be wrapped say v0.0.1
const {wrapMainKnexAsMiddleware} = require('@google-cloud/sqlcommenter-knex');
const express = require('express');

const options = {
    client: 'postgresql',
    connection: {
        host: '127.0.0.1',
        password: '$postgres$',
        database: 'quickstart_nodejs'
    }
};
const knex = Knex(options); // knex instance

const app = express();
const port = process.env.APP_PORT || 3000;

// Use the knex+express middleware with db driver
app.use(wrapMainKnexAsMiddleware(Knex, {db_driver: true}));

app.get('/', (req, res) => res.send('Hello, sqlcommenter-nodejs!!'));
app.get('^/polls/:param', function(req, res) {
    knex.raw('SELECT * from polls_question').then(function(polls) {
        const blob = JSON.stringify(polls);
        res.send(blob);
    }).catch(function(err) {
        console.log(err);
        res.send(500);
    });
});
app.listen(port, () => console.log(`Application listening on ${port}`));
{{</highlight>}}

{{<highlight javascript>}}
// In file app.js.
const tracing = require('@opencensus/nodejs');
const {B3Format} = require('@opencensus/propagation-b3');
const {ZipkinTraceExporter} = require('@opencensus/exporter-zipkin');
const Knex = require('knex'); // Knex to be wrapped say v0.0.1
const {wrapMainKnexAsMiddleware} = require('@google-cloud/sqlcommenter-knex');
const express = require('express');

const exporter = new ZipkinTraceExporter({
    url: process.env.ZIPKIN_TRACE_URL || 'localhost:9411/api/v2/spans',
    serviceName: 'trace-542'
});

const b3 = new B3Format();
const traceOptions = {
    samplingRate: 1, // Always sample
    propagation: b3,
    exporter: exporter
};

// start tracing
tracing.start(traceOptions);

const knexOptions = {
    client: 'postgresql',
    connection: {
        host: '127.0.0.1',
        password: '$postgres$',
        database: 'quickstart_nodejs'
    }
};
const knex = Knex(knexOptions); // knex instance

const app = express();
const port = process.env.APP_PORT || 3000;

// Use the knex+express middleware with all attributes set
app.use(wrapMainKnexAsMiddleware(Knex, {
    traceparent: true,
    tracestate: true,
    route: true,
    db_driver: true
}));

app.get('/', (req, res) => res.send('Hello, sqlcommenter-nodejs!!'));
app.get('^/polls/:param', function(req, res) {
    knex.raw('SELECT * from polls_question').then(function(polls) {
        const blob = JSON.stringify(polls);
        res.send(blob);
    }).catch(function(err) {
        console.log(err);
        res.send(500);
    });
});
app.listen(port, () => console.log(`Application listening on ${port}`));
{{</highlight>}}

{{</tabs>}}

which after running by
```shell
$ node app.js 
Application listening on 3000
```

#### Results

On making a request to that server at `http://localhost:3000/polls/1000`, the PostgreSQL logs show:

{{<tabs "With OpenCensus" "With Route" "With DB Driver" "With All Options Set">}}

{{<highlight shell>}}
2019-06-03 14:32:10.842 PDT [32004] LOG:  statement: SELECT * from polls_question 
/*traceparent='00-11000000000000ff-020000ee-01',tracestate='brazzaville=t61rcWkgMzE,rondo=00f067aa0ba902b7'*/
{{</highlight>}}

{{<highlight shell>}}
2019-06-03 14:32:10.842 PDT [32004] LOG:  statement: SELECT * from polls_question 
/*route='%5E%2Fpolls%2F%1000'*/
{{</highlight>}}

{{<highlight shell>}}
2019-06-03 14:32:10.842 PDT [32004] LOG:  statement: SELECT * from polls_question 
/*db_driver='knex%3A0.0.1'*/
{{</highlight>}}

{{<highlight shell>}}
2019-06-03 14:32:10.842 PDT [32004] LOG:  statement: SELECT * from polls_question 
/*db_driver='knex%3A0.0.1',route='%5E%2Fpolls%2F%1000',traceparent='00-11000000000000ff-020000ee-01',tracestate='brazzaville=t61rcWkgMzE,rondo=00f067aa0ba902b7'*/
{{</highlight>}}

{{</tabs>}}

#### References

| Resource                               | URL                                                           |
|----------------------------------------|---------------------------------------------------------------|
| @google-cloud/sqlcommenter-knex on npm | https://www.npmjs.com/package/@google-cloud/sqlcommenter-knex |
| express.js                             | https://expressjs.com/                                        |
