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

/**
 * fields represent variables that can be made optional for commenter output
 */
exports.fields = {
    CLIENT_TIMEZONE: "client_timezone",
    DB_DRIVER: "db_driver",
    ROUTE: "route",
    TRACE_STATE: "tracestate",
    TRACE_PARENT: "traceparent"
};

/**
 * Inspects the provided sql statement for existing comments.
 * 
 * @param {String} sql The SQL string to inspect
 * @return {Boolean} true if a comment exists, false otherwise
 */
exports.hasComment = (sql) => {

    if (!sql)
        return false;

    // See https://docs.oracle.com/cd/B12037_01/server.101/b10759/sql_elements006.htm
    // for how to detect comments.
    const indexOpeningDashDashComment = sql.indexOf('--');
    if (indexOpeningDashDashComment >= 0)
        return true;

    const indexOpeningSlashComment = sql.indexOf('/*');
    if (indexOpeningSlashComment < 0)
        return false;

    // Check if it is a well formed comment.
    const indexClosingSlashComment = sql.indexOf('*/');

    /* c8 ignore next */
    return indexOpeningSlashComment < indexClosingSlashComment;
}

const latestSpan = (span) => {

    if (!span || !span.isRootSpan())
        return span;

    // Otherwise if it is a root span, we'll try to grab its last child.
    const children = span.spans;
    if (children.length < 1)
        return span;

    /* c8 ignore next */
    return children[children.length - 1];
}
