const { NodeTracerProvider } = require("@opentelemetry/node");
const { BatchSpanProcessor } = require("@opentelemetry/tracing");
const {
  TraceExporter,
} = require("@google-cloud/opentelemetry-cloud-trace-exporter");
const { logger, sleep } = require("./util");

const tracerProvider = new NodeTracerProvider();
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
const {
  Todo,
  sequelize,
  sqlcommenterMiddleware,
} = require("./sequelizeModels");

const tracer = tracerProvider.getTracer(__filename);

async function main() {
  await sequelize.sync();

  const PORT = 8000;

  // SQLCommenter express middleware injects the route into the traces
  app.use(sqlcommenterMiddleware);

  app.get("/", async (req, res) => {
    const span = tracer.startSpan("sleep for no reason (parent)");
    await sleep(250);
    span.end();

    const getRecordsSpan = tracer.startSpan("query records");
    await tracer.withSpan(getRecordsSpan, async () => {
      const todos = await Todo.findAll({ limit: 20 });
      const countByDone = await Todo.findAll({
        attributes: [
          "done",
          [sequelize.fn("COUNT", sequelize.col("id")), "count"],
        ],
        group: "done",
      });
      res.json({ countByDone, todos });
    });
    getRecordsSpan.end();
  });
  app.listen(PORT, async () => {
    console.log(`⚡️[server]: Server is running at https://localhost:${PORT}`);
  });
}

main().catch(console.error);
