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

const { hasComment, toW3CTraceContext } = require('../../util');
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
});
