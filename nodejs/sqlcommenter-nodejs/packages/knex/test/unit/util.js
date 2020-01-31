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

const {hasComment, toW3CTraceContext} = require('../../util');
const chai = require("chai");
const expect = chai.expect;

describe("Unit", () => {

    describe("hasComment", () => {

        it("should return true for well-formed comments", () => {
            
            const queries = [
                `SELECT * FROM foo /* existing */`,
                `SELECT * FROM foo -- existing`
            ];

            const want = true;
            queries.forEach(q => {
                expect(hasComment(q)).to.equal(want)
            });
        });
        
        it("should return false when comment is undefined", () => {
            let comment;
            expect(hasComment(comment)).to.equal(false);
        });
    
        it("should return false for malformed comments", () => {
            const queries = [
                "SELECT * FROM people /*",
                "SELECT * FROM people */ /*"
            ];
    
            queries.forEach(q => {
                expect(hasComment(q)).to.equal(false);
            });
        });
    });

    describe("toW3CTraceContext", () => {

        it("should return undefined/null for undefined/null span", () => {
            const badParams = [
                null,
                undefined
            ];
            const obj = {};
            expect(toW3CTraceContext(badParams[0], null)).to.be.null;
            expect(toW3CTraceContext(badParams[1], null)).to.be.null;
            expect(toW3CTraceContext(badParams[0], obj)).to.equal(obj);
            expect(toW3CTraceContext(badParams[1], obj)).to.equal(obj);
        });

        it("should still extract the topSpanContext even if not root span", () => {
            const fakeSpan = {
                spanContext: {
                    traceId: 'ff00000000000001',
                    spanId: 'ee000002',
                    traceState: 'congo=t61rcWkgMzE,rojo=00f067aa0ba902b7',
                },
                isRootSpan() { return false; }
            };
            expect(toW3CTraceContext(fakeSpan, {}).traceparent).to.equal('00-ff00000000000001-ee000002-00');
            expect(toW3CTraceContext(fakeSpan, {}).tracestate).to.equal('congo=t61rcWkgMzE,rojo=00f067aa0ba902b7');
        });

        it("should nothing if the latest span doesn't have children", () => {
            const fakeSpan = {
                spanContext: {},
                spans: [],
                isRootSpan() { return true; }
            };

            const got = toW3CTraceContext(fakeSpan, {});
            expect(got.traceparent).to.equal(undefined);
            expect(got.tracestate).to.equal(undefined);
        });

        it("should return the last span it has children", () => {
            const fakeSpan = {
                spans: [{
                    spanContext: {
                        traceId: 'ff00000000000001',
                        spanId: 'ee000002',
                        traceState: 'congo=t61rcWkgMzE,rojo=00f067aa0ba902b7',
                    },
                },
                {
                    spanContext: {
                        traceId: '11000000000000ff',
                        spanId: '020000ee',
                        traceState: 'brazzaville=t61rcWkgMzE,rondo=00f067aa0ba902b7',
                        options: 0x01, // It is sampled.
                    },
                }],
                isRootSpan() { return true; }
            };

            const got = toW3CTraceContext(fakeSpan, {});
            expect(got.traceparent).to.equal('00-11000000000000ff-020000ee-01');
            expect(got.tracestate).to.equal('brazzaville=t61rcWkgMzE,rondo=00f067aa0ba902b7');
        });
    });
});
