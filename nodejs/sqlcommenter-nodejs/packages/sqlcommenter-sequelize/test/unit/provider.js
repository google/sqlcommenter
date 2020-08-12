// Copyright 2020 Google LLC
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

const provider = require('../../provider');
const chai = require('chai');
const expect = chai.expect;
const sinon = require('sinon');
const rewiremock = require('rewiremock');

// Mock the trace provider implementations
const openCensusMock = class openCensusMock {
    static addW3CTraceContext() {}
}
const openTelemetryMock = class openTelemetryMock {
    static addW3CTraceContext() {}
}
openCensusMock.addW3CTraceContext = sinon.stub();
openTelemetryMock.addW3CTraceContext = sinon.stub();

// A helper method to test which provider was called
const verifyProviderUsed = function(traceProvider, mock, used) {
    let comments = {};
    provider.attachComments(traceProvider, comments);
    expect(mock.called()).to.equal(used);
};

describe("Provider", () => {
    describe("attachComment", () => {
        before(() => {
            /*
            rewiremock('../../provider/opencensus').with({
                OpenCensus: openCensusMock
            });
            rewiremock('../../provider/opentelemetry').with({
                OpenTelemetry: openTelemetryMock
            }); */
        });
        it("should default to OpenCensus when no options are provided", () => {
            verifyProviderUsed(undefined, openCensusMock.addW3CTraceContext, true);
            verifyProviderUsed(undefined, openTelemetryMock.addW3CTraceContext, false)
        }); 
    });
});