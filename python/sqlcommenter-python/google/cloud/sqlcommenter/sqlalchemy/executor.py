#!/usr/bin/python
#
# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This integration implements the before_cursor_execute hook factory as per:
#   https://kite.com/python/docs/sqlalchemy.events.ConnectionEvents.before_execute
from __future__ import absolute_import

import sqlalchemy
from google.cloud.sqlcommenter import generate_sql_comment
from google.cloud.sqlcommenter.flask import get_flask_info
from google.cloud.sqlcommenter.opencensus import get_opencensus_values


def BeforeExecuteFactory(
        with_framework=True, with_controller=True, with_route=True,
        with_opencensus=False, with_db_driver=False, with_db_framework=False):

    attributes = {
        'framework': with_framework,
        'controller': with_controller,
        'route': with_route,
        'db_driver': with_db_driver,
        'db_framework': with_db_framework,
    }

    def before_cursor_execute(conn, cursor, sql, parameters, context, executemany):
        data = dict(
            # TODO: Figure out how to retrieve the exact driver version.
            db_driver=conn.engine.driver,
            # Driver/framework centric information.
            db_framework='sqlalchemy:%s' % sqlalchemy.__version__,
        )

        # Because sqlalchemy is a plain database connectivity module,
        # folks using it in a web framework such as flask will
        # use it in unison with flask but initialize the parts disjointly,
        # unlike Django which uses ORMs directly as part of the framework.
        data.update(get_flask_info())

        # Filter down to just the requested attributes.
        data = {k: v for k, v in data.items() if attributes.get(k)}

        if with_opencensus:
            data.update(get_opencensus_values())

        sql_comment = generate_sql_comment(**data)

        # TODO: Check if the database type is MySQL and figure out
        # if we should prepend comments because MySQL server truncates
        # logs greater than 1kB.
        # See:
        #  * https://github.com/basecamp/marginalia/issues/61
        #  * https://github.com/basecamp/marginalia/pull/80
        sql += sql_comment

        return sql, parameters

    return before_cursor_execute
