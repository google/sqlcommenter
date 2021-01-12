const { LogLevel } = require("@opentelemetry/core");
const { NodeTracerProvider } = require("@opentelemetry/node");
const { BatchSpanProcessor } = require("@opentelemetry/tracing");
const {
  TraceExporter,
} = require("@google-cloud/opentelemetry-cloud-trace-exporter");
const { logger, sleep } = require("./util");

const tracerProvider = new NodeTracerProvider({
  logLevel: LogLevel.DEBUG,
  logger,
});
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
const app = express();
const { sqlcommenterMiddleware, knex } = require("./knexConnection");

const tracer = tracerProvider.getTracer(__filename);

async function main() {
  if (!(await knex.schema.hasTable("Todos"))) {
    console.error("Todos table does not exist. Run `npm run createTodos`");
    process.exit(1);
  }

  const PORT = 8000;

  // SQLCommenter express middleware injects the route into the traces
  app.use(sqlcommenterMiddleware);

  app.get("/", async (req, res) => {
    const span = tracer.startSpan("sleep for no reason (parent)");
    await sleep(250);
    span.end();

    const getRecordsSpan = tracer.startSpan("query records with knex");
    await tracer.withSpan(getRecordsSpan, async () => {
      tracer.startSpan("testing 123").end();
      const todos = await knex.select().table("Todos").limit(20);
      const countByDone = await knex
        .select("done", knex.raw("COUNT(id) as count"))
        .table("Todos")
        .groupBy("done");
      res.json({ countByDone, todos });
    });
    getRecordsSpan.end();
  });
  app.listen(PORT, async () => {
    console.log(`⚡️[server]: Server is running at https://localhost:${PORT}`);
  });
}

main().catch(console.error);
