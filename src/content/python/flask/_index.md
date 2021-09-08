---
title: "Flask"
date: 2019-05-31T18:22:11-07:00
draft: false
weight: 40
tags: ["python", "flask"]
---

![](/images/flask-logo.png)

- [Introduction](#introduction)
- [Using it](#using-it)
- [Expected fields](#expected-fields)
- [End to examples](#end-to-end-examples)
    - [With psycopg2](#with-psycopg2)
    - [With sqlalchemy](#with-sqlalchemy)
- [References](#references)

### Introduction

We provide Flask middleware which when coupled with:

* [sqlcommenter-psycopg2](/python/psycopg2)
* [sqlcommenter-sqlalchemy](/python/sqlalchemy)

allow us to retrieve the `controller` and `route` correlated with your source code in your web app.

### Using it

Having successfully installed [sqlcommenter-python's sqlcommenter](/python#install)

```python
from sqlcommenter import FlaskMiddleware

# Then in your flask programs just pass in the app
FlaskMiddleware(app)
```

### Expected fields

This Flask integration when coupled with compatible drivers will place the following fields

Field|Included by default|Description|Turn if off by
---|---|---|---
controller|<div style="text-align: center">&#10004;</div>|The function being used to service an HTTP request|`with_controller=False`
framework|<div style="text-align: center">&#10004;</div>|"flask:<FLASK_VERSION>"|`with_framework=False`
route|<div style="text-align: center">&#10004;</div>|The pattern used to match an HTTP request|`with_route=False`

### End to end examples
#### With psycopg2
{{<highlight python>}}
#!/usr/bin/env python3

import psycopg2
import json
import flask
app = flask.Flask(__name__)

from google.cloud.sqlcommenter import FlaskMiddleware
from google.cloud.sqlcommenter.psycopg2.extension import CommenterCursorFactory

conn = None

@app.route('/polls')
def get_polls():
    cursor = conn.cursor()
    cursor.execute("SELECT * FROM polls_question")
    str_polls = list(map(lambda s: str(s), cursor))

    cursor.close()
    return json.dumps(str_polls)

def main():
    global conn

    try:
        conn = psycopg2.connect(user='', password='$postgres$',
                host='127.0.0.1', port='5432', database='quickstart_py',
                cursor_factory=CommenterCursorFactory())

        # Now enable the middleware.
        FlaskMiddleware(app)

        # Finally run the Flask web app.
        app.run(host='localhost', port=8088, threaded=True)
    except Exception as e:
        print('Encountered exception %s'%(e))

    finally:
        if conn:
            conn.close()


if __name__ == '__main__':
    main()
{{</highlight>}}

which when run by `python3 main.py` and on visiting http://localhost:8088/polls we can see on our database logs

```shell
2019-06-08 12:19:11.284 PDT [70984] LOG:  statement: SELECT * FROM polls_question
/*controller='get_polls',db_driver='psycopg2',framework='sqlalchemy%3A1.3.4',
route='/polls',traceparent='00-5b3df77064f35f091e89fb40022e2a1d-9bbd4868cf0ba2c3-01'*/
```

#### With sqlalchemy

Having successfully installed [google-cloud-sqlcommenter](/python/sqlalchemy) you can now just run
{{<highlight python>}}
#!/usr/bin/env python3

import json
import flask
app = flask.Flask(__name__)

from sqlalchemy import create_engine, event
from sqlcommenter import FlaskMiddleware
from sqlcommenter.sqlalchemy.executor import BeforeExecuteFactory

engine = None

@app.route('/polls')
def get_polls():
    result_proxy = engine.execute("SELECT * FROM polls_question")
    str_polls = list(map(lambda s: str(s), result_proxy))

    result_proxy.close()
    return json.dumps(str_polls)

def main():
    global engine
    engine = create_engine("postgresql://:$postgres$@127.0.0.1:5432/quickstart_py")
    event.listen(engine, 'before_cursor_execute', BeforeExecuteFactory(), retval=True)

    FlaskMiddleware(app)
    app.run(host='localhost', port=8089, threaded=True)

if __name__ == '__main__':
    main()
{{</highlight>}}

which when run by `python3 main.py` and on visiting http://localhost:8089/polls we can see on our database logs

```shell
2019-06-08 12:17:59.518 PDT [73546] LOG:  statement: SELECT * FROM polls_question
/*controller='get_polls',db_driver='psycopg2%3A2.8.2%20%28dt%20dec%20pq3%20ext%20lo64%29',
dbapi_level='2.0',dbapi_threadsafety=2,driver_paramstyle='pyformat',
framework='flask%3A1.0.3',libpq_version=100001,route='/polls'*/
```

### References

Resource|URL
---|---
flask web framework|http://flask.pocoo.org/
sqlcommenter-psycopg2+flask|[/python/psycopg2#with-flask](/python/psycopg2#with-flask)
sqlcommenter-sqlalchemy+flask|[/python/sqlalchemy#with-flask](/python/sqlalchemy#with-flask)
