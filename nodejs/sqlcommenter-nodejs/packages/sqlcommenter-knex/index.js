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

const {tracer} = require('@opencensus/nodejs');
const {hasComment, toW3CTraceContext} = require('./util');

const defaultFields = {
    'route': true,
    'tracestate': false,
    'traceparent': false,
};

/**
 * All available variables for the commenter are on the `util.fields` object
 * passing the include parameter will result in each item being excluded from
 * the commenter output
 * 
 * @param {Object} Knex
 * @param {Object} includes a map of values to be optionally included.
 * @return {void}
 */
exports.wrapMainKnex = (Knex, include={}) => {

    /* c8 ignore next 2 */
    if (Knex.___alreadySQLCommenterWrapped___)
        return;

    const query = Knex.Client.prototype.query;

    // TODO: Contemplate patch for knex's stream prototype
    // in addition to the query for commenterization.

    // Please don't change this prototype from an explicit function
    // to use arrow functions lest we'll get bugs with not resolving "this".
    Knex.Client.prototype.query = function(conn, obj) {

        // If Knex.VERSION() is available, that means they are using a version of knex.js < 0.16.1
        // because as per https://github.com/tgriesser/knex/blob/master/CHANGELOG.md#0161---28-nov-2018
        // Knex.VERSION was removed in favour of `require('knex/package').version`

        const sqlStmt = typeof obj === 'string' ? obj : obj.sql;

        // If a comment already exists, do not insert a new one.
        // See internal issue #20.
        if (hasComment(sqlStmt)) // Just proceed with the next function ASAP
            return query.apply(this, [conn, obj]);

        const knexVersion = getKnexVersion(Knex);
        const comments = {
            db_driver: `knex:${knexVersion}`
        };

        // TODO: Perhaps remove uuid as it is highly ephemeral?
        // comments.uuid = obj.__knexQueryUid;

        if (Knex.__middleware__) {
            const req = Knex.__req__;

            comments['route'] = req.path;

            // TODO: Clear out __req__ for the next usage.
            // Knex.__req__ = null;
        }

        if (tracer.active) {
            // TODO: check if tracer.currentRootSpan can be null/undefined when tracer is active
            toW3CTraceContext(tracer.currentRootSpan, comments);
        }

        const filtering = typeof include === 'object' && include && Object.keys(include).length > 0; 
        // Filter out keys whose values are undefined or aren't to be included by default.
        const keys = Object.keys(comments).filter((key) => {
            /* c8 ignore next 6 */
            if (!filtering)
                return defaultFields[key] && comments[key];

            // Otherwise since we are filtering, we have to
            // see if the field is included and if it set.
            return include[key] && comments[key];
        });

        // Finally sort the keys alphabetically.
        keys.sort();

        const commentStr = keys.map((key) => {
            const uri_encoded_key = encodeURIComponent(key);
            const uri_encoded_value = encodeURIComponent(comments[key]);
            return `${uri_encoded_key}='${uri_encoded_value}'`;
        }).join(',');
        
        if (typeof obj === 'string') {
            obj = {sql: `${sqlStmt} /*${commentStr}*/`};
        } else {
            obj.sql = `${sqlStmt} /*${commentStr}*/`;
        }

        return query.apply(this, [conn, obj]);
    }

    // Finally mark the object as having already been wrapped.
    Knex.___alreadySQLCommenterWrapped___ = true;
}

const resolveKnexVersion = () => {

    try {
        return require('knex/package').version;
    } catch (err) {
        // Perhaps they are using an old version of knex.js.
        // That is because knex.js as per
        // https://github.com/tgriesser/knex/blob/master/CHANGELOG.md#0161---28-nov-2018
        // Knex.VERSION() was removed in favor of `require('knex/package').version`
        return null;
    }
};

// Since resolveKnexVersion performs expensive lookups by imports,
// we should ensure that this is resolved only once at start time.
const resolvedKnexVersion = resolveKnexVersion();

// Use getKnexVersion to find out the version of knex being used.
const getKnexVersion = (Knex) => {
    return Knex && typeof Knex.VERSION === 'function' ? Knex.VERSION() : resolvedKnexVersion;
}

/**
 * All available variables for the commenter are on the `util.fields` object
 * passing the include parameter will result items not available in that map
 * only being included in the comment.
 * 
 * @param {Object} Knex 
 * @param {Object} include A map of variables to include. If unset, we'll use default attributes.
 * @return {Function} A middleware that is compatible with the express framework.
 */
exports.wrapMainKnexAsMiddleware = (Knex, include=null) => {

    exports.wrapMainKnex(Knex, include);

    return (req, res, next) => {

        Knex.__middleware__ = true;
        Knex.__req__ = req;

        // TODO: Perhaps grab the view/controller name.
        //       Usually req.route.path is useful enough to correlate
        //       where/what source code originates.
        next();
    }
}
