---
title: "Postgresql"
date: 2019-05-31T16:23:01-07:00
draft: false
weight: 40
tags: ["databases", "postgresql"]
---

![](/images/postgresql-logo.png)

If you host your own instance of Postgresql, you can run it with logs to examine statements by following the guide at
[Runtime config logging](https://www.postgresql.org/docs/11/runtime-config-logging.html#RUNTIME-CONFIG-LOGGING-WHERE)

i.e. on macOS, edit `/usr/local/var/postgres/postgresql.conf` and set `log_destination` to `'stderr'` as per:
```python
#------------------------------------------------------------------------------
# REPORTING AND LOGGING
#------------------------------------------------------------------------------

# - Where to Log -

log_destination = 'stderr'             # Valid values are combinations of
                                    # stderr, csvlog, syslog, and eventlog,
                                    # depending on platform.  csvlog
                                    # requires logging_collector to be on.
```

and then when run as per
```shell
$ PGDATA=/usr/local/var/postgres postgres
```

produces such output
```shell
2019-05-31 16:27:27.482 PDT [19175] LOG:  listening on IPv4 address "127.0.0.1", port 5432
2019-05-31 16:27:27.482 PDT [19175] LOG:  listening on IPv6 address "::1", port 5432
2019-05-31 16:27:27.482 PDT [19175] LOG:  listening on Unix socket "/tmp/.s.PGSQL.5432"
2019-05-31 16:27:27.503 PDT [19176] LOG:  database system was shut down at 2019-05-31 16:27:06 PDT
2019-05-31 16:27:27.508 PDT [19175] LOG:  database system is ready to accept connections
2019-05-31 16:27:31.190 PDT [19183] LOG:  statement: SET TIME ZONE 'UTC'
2019-05-31 16:27:31.195 PDT [19183] LOG:  statement: INSERT INTO "polls_question"
("question_text", "pub_date") VALUES ('Wassup?', '2019-05-31T23:27:31.175952+00:00'::timestamptz)
RETURNING "polls_question"."id" /*controller='index',db_driver='django.db.backends.postgresql',
framework='django%3A2.2.1',route='%5Epolls/%24'*/
```

### References

Resource|URL
---|---
Runtime config logging|https://www.postgresql.org/docs/11/runtime-config-logging.html#RUNTIME-CONFIG-LOGGING-WHERE
