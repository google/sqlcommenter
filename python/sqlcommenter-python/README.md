# sqlcommenter

Python modules for popular projects that add meta info to your SQL queries as comments.

 * [Django](#django)
 * [SQLAlchemy](#sqlalchemy)
 * [Psycopg2](#psycopg2)

## Local Install

```shell
pip3 install --user google-cloud-sqlcommenter
```

If you'd like to record some OpenCensus information as well, just install it:

```shell
pip3 install opencensus
```

## Usage

### Django

Add the provided Django middleware to your Django project's settings. All queries executed within the standard request→response cycle will have the SQL comment prepended to them.

```python
MIDDLEWARE = [
  'google.cloud.sqlcommenter.django.middleware.SqlCommenter',
  ...
]
```

which when viewed say on Postgresql logs, produces
```shell
2019-05-28 11:54:50.780 PDT [64128] LOG:  statement: INSERT INTO "polls_question"
("question_text", "pub_date") VALUES
('Wassup?', '2019-05-28T18:54:50.767481+00:00'::timestamptz) RETURNING
"polls_question"."id" /*controller='index',framework='django%3A2.2.1',route='%5Epolls/%24'*/
```
If you want the opencensus attributes included, you must set the
``SQLCOMMENTER_WITH_OPENCENSUS`` setting to ``True`` and include
``'opencensus.ext.django.middleware.OpencensusMiddleware'`` before
``'google.cloud.sqlcommenter.django.middleware.SqlCommenter',`` in your ``MIDDLEWARE``
setting.

### SQLAlchemy

Attach the provided event listener to the `before_cursor_execute` event of the database engine, with `retval=True`. All queries executed with that engine will have the SQL comment prepended to them.

```python
import sqlalchemy
from google.cloud.sqlcommenter.sqlalchemy.executor import BeforeExecuteFactory

engine = sqlalchemy.create_engine(...)
listener = BeforeExecuteFactory(with_db_driver=True, with_db_framework=True, with_opencensus=True)
sqlalchemy.event.listen(engine, 'before_cursor_execute', listener, retval=True)
engine.execute(...) # comment will be added before execution
```

which will produce a backend log such as when viewed on Postgresql
```shell
2019-05-28 11:52:06.527 PDT [64087] LOG:  statement: SELECT * FROM polls_question
/*db_driver='psycopg2',framework='sqlalchemy%3A1.3.4',
traceparent='00-5bd66ef5095369c7b0d1f8f4bd33716a-c532cb4098ac3dd2-01',
tracestate='congo%%3Dt61rcWkgMzE%%2Crojo%%3D00f067aa0ba902b7'*/
```

### Psycopg2

Use the provided cursor factory to generate database cursors. All queries executed with such cursors will have the SQL comment prepended to them.

```python
import psycopg2
from google.cloud.sqlcommenter.psycopg2.extension import CommenterCursorFactory

cursor_factory = CommenterCursorFactory(
    with_db_driver=True, with_dbapi_level=True, with_dbapi_threadsafety=True,
    with_driver_paramstyle=True, with_libpq_version=True, with_opencensus=True)
conn = psycopg2.connect(..., cursor_factory=cursor_factory)
cursor = conn.cursor()
cursor.execute(...) # comment will be added before execution
```

which will produce a backend log such as when viewed on Postgresql
```shell
2019-05-28 02:33:25.287 PDT [57302] LOG:  statement: SELECT * FROM
polls_question /*db_driver='psycopg2%%3A2.8.2%%20%%28dt%%20dec%%20pq3%%20ext%%20lo64%%29',
dbapi_level='2.0',dbapi_threadsafety=2,driver_paramstyle='pyformat',
libpq_version=100001,traceparent='00-5bd66ef5095369c7b0d1f8f4bd33716a-c532cb4098ac3dd2-01',
tracestate='congo%%3Dt61rcWkgMzE%%2Crojo%%3D00f067aa0ba902b7'*/
```

## Options

With Django, each option translates to a Django setting by uppercasing it and prepending `SQLCOMMENTER_`. For example, `with_framework` is controlled by the django setting `SQLCOMMENTER_WITH_FRAMEWORK`.

| Options | Included by default? | Django | SQLAlchemy | psycopg2 | Notes |
| ------- | :------------------: | ------ | ---------- | -------- | :---: |
| `with_framework` | :heavy_check_mark: | [Django version](https://docs.djangoproject.com/en/stable/releases/)  | [Flask version](http://flask.pocoo.org/) | [Flask version](http://flask.pocoo.org/) |
| `with_controller` | :heavy_check_mark: | [Django view](https://docs.djangoproject.com/en/stable/ref/urlresolvers/#django.urls.ResolverMatch.view_name)  | [Flask endpoint](http://flask.pocoo.org/docs/1.0/api/#flask.Flask.endpoint) | [Flask endpoint](http://flask.pocoo.org/docs/1.0/api/#flask.Flask.endpoint) |
| `with_route` | :heavy_check_mark: | [Django route](https://docs.djangoproject.com/en/stable/ref/urlresolvers/#django.urls.ResolverMatch.route)  | [Flask route](http://flask.pocoo.org/docs/1.0/api/#flask.Flask.route) | [Flask route](http://flask.pocoo.org/docs/1.0/api/#flask.Flask.route) |
| `with_app_name ` | | [Django app name](https://docs.djangoproject.com/en/stable/ref/urlresolvers/#django.urls.ResolverMatch.app_name) | | |
| `with_opencensus` | | [W3C TraceContext.Traceparent](https://www.w3.org/TR/trace-context/#traceparent-field), [W3C TraceContext.Tracestate](https://www.w3.org/TR/trace-context/#tracestate-field) | [W3C TraceContext.Traceparent](https://www.w3.org/TR/trace-context/#traceparent-field), [W3C TraceContext.Tracestate](https://www.w3.org/TR/trace-context/#tracestate-field) | [W3C TraceContext.Traceparent](https://www.w3.org/TR/trace-context/#traceparent-field), [W3C TraceContext.Tracestate](https://www.w3.org/TR/trace-context/#tracestate-field) | [[1]](#1-opencensus)
| `with_db_driver` | | [Django DB engine](https://docs.djangoproject.com/en/stable/ref/settings/#engine) | [SQLAlchemy DB driver](https://docs.sqlalchemy.org/en/13/core/engines.html#database-urls) | [psycopg2 version](http://initd.org/psycopg/docs/) |
| `with_db_framework` | | | [SQLAlchemy version](https://www.sqlalchemy.org/) | |
| `with_dbapi_threadsafety` | | | | [psycopg2 thread safety](http://initd.org/psycopg/docs/module.html#psycopg2.threadsafety) |
| `with_dbapi_level` | | | | [psycopg2 api level](http://initd.org/psycopg/docs/module.html#psycopg2.apilevel) |
| `with_libpq_version` | | | | [psycopg2 libpq version](http://initd.org/psycopg/docs/module.html#psycopg2.__libpq_version__) |
| `with_driver_paramstyle` | | | | [psycopg2 parameter style](http://initd.org/psycopg/docs/module.html#psycopg2.paramstyle) |

#### [1] `opencensus`

For `opencensus` to work correctly, note that:

* [OpenCensus for Python](https://github.com/census-instrumentation/opencensus-python) must be installed in the python environment.
* Because the W3C TraceContext's `traceparent` and `tracestate` are quite ephemeral per request, including these attributes can have a negative impact on query caching.
