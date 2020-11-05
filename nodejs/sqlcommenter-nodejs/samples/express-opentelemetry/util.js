module.exports = {
  sleep(time) {
    return new Promise((resolve) => setTimeout(resolve, time));
  },

  logger: {
    warn: console.log,
    info: console.log,
    error: console.log,
    debug: console.log,
  },
};
