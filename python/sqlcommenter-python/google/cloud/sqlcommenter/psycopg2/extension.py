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

import psycopg2
import psycopg2.extensions
from google.cloud.sqlcommenter import generate_sql_comment
from google.cloud.sqlcommenter.flask import get_flask_info
from google.cloud.sqlcommenter.opencensus import get_opencensus_values


# This integration extends psycopg2.extensions.cursor
# by implementing a custom execute method.
#
# By default, it doesn't enable adding trace_id and span_id
# to SQL comments due to their ephemeral nature. You can opt-in
# by instead setting with_opencensus=True
def CommenterCursorFactory(
        with_framework=True, with_controller=True, with_route=True,
        with_opencensus=False, with_db_driver=False, with_dbapi_threadsafety=False,
        with_dbapi_level=False, with_libpq_version=False, with_driver_paramstyle=False):

    attributes = {
        'framework': with_framework,
        'controller': with_controller,
        'route': with_route,
        'db_driver': with_db_driver,
        'dbapi_threadsafety': with_dbapi_threadsafety,
        'dbapi_level': with_dbapi_level,
        'libpq_version': with_libpq_version,
        'driver_paramstyle': with_driver_paramstyle,
    }

    class CommenterCursor(psycopg2.extensions.cursor):

        def execute(self, sql, args=None):
            data = dict(
                # Psycopg2/framework information
                db_driver='psycopg2:%s' % psycopg2.__version__,
                dbapi_threadsafety=psycopg2.threadsafety,
                dbapi_level=psycopg2.apilevel,
                libpq_version=psycopg2.__libpq_version__,
                driver_paramstyle=psycopg2.paramstyle,
            )

            # Because psycopg2 is a plain database connectivity module,
            # folks using it in a web framework such as flask will
            # use it in unison with flask but initialize the parts disjointly,
            # unlike Django which uses ORMs directly as part of the framework.
            data.update(get_flask_info())

            # Filter down to just the requested attributes.
            data = {k: v for k, v in data.items() if attributes.get(k)}

            if with_opencensus:
                data.update(get_opencensus_values())

            sql += generate_sql_comment(**data)

            return psycopg2.extensions.cursor.execute(self, sql, args)

    return CommenterCursor
