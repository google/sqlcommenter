---
title: "SQLAlchemy"
date: 2019-05-29T08:40:11-06:00
draft: false
logo: /images/sqlalchemy-logo.png
weight: 3
---

- [Introduction](#introduction)
- [Requirements](#requirements)
- [BeforeExecuteFactory](#BeforeExecuteFactory)
- [Fields](#fields)
- [End to end examples](#end-to-end-examples)
- [With flask](#with-flask)
- [References](#references)

### Introduction
sqlcommenter-sqlalchemy provides a factory to create `before_cursor_execute`, called `BeforeExecuteFactory`

We provide a `BeforeExecuteFactory` that takes options such as
```python
CommenterCursorFactory(with_opencensus=<True or False>)
```

We provide options such as `with_opencensus` because
{{% notice warning%}}
Since OpenCensus [`trace_id`](https://opencensus.io/tracing/span/traceid) and [`span_id`](https://opencensus.io/tracing/span/spanid/) are highly ephemeral, including them in SQL comments will likely break any form of statement-based caching that doesn't strip out comments.
{{% /notice %}}

### Requirements

* Python X: any version of Python that is supported by SQLAlchemy
* [OpenCensus](https://opencensus.io/) optionally

### Installation

{{<tabs Pip Source>}}
{{<highlight shell>}}
pip3 install google-cloud-sqlcommenter
{{</highlight>}}

{{<highlight shell>}}
git clone https://github.com/google/sqlcommenter.git
cd python/sqlcommenter-python && python3 setup.py install
{{</highlight>}}
{{</tabs>}}

and then we shall perform the following imports in our source code:


### BeforeExecuteFactory

`BeforeExecuteFactory` is a factory that creates a `before_cursor_execute` hook to your engine to grab information about your application and augment it as a comment to your SQL statement.

```python
from sqlalchemy import create_engine, event
from sqlcommenter.sqlalchemy.executor import BeforeExecuteFactory

engine = create_engine(...) # Create the engine with your dialect of SQL
event.listen(engine, 'before_cursor_execute', BeforeExecuteFactory(), retval=True)
engine.execute(...) # comment will be appended to SQL before execution
```

**NOTE**
Please ensure that you set `retval=True` when listening for events

and this will produce such output on for example a Postgresql database logs:
```shell
2019-06-30 18:01:16.315 PDT [96973] LOG:  statement: SELECT * FROM polls_question
/*traceparent='00-ade4c36dc5e43b503a5bba237ea11746-578a74a562044332-01'*/
```

#### <a name="with-opencensus"></a> with_openCensus=True

To enable the comment cursor to also attach information about the current OpenCensus span (if any exists), pass in option `with_opencensus=True` when invoking `BeforeExecuteFactory`, so


```python
engine = create_engine("postgresql://:$postgres$@127.0.0.1:5432/quickstart_py")
event.listen(engine, 'before_cursor_execute', BeforeExecuteFactory(with_opencensus=True), retval=True)
engine.execute(...) # comment will be appended to SQL before execution
```

**NOTE**
Please ensure that you set `retval=True` when listening for events


### Fields

Field|Description|Included by default
---|---|---
`db_driver`|The underlying database driver e.g. `'psycopg2'`|<div style="text-align: center">&#10060;</div>
`db_framework`|The version of SQLAlchemy in the form `'sqlalchemy:<sqlalchemy_version>'`|<div style="text-align: center">&#10060;</div>
`traceparent`|The [W3C TraceContext.Traceparent field](https://www.w3.org/TR/trace-context/#traceparent-field) of the OpenCensus trace -- optionally defined with [`with_opencensus=True`](#with-opencensus)|<div style="text-align: center">&#10060;</div>
`tracestate`|The [W3C TraceContext.Tracestate field](https://www.w3.org/TR/trace-context/#tracestate-field) of the OpenCensus trace -- optionally defined with [`with_opencensus=True`](#with-opencensus)|<div style="text-align: center">&#10060;</div>

### End to end examples

#### Source code

{{<tabs "With OpenCensus" "With DB Framework" "With DB Driver">}}

{{<highlight python>}}
#!/usr/bin/env python3

from sqlalchemy import create_engine, event
from google.cloud.sqlcommenter.sqlalchemy.executor import BeforeExecuteFactory

def main():
    tracer = Tracer(exporter=NoopExporter, sampler=AlwaysOnSampler())
    engine = create_engine(DB_URL)

    listener = BeforeExecuteFactory(with_opencensus=True)
    event.listen(engine, 'before_cursor_execute', listener, retval=True)

    with tracer.span():
        result = engine.execute('SELECT * FROM polls_question')
        for row in result:
            print(row)

if __name__ == '__main__':
    main()
{{</highlight>}}

{{<highlight python>}}
#!/usr/bin/env python3

from sqlalchemy import create_engine, event
from google.cloud.sqlcommenter.sqlalchemy.executor import BeforeExecuteFactory

DB_URL = '...'  # DB connection info

def main():
    engine = create_engine(DB_URL)

    listener = BeforeExecuteFactory(with_db_framework=True)
    event.listen(engine, 'before_cursor_execute', listener, retval=True)

    result = engine.execute('SELECT * FROM polls_question')
    for row in result:
        print(row)

if __name__ == '__main__':
    main()
{{</highlight>}}

{{<highlight python>}}
#!/usr/bin/env python3

from sqlalchemy import create_engine, event
from google.cloud.sqlcommenter.sqlalchemy.executor import BeforeExecuteFactory

DB_URL = '...'  # DB connection info

def main():
    engine = create_engine(DB_URL)

    listener = BeforeExecuteFactory(with_db_driver=True)
    event.listen(engine, 'before_cursor_execute', listener, retval=True)

    result = engine.execute('SELECT * FROM polls_question')
    for row in result:
        print(row)

if __name__ == '__main__':
    main()
{{</highlight>}}

{{</tabs>}}

```shell
python3 main.py
(1, 'Wassup?', datetime.datetime(2019, 5, 30, 13, 51, 12, 910545, tzinfo=psycopg2.tz.FixedOffsetTimezone(offset=-420, name=None)))
(2, 'Wassup?', datetime.datetime(2019, 5, 30, 13, 57, 45, 905771, tzinfo=psycopg2.tz.FixedOffsetTimezone(offset=-420, name=None)))
(3, 'Wassup?', datetime.datetime(2019, 5, 30, 13, 57, 46, 908185, tzinfo=psycopg2.tz.FixedOffsetTimezone(offset=-420, name=None)))
(4, 'Wassup?', datetime.datetime(2019, 5, 30, 13, 57, 47, 557196, tzinfo=psycopg2.tz.FixedOffsetTimezone(offset=-420, name=None)))
(5, 'Wassup?', datetime.datetime(2019, 5, 30, 13, 57, 47, 853424, tzinfo=psycopg2.tz.FixedOffsetTimezone(offset=-420, name=None)))
```

#### Results

Examining our Postgresql server logs

{{<tabs "With OpenCensus" "With DB Framework" "With DB Driver">}}

{{<highlight shell>}}
2019-07-18 14:10:15.228 -03 [30071] LOG:  statement: SELECT * FROM polls_question
/*traceparent='00-bf66750ad4c76f614c0a99d843758cbb-e6b27c3caf35de73-01'*/
{{</highlight>}}

{{<highlight shell>}}
2019-07-18 14:11:19.576 -03 [30108] LOG:  statement: SELECT * FROM polls_question
/*db_framework='sqlalchemy%3A1.3.5'*/
{{</highlight>}}

{{<highlight shell>}}
2019-07-18 14:03:33.426 -03 [29858] LOG:  statement: SELECT * FROM polls_question
/*db_driver='psycopg2'*/
{{</highlight>}}

{{</tabs>}}

### With flask
When coupled with the web framework [flask](http://flask.pocoo.org), we still provide middleware to correlate
your web applications with your SQL statements from sqlalchemy. Please see this end-to-end guide below:
{{<card-vendor href="/python/flask#with-sqlalchemy" src="/images/flask-logo.png">}}

### References

Resource|URL
---|---
sqlcommenter-sqlalchemy on PyPi|https://pypi.org/project/google-cloud-sqlcommenter
sqlcommenter-sqlalchemy on Github|https://github.com/google/sqlcommenter
OpenCensus|https://opencensus.io/
OpenCensus SpanID|https://opencensus.io/tracing/span/spanid
OpenCensus TraceID|https://opencensus.io/tracing/span/traceid
