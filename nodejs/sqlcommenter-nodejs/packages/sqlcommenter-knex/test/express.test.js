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

'use strict';

const express = require('express');
const {wrapMainKnexAsMiddleware} = require('../index');

const chaiHttp = require('chai-http');
const chai = require("chai");
chai.use(chaiHttp);

const expect = chai.expect;

describe("Comments for Knex - Express Middleware", function() {
    
    let fakeKnex = {
        Client: {
            prototype: {
                config: { connection: { database: 'fake'}, client: 'fakesql'},
                version: 'fake-server:0.0.X',
                query: function(conn, obj) {
                    return Promise.resolve(obj); // simply returns a resolved promise for inspection.
                },
            },
        },
        VERSION: function() {
            return 'fake:0.0.1';
        }
    };

    // Test Express server app
    const app = express();
    
    before(function() {
        // use as an express middleware
        app.use(wrapMainKnexAsMiddleware(fakeKnex, {db_driver: true, route: true}));
        app.get('/test-comment', (req, res) => {
            fakeKnex.Client.prototype.query(null, 'SELECT * from people')
            .then(({sql}) => res.json({query: sql}))
            .catch((error) => {
                res.json({error: error.message})
            })
        });
    })

    describe("Cases", function() {
        it("should add comment to generated sql", function(done) {
            chai.request(app)
                .get('/test-comment')
                .end(function (err, res) {
                    const sql = res.body.query;
                    const expected = "SELECT * from people /*db_driver='knex%3Afake%3A0.0.1',route='%2Ftest-comment'*/";
                    expect(sql).to.equals(expected);
                });
                done();
        });
        
        it("should add expected database/driver and http request properties", function(done) {
            const expected = [

                // db/driver properties
                "db_driver",

                // http properties 
                "route",
            ];

            chai.request(app)
                .get('/test-comment')
                .end(function (err, res) {
                    const query = res.body.query;
                    expected.forEach(function(prop) {
                        expect(query.indexOf(prop)).to.be.greaterThan(-1);
                    });
                });
                done();
        })
    });
});
