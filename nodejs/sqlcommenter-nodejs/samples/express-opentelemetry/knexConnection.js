const { wrapMainKnexAsMiddleware } = require("@google-cloud/sqlcommenter-knex");
const Knex = require("knex");

const sqlcommenterMiddleware = wrapMainKnexAsMiddleware(
  Knex,
  {
    traceparent: true,
    tracestate: true,

    // These are optional and will cause a high cardinality burst traced queries
    db_driver: false,
    route: false,
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
