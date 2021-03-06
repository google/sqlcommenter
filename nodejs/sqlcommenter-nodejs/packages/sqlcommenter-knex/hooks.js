const hooks = require('async_hooks')
const store = new Map();

const hook = hooks.createHook({
    init: (asyncId, _, triggerAsyncId) => {
        if (store.has(triggerAsyncId)) {
            store.set(asyncId, store.get(triggerAsyncId))
        }
    },
    destroy: (asyncId) => {
        if (store.has(asyncId)) {
            store.delete(asyncId)
        }
    }
});

hook.enable();

const createContext = (data) => {
    store.set(hooks.executionAsyncId(), data);
    return data;
};

const getContext = () => {
    return store.get(hooks.executionAsyncId());
};

module.exports = { createContext, getContext };
