# Express + OpenTelemetry example app

This example Express app connects to a Cloud SQL database (Postgres or MySQL)
using Sequelize or Knex as the ORM. It is instrumented with OpenTelemetry to
collect trace data and send it to Cloud Trace. SQLCommenter attaches the
OpenTelemetry trace and span ID to the queries for correlation.

## Running

The examples read some environment variables to connect to a Cloud SQL
database.

- `DBHOST` - The hostname, ip, or path to unix socket directory for Cloud SQL
database. For example, if you are using Cloud SQL proxy with a unix socket,
you might pass `/cloudsql/my-project:us-central1:my-db`
- `DBUSERNAME` - Database username
- `DBPASSWORD` - Database password
- `DBDIALECT` - Optional. Should be `postgres` (default) or `mysql`

```bash
# in this directory
$ export DBHOST="<hostname or unix socket path>"
$ export DBUSERNAME="<database user>"
$ export DBPASSWORD="<database password>"
$ export DBDIALECT="<postgres | mysql>"
$ npm install

# before the first run, create some data
$ npm run createTodos
```

To run the sequelize server example:
```bash
$ npm run sequelizeServer
```

Or run the knex server:
```bash
$ npm run knexServer
```

In a separate terminal, hit the server with `curl` and check the output:
```bash
$ curl http://localhost:8000
{"todos":[{"id":11,"title":"Do dishes","description":null,"done":false,"createdAt":"2020-11-06T01:59:50.111Z","updatedAt":"2020-11-06T01:59:50.111Z"},{"id":12,"title":"Buy groceries","description":null,"done":false,"createdAt":"2020-11-06T01:59:50.111Z","updatedAt":"2020-11-06T01:59:50.111Z"},{"id":13,"title":"Do laundry","description":"Finish before Thursday!","done":false,"createdAt":"2020-11-06T01:59:50.111Z","updatedAt":"2020-11-06T01:59:50.111Z"},{"id":14,"title":"Clean room","description":null,"done":false,"createdAt":"2020-11-06T01:59:50.111Z","updatedAt":"2020-11-06T01:59:50.111Z"},{"id":15,"title":"Wash car","description":null,"done":false,"createdAt":"2020-11-06T01:59:50.112Z","updatedAt":"2020-11-06T01:59:50.112Z"}]}

# Original terminal will output
[server]: Server is running at https://localhost:8000
Google Cloud Trace export
Google Cloud Trace batch writing traces
Google Cloud Trace authenticating
Google Cloud Trace got authentication. Initializaing rpc client
batchWriteSpans successfully
Executing (default): SELECT "id", "title", "description", "done", "createdAt", "updatedAt" FROM "Todos" AS "Todo" LIMIT 20; /*client_timezone='%2B00%3A00',db_driver='sequelize%3A6.3.3',route='%2F',traceparent='00-3e2914ebce6af09508dd1ff1128493a8-81d09ab4d8cde7cf-01'*/
Executing (default): SELECT "done", COUNT("id") AS "count" FROM "Todos" AS "Todo" GROUP BY "done"; /*client_timezone='%2B00%3A00',db_driver='sequelize%3A6.3.3',route='%2F',traceparent='00-3e2914ebce6af09508dd1ff1128493a8-81d09ab4d8cde7cf-01'*/
Google Cloud Trace export
Google Cloud Trace batch writing traces
batchWriteSpans successfully
```
