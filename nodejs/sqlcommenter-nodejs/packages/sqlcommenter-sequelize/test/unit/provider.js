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
const openTelemetryMock = sinon.spy();

// Mock the dependencies of provider to make use of faked methods for adding trace context
rewiremock('../../provider/opentelemetry').with({ addW3CTraceContext: openTelemetryMock });

// Load the provider module with the appropriate mocks
rewiremock.enable();
const provider = require('../../provider');

// A helper method to test which provider was called
const verifyProviderUsed = function (traceProvider, spy, used) {
    provider.attachComments(traceProvider, {});
    expect(spy.called).to.equal(used);
};

describe("Provider", () => {
    describe("attachComment", () => {

        beforeEach(() => {
            openTelemetryMock.resetHistory();
        })

        it("should default to OpenTelemetry when no options are provided", () => {
            verifyProviderUsed(undefined, openTelemetryMock, true);
        });

        it("should default to OpenTelemetry when null is provided", () => {
            verifyProviderUsed(null, openTelemetryMock, true);
        });

        it("should default to OpenTelemetry when invalid options are provided", () => {
            verifyProviderUsed("bad trace library name", openTelemetryMock, true);
        });

        it("should use OpenTelemetry when the name is provided", () => {
            verifyProviderUsed("opentelemetry", openTelemetryMock, true);
        });

        it("should accept an arbitrary capitalization of OpenTelemetry", () => {
            verifyProviderUsed("OpenTeleMetRY", openTelemetryMock, true);
        });

    });
});

// Unload module mocks
rewiremock.disable();