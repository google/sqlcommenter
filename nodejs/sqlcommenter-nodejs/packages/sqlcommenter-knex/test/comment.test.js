// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

"use strict";

const {wrapMainKnex} = require('../index');
const opencensus_tracing = require('@opencensus/nodejs');
const chai = require("chai");
const {fields} = require('../util');
const {context, trace} = require('@opentelemetry/api');
const {NodeTracerProvider} = require('@opentelemetry/node');
const {AsyncHooksContextManager} = require('@opentelemetry/context-async-hooks');
const {InMemorySpanExporter, SimpleSpanProcessor} = require('@opentelemetry/tracing');

const expect = chai.expect;

describe("Comments for Knex", () => {

    let fakeKnex = {
        Client: {
            prototype: {
                config: { connection: { database: 'fake'}, client: 'fakesql'},
                version: 'fake-server:0.0.X',
                query: (conn, obj) => {
                    return Promise.resolve(obj); // simply returns a resolved promise for inspection.
                }
            }
        },
        VERSION: () => {
            return 'fake:0.0.1';
        }
    };

    before(() => {
        wrapMainKnex(fakeKnex, {db_driver: true})
    });

    describe("Cases", () => {

        it("should add comment to generated sql", (done) => {
            
            const want = "SELECT CURRENT_TIMESTAMP /*db_driver='knex%3Afake%3A0.0.1'*/";
            const obj = {sql: 'SELECT CURRENT_TIMESTAMP'};

            fakeKnex.Client.prototype.query(null, obj).then(({sql}) => {
                expect(sql).equals(want);
            });
            done();
        });

        it("should NOT affix comments to statements with existing comments", (done) => {
            
            const queries = [
                'SELECT * FROM people /* existing */',
                'SELECT * FROM people -- existing'
            ];

            Promise.all([
                fakeKnex.Client.prototype.query(null, queries[0]),
                fakeKnex.Client.prototype.query(null, queries[1])
            ]).then(([a, b]) => {
                expect(a).to.equal(queries[0]);
                expect(b).to.equal(queries[1]);
            });
            done();
        });

        it("should add expected database/driver properties", (done) => {
            const want = [
                "db_driver",
            ];
            fakeKnex.Client.prototype.query(null, 'SELECT * from foo').then(({sql}) => {
                want.forEach((key) => {
                    expect(sql.indexOf(key)).to.be.gt(-1);
                });
            });
            done();
        });

        it("should deterministically sort keys alphabetically", (done) => {
            const want = "SELECT * from foo /*db_driver='knex%3Afake%3A0.0.1'*/";
            fakeKnex.Client.prototype.query(null, {sql: 'SELECT * from foo'}).then(({sql}) => {
                expect(sql).equals(want);
            });
            done();
        });

        it("chaining and repeated calls should NOT indefinitely chain SQL", (done) => {
            
            const want = "SELECT * from foo /*db_driver='knex%3Afake%3A0.0.1'*/";
            
            const obj = {sql: 'SELECT * from foo'};

            fakeKnex.Client.prototype.query(null, obj)
                .then((a) => fakeKnex.Client.prototype.query(null, a))
                .then((b) => fakeKnex.Client.prototype.query(null, b))
                .then((c) => fakeKnex.Client.prototype.query(null, c))
                .then((d) => {
                    expect(d.sql).equals(want);
                });
            
            done();
        });
    });
});


describe("With OpenCensus tracing", () => {

    let fakeKnex = {
        Client: {
            prototype: {
                config: { connection: { database: 'fake'}, client: 'fakesql'},
                version: 'fake-server:0.0.X',
                query: (conn, obj) => {
                    return Promise.resolve(obj); // simply returns a resolved promise for inspection.
                }
            }
        },
        VERSION: () => {
            return 'fake:0.0.1';
        }
    };

    before(() => {
        wrapMainKnex(fakeKnex, {db_driver: true, traceparent: true, tracestate: true}, {TraceProvider: "OpenCensus"});
    });

    it('Starting an OpenCensus trace should produce `traceparent`', (done) => {
            // Let's remember  https://github.com/census-instrumentation/opencensus-node/issues/580

            const traceOptions = {
                samplingRate: 1, // Always sample
            };
            const tracer = opencensus_tracing.start(traceOptions).tracer;

            tracer.startRootSpan({ name: 'with-tracing' }, rootSpan => {
                const obj = {sql: 'SELECT * FROM foo'};
                fakeKnex.Client.prototype.query(null, obj).then((got) => {
                    const augmentedSQL = got.sql;
                    const wantSQL = `SELECT * FROM foo /*db_driver='knex%3Afake%3A0.0.1',traceparent='00-${rootSpan.traceId}-${rootSpan.id}-01'*/`;
                    expect(augmentedSQL).equals(wantSQL);
                    opencensus_tracing.tracer.stop();
                    done();
                });
            });
    });
});

describe("With OpenTelemetry tracing", () => {

    let fakeKnex = {
        Client: {
            prototype: {
                config: { connection: { database: 'fake'}, client: 'fakesql'},
                version: 'fake-server:0.0.X',
                query: (conn, obj) => {
                    return Promise.resolve(obj); // simply returns a resolved promise for inspection.
                }
            }
        },
        VERSION: () => {
            return 'fake:0.0.1';
        }
    };

    // Load OpenTelemetry components
    const provider = new NodeTracerProvider();
    const memoryExporter = new InMemorySpanExporter();
    const spanProcessor = new SimpleSpanProcessor(memoryExporter);
    provider.addSpanProcessor(spanProcessor);
    const tracer = provider.getTracer('default');
    trace.setGlobalTracerProvider(provider);
    let contextManager;

    before(() => {
        contextManager = new AsyncHooksContextManager();
        context.setGlobalContextManager(contextManager.enable());
        wrapMainKnex(fakeKnex, {db_driver: true, traceparent: true, tracestate: true}, {TraceProvider: "OpenTelemetry"});
    });

    after(() => {
        memoryExporter.reset();
        context.disable();
    });

    it('Starting an OpenTelemetry trace should produce `traceparent`', (done) => {
        const rootSpan = tracer.startSpan('rootSpan');
        context.with(trace.setSpan(context.active(), rootSpan), async () => {
            const obj = {sql: 'SELECT * FROM foo'};
            fakeKnex.Client.prototype.query(null, obj).then((got) => {
                const augmentedSQL = got.sql;
                const wantSQL = `SELECT * FROM foo /*db_driver='knex%3Afake%3A0.0.1',traceparent='00-${rootSpan.spanContext().traceId}-${rootSpan.spanContext().spanId}-01'*/`;
                expect(augmentedSQL).equals(wantSQL);
                rootSpan.end();
                done();
            });
        });
    });
});
