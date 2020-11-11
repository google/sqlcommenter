const { Sequelize, DataTypes } = require("sequelize");
const {
  wrapSequelizeAsMiddleware,
} = require("@google-cloud/sqlcommenter-sequelize");

const sequelize = new Sequelize(
  "postgres",
  process.env.DBUSERNAME,
  process.env.DBPASSWORD,
  {
    host: process.env.DBHOST,
    dialect: process.env.DBDIALECT || "postgres",
  }
);
const sqlcommenterMiddleware = wrapSequelizeAsMiddleware(
  sequelize,
  {
    client_timezone: true,
    db_driver: true,
    route: true,
    traceparent: true,
    tracestate: true,
  },
  { TraceProvider: "OpenTelemetry" }
);

const Todo = sequelize.define(
  "Todo",
  {
    // Model attributes are defined here
    title: {
      type: DataTypes.STRING,
      allowNull: false,
    },
    description: {
      type: DataTypes.STRING,
    },
    done: {
      type: DataTypes.BOOLEAN,
      defaultValue: false,
    },
  },
  {}
);

async function createSomeTodos() {
  await sequelize.sync();

  const boringTasks = [];
  for (let i = 0; i < 1000; ++i) {
    boringTasks.push({
      title: `Boring task ${i}`,
      description: "A mundane task",
      done: true,
    });
  }

  await Todo.bulkCreate([
    { title: "Do dishes" },
    { title: "Buy groceries" },
    {
      title: "Do laundry",
      description: "Finish before Thursday!",
    },
    { title: "Clean room" },
    { title: "Wash car" },
    ...boringTasks,
  ]);
}

module.exports = {
  createSomeTodos,
  sequelize,
  sqlcommenterMiddleware,
  Todo,
};
