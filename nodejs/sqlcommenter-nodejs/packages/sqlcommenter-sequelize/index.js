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

const { hasComment } = require('./util');
const provider = require('./provider');
const hook = require('./hooks');

const defaultFields = {
    'route': true,
    'tracestate': false,
    'traceparent': false,
};

/**
 * All available variables for the commenter are on the `util.fields` object
 * passing the excludes parameter will result in each item being excluded from
 * the commenter output
 * 
 * @param {Object} sequelize 
 * @param {Object} include - A map of values to be optionally included.
 * @param {Object} options - A configuration object specifying where to collect trace data from. Accepted fields are:
 *  TraceProvider: Should be either 'OpenCensus' or 'OpenTelemetry', indicating which library to collect trace data from.
 * @return {void}
 */
exports.wrapSequelize = (sequelize, include = {}, options = {}) => {

    /* c8 ignore next 2 */
    if (sequelize.___alreadySQLCommenterWrapped___)
        return;

    const run = sequelize.dialect.Query.prototype.run;

    // Please don't change this prototype from an explicit function
    // to use arrow functions lest we'll get bugs with not resolving "this".
    sequelize.dialect.Query.prototype.run = function (sql, sql_options) {

        // If a comment already exists, do not insert a new one.
        // See internal issue #20.
        if (hasComment(sql)) // Just proceed with the next function ASAP
            return run.apply(this, [sql, sql_options]);

        const comments = {
            client_timezone: this.sequelize.options.timezone,
            db_driver: `sequelize:${sequelizeVersion}`
        };

        if (sequelize.__middleware__) {
            const context = hook.getContext();
            if (context && context.req) {
                comments.route = context.req.route.path;
            }
        }

        // Add trace context to comments, depending on the provider.
        provider.attachComments(options.TraceProvider, comments);

        // Filter out keys whose values are undefined or aren't to be included by default.
        const filtering = typeof include === 'object' && include && Object.keys(include).length > 0;
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

        sql = sql.trim()
        if (commentStr && commentStr.length > 0) {
            if (sql.slice(-1) === ';') {
                var trimmedSql = sql.slice(0, -1);
                sql = `${trimmedSql} /*${commentStr}*/;`;

            }
            else {
                sql = `${sql} /*${commentStr}*/`;
            }

        }

        return run.apply(this, [sql, sql_options]);
    }

    // Finally mark the object as having already been wrapped.
    sequelize.___alreadySQLCommenterWrapped___ = true;
}

const resolveSequelizeVersion = () => {
    const sv = require('sequelize').version;
    if (sv)
        return sv;

    return require('sequelize/package').version;
};

// Since resolveSequelizeVersion performs expensive lookups by imports,
// we should ensure that this is resolved only once at start time.
const sequelizeVersion = resolveSequelizeVersion();

/**
 * All available variables for the commenter are on the `util.fields` object
 * passing the include parameter will result items not available in that map
 * only being included in the comment.
 * 
 * @param {Object} sequelize
 * @param {Object} include - A map of variables to include. If unset, we'll use default attributes.
 * @param {Object} options - A configuration object specifying where to collect trace data from. Accepted fields are:
 *  TraceProvider: Should be either 'OpenCensus' or 'OpenTelemetry', indicating which library to collect trace data from.
 * @return {Function} A middleware that is compatible with the express framework. 
 */
exports.wrapSequelizeAsMiddleware = (sequelize, include = null, options) => {

    exports.wrapSequelize(sequelize, include, options);

    return (req, res, next) => {
        data = { req: req };
        hook.createContext(data);
        sequelize.__middleware__ = true;
        next();
    }
}
