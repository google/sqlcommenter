# sqlcommenter-http

This is a sample application used to test out `sqlcommenter-go-http` instrumentation libraries.

## Installation

We have containerized this application and its related dependencies (i.e. Postgres and MySQL servers) using `docker-compose`.

To start the application:
```sh
docker compose build
docker compose up
```

This will start the application server as well as MySQL and Postgres databases 
(NOTE: At present both database will start even if one is used). 

By default the application server will connect with `Postgres` database. To change that,
update the `db_engine` parameter in the [`Dockerfile`](https://github.com/google/sqlcommenter/blob/master/go/samples/http/Dockerfile) to `mysql`.

To view postgres logs, just use the `docker logs` command:

```sh
docker logs --follow postgres
```

To view MySQL logs, use the provided script `./tail_mysql_log.sh`.


## Execution
Hit the url http://localhost:8081/ and observe the MySQL/Postgres logs to see comments appended.

Alternatively, there are few scripts present in the `curls` folder to execute CRUD operations.
