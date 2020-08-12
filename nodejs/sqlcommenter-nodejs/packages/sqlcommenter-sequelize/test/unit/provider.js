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

const chai = require('chai');
const expect = chai.expect;
const sinon = require('sinon');
const sinonChai = require("sinon-chai");
const rewiremock = require('../rewiremock');
chai.use(sinonChai);

// Mock the trace provider implementations
const openCensusMock = sinon.spy();
const openTelemetryMock = sinon.spy();

// Mock the dependencies of provider to make use of faked methods for adding trace context
rewiremock('../../provider/opencensus').with({addW3CTraceContext: openCensusMock});
rewiremock('../../provider/opentelemetry').with({addW3CTraceContext: openTelemetryMock});

// Load the provider module with the appropriate mocks
rewiremock.enable();
const provider = require('../../provider');

// A helper method to test which provider was called
const verifyProviderUsed = function(traceProvider, spy, used) {
    provider.attachComments(traceProvider, {});
    expect(spy.called).to.equal(used);
};

describe("Provider", () => {
    describe("attachComment", () => {

        beforeEach(() => {
            openCensusMock.resetHistory();
            openTelemetryMock.resetHistory();
        })

        it("should default to OpenCensus when no options are provided", () => {
            verifyProviderUsed(undefined, openCensusMock, true);
            verifyProviderUsed(undefined, openTelemetryMock, false);
        });

        it("should default to OpenCensus when null is provided", () => {
            verifyProviderUsed(null, openCensusMock, true);
            verifyProviderUsed(null, openTelemetryMock, false);
        });

        it("should default to OpenCensus when invalid options are provided", () => {
            verifyProviderUsed("bad trace library name", openCensusMock, true);
            verifyProviderUsed("bad trace library name", openTelemetryMock, false);
        });

        it("should use OpenCensus when the name is provided", () => {
            verifyProviderUsed("opencensus", openCensusMock, true);
            verifyProviderUsed("opencensus", openTelemetryMock, false);
        });

        it("should accept an arbitrary capitalization of OpenCensus", () => {
            verifyProviderUsed("oPeNceNSus", openCensusMock, true);
            verifyProviderUsed("oPeNceNSus", openTelemetryMock, false);
        });

        it("should use OpenTelemetry when the name is provided", () => {
            verifyProviderUsed("opentelemetry", openCensusMock, false);
            verifyProviderUsed("opentelemetry", openTelemetryMock, true);
        });

        it("should accept an arbitrary capitalization of OpenTelemetry", () => {
            verifyProviderUsed("OpenTeleMetRY", openCensusMock, false);
            verifyProviderUsed("OpenTeleMetRY", openTelemetryMock, true);
        });

    });
});

// Unload module mocks
rewiremock.disable();