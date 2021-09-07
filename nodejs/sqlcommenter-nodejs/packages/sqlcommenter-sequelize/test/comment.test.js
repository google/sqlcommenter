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

let sequelize_version = require('sequelize').version;
/* c8 ignore next 3 */
if (!sequelize_version)
    sequelize_version = require('sequelize/package').version;

const {wrapSequelize} = require('../index');
const {fields} = require('../util');
const chai = require("chai");
const expect = chai.expect;
const seq_version = require('sequelize').version;
const opencensus_tracing = require('@opencensus/nodejs');
const {context, trace} = require('@opentelemetry/api');
const {NodeTracerProvider} = require('@opentelemetry/node');
const {AsyncHooksContextManager} = require('@opentelemetry/context-async-hooks');
const {InMemorySpanExporter, SimpleSpanProcessor} = require('@opentelemetry/tracing');

const createFakeSequelize = () => {
    return {
        dialect: {
            Query: {
                prototype: {
                    run: (sql, options) => {
                        return Promise.resolve(sql);
                    },
                    sequelize: {
                        config: {
                            database: 'fake', client: 'fakesql',
                        },
                            options: {
                                databaseVersion: 'fakesql-server:0.0.X',
                                dialect: 'fakesql',
                                timezone: '+00:00',
                            },
                    },
                },
            },
        },
    };
}

describe("Comments for Sequelize", () => {

    const fakeSequelize = createFakeSequelize();

    before(() => {
        wrapSequelize(fakeSequelize, {client_timezone:true, db_driver:true});
    });

    after(() => {
        opencensus_tracing.stop();
    });

    describe("Cases", () => {

        it("should add comment to generated sql", (done) => {
            
            const want = `SELECT * FROM foo /*client_timezone='%2B00%3A00',db_driver='sequelize%3A${seq_version}'*/`;
            const query = 'SELECT * FROM foo';

            fakeSequelize.dialect.Query.prototype.run(query).then((sql) => {
                expect(sql).equals(want);
            });

            done();
        });

        it("should NOT affix comments to statements with existing comments", (done) => {
            
            const q = [
                'SELECT * FROM people /* existing */',
                'SELECT * FROM people -- existing'
            ];

            Promise.all([
                fakeSequelize.dialect.Query.prototype.run(q[0]),
                fakeSequelize.dialect.Query.prototype.run(q[1])
            ]).then(([a, b]) => {
                expect(a).to.equal(q[0]);
                expect(b).to.equal(q[1]);
            });

            done();
        });

        it("should add expected database/driver properties", (done) => {
            const want = [
                "db_driver",
                "client_timezone"
            ];
            fakeSequelize.dialect.Query.prototype.run('SELECT * FROM foo').then((sql) => {
                want.forEach((key) => {
                    expect(sql.indexOf(key)).to.be.gt(-1);
                });
            });
            done();
        });

        it("should deterministically sort keys alphabetically", (done) => {
            const want = `SELECT * FROM foo /*client_timezone='%2B00%3A00',db_driver='sequelize%3A${seq_version}'*/`;
            fakeSequelize.dialect.Query.prototype.run('SELECT * FROM foo').then((got) => {
                expect(got).equals(want);
            });
            done();
        });

        it("chaining and repeated calls should NOT indefinitely chain SQL", (done) => {
            
            const want = `SELECT * FROM foo /*client_timezone='%2B00%3A00',db_driver='sequelize%3A${seq_version}'*/`;
            
            const sql = 'SELECT * FROM foo';

            fakeSequelize.dialect.Query.prototype.run(sql)
                .then((a) => fakeSequelize.dialect.Query.prototype.run(a))
                .then((b) => fakeSequelize.dialect.Query.prototype.run(b))
                .then((c) => fakeSequelize.dialect.Query.prototype.run(c))
                .then((d) => {
                    expect(d).equals(want);
                    done();
                });
        });
    });
});

describe("Excluding all variables", () => {

    const fakeSequelize = createFakeSequelize();

    before(() => {
        wrapSequelize(fakeSequelize, {non_existent: true});
    });

    after(() => {
        opencensus_tracing.stop();
    });

    it("when all variables are excluded, no comment should be generated", (done) => {
        // Allow a re-wrap.
        fakeSequelize.___alreadySQLCommenterWrapped___ = false;
        wrapSequelize(fakeSequelize, {foo:true});

        const want = `SELECT * FROM foo`;
        const sql =  `SELECT * FROM foo`;

        fakeSequelize.dialect.Query.prototype.run(sql).then((got) => {
            expect(got).equals(want);
        });
        done();
    });
});

describe("With OpenCensus tracing", () => {

    const fakeSequelize = createFakeSequelize();

    before(() => {
        wrapSequelize(fakeSequelize, {traceparent: true, tracestate: true}, {TraceProvider: "OpenCensus"});
    });

    after(() => {
            opencensus_tracing.stop();
    });

    it('Starting an OpenCensus trace should produce `traceparent`', (done) => {
            // TODO: Follow-up with https://github.com/census-instrumentation/opencensus-node/issues/580
            // and get a proper guide or file bugs against the project to get the proper
            // way to retrieve spans. For now let's skip this test.
            //  Remember: https://github.com/census-instrumentation/opencensus-node/issues/580

            const traceOptions = {
                samplingRate: 1, // Always sample
            };
            const tracer = opencensus_tracing.start(traceOptions).tracer;

            tracer.startRootSpan({ name: 'with-tracing' }, rootSpan => {
                const sql = 'SELECT * FROM foo';
                fakeSequelize.dialect.Query.prototype.run(sql).then((augmentedSQL) => {
                    const wantSQL = `SELECT * FROM foo /*traceparent='00-${rootSpan.traceId}-${rootSpan.id}-01'*/`;
                    expect(augmentedSQL).equals(wantSQL);
                    opencensus_tracing.tracer.stop();
                    done();
                });
            });
    });
});

describe("With OpenTelemetry tracing", () => {

    const fakeSequelize = createFakeSequelize();

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
        wrapSequelize(fakeSequelize, {traceparent: true, tracestate: true}, {TraceProvider: "OpenTelemetry"});
    });

    after(() => {
        memoryExporter.reset();
        context.disable();
    });

    it('Starting an OpenTelemetry trace should produce `traceparent`', (done) => {
        const rootSpan = tracer.startSpan('rootSpan');

        context.with(trace.setSpan(context.active(), rootSpan),  async () => {
            const sql = 'SELECT * FROM foo';
            let augmentedSQL = await fakeSequelize.dialect.Query.prototype.run(sql);
            const wantSQL = `SELECT * FROM foo /*traceparent='00-${rootSpan.spanContext().traceId}-${rootSpan.spanContext().spanId}-01'*/`;
            expect(augmentedSQL).equals(wantSQL);
            rootSpan.end();
            done();
        });
    });
});
