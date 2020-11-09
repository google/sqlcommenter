const { wrapMainKnexAsMiddleware } = require("@google-cloud/sqlcommenter-knex");
const Knex = require("knex");

const sqlcommenterMiddleware = wrapMainKnexAsMiddleware(
  Knex,
  {
    client_timezone: true,
    db_driver: true,
    route: true,
    traceparent: true,
    tracestate: true,
  },
  { TraceProvider: "OpenTelemetry" }
);

const knex = Knex({
  client: "pg",
  connection: {
    host: process.env.DBHOST,
    user: process.env.DBUSERNAME,
    password: process.env.DBPASSWORD,
    database: "postgres",
  },
});

module.exports = {
  knex,
  sqlcommenterMiddleware,
};
