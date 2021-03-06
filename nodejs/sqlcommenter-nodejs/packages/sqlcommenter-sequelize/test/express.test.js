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
const {wrapSequelizeAsMiddleware} = require('../index');
const chaiHttp = require('chai-http');
const chai = require("chai");
chai.use(chaiHttp);

const expect = chai.expect;
const seq_version = require('sequelize').version;


describe("Comments for Sequelize - Express Middleware", function() {

    const fakeSequelize = {
        dialect: {
            Query: {
                prototype: {
                    run: function(sql, options) {
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

    // Test Express server app
    const app = express();

    before(function() {
        app.use(wrapSequelizeAsMiddleware(fakeSequelize, {client_timezone:true, db_driver:true, route:true}));
        app.get('/test-comment', (req, res) => {
            fakeSequelize.dialect.Query.prototype.run('SELECT * FROM foo')
            .then((sql) => res.json({query: sql}))
            .catch((error) => {
                res.json({error: error.message})
            })
        });
    });

    describe("Cases", function() {

        it("should add comment to generated sql", function(done) {
            chai.request(app)
                .get('/test-comment')
                .end(function (err, res) {
                    const sql = res.body.query;
                    const expected = `SELECT * FROM foo /*client_timezone='%2B00%3A00',db_driver='sequelize%3A${seq_version}',route='%2Ftest-comment'*/`;
                    expect(sql).to.equals(expected);
                });
                done();
        });

        it("should add expected database/driver and http request properties", function(done) {

            const expected = [
                // db/driver properties
                "db_driver",
                "client_timezone",

                // http properties
                "route"
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
        });
    });
});
