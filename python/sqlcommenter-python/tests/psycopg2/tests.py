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

from unittest import TestCase

import psycopg2.extensions
from forbiddenfruit import curse
from google.cloud.sqlcommenter import url_quote
from google.cloud.sqlcommenter.psycopg2.extension import CommenterCursorFactory

from ..compat import mock
from ..opencensus_mock import mock_opencensus_tracer


class Psycopg2TestCase(TestCase):

    def assertSQL(self, sql, **kwargs):
        def execute(self, sql, args=None):
            pass
        mocked_execute = mock.create_autospec(execute, return_value='worked')
        curse(psycopg2.extensions.cursor, 'execute', mocked_execute)
        cursor = CommenterCursorFactory(**kwargs)
        self.assertEqual(cursor.execute(None, 'SELECT 1;'), 'worked')
        mocked_execute.assert_called_with(None, sql, None)


class Tests(Psycopg2TestCase):

    def test_no_args(self):
        self.assertSQL('SELECT 1;')

    def test_db_driver(self):
        self.assertSQL(
            "SELECT 1; /*db_driver='psycopg2%%3A{}'*/".format(url_quote(psycopg2.__version__)),
            with_db_driver=True,
        )

    def test_dbapi_threadsafety(self):
        self.assertSQL(
            "SELECT 1; /*dbapi_threadsafety={}*/".format(psycopg2.threadsafety),
            with_dbapi_threadsafety=True,
        )

    def test_driver_paramstyle(self):
        self.assertSQL(
            "SELECT 1; /*driver_paramstyle='{}'*/".format(psycopg2.paramstyle),
            with_driver_paramstyle=True,
        )

    def test_dbapi_level(self):
        self.assertSQL(
            "SELECT 1; /*dbapi_level='{}'*/".format(url_quote(psycopg2.apilevel)),
            with_dbapi_level=True,
        )

    def test_libpq_version(self):
        self.assertSQL(
            "SELECT 1; /*libpq_version={}*/".format(url_quote(psycopg2.__libpq_version__)),
            with_libpq_version=True,
        )

    def test_opencensus(self):
        with mock_opencensus_tracer():
            self.assertSQL(
                "SELECT 1; /*traceparent='00-trace%%20id-span%%20id-00',"
                "tracestate='congo%%3Dt61rcWkgMzE%%2Crojo%%3D00f067aa0ba902b7'*/",
                with_opencensus=True,
            )


class FlaskTests(Psycopg2TestCase):
    flask_info = {
        'framework': 'flask',
        'controller': 'c',
        'route': '/',
    }

    @mock.patch('google.cloud.sqlcommenter.psycopg2.extension.get_flask_info', return_value=flask_info)
    def test_all_data(self, get_info):
        self.assertSQL(
            "SELECT 1; /*controller='c',framework='flask',route='/'*/",
        )

    @mock.patch('google.cloud.sqlcommenter.psycopg2.extension.get_flask_info', return_value=flask_info)
    def test_framework_disabled(self, get_info):
        self.assertSQL(
            "SELECT 1; /*controller='c',route='/'*/",
            with_framework=False,
        )

    @mock.patch('google.cloud.sqlcommenter.psycopg2.extension.get_flask_info', return_value=flask_info)
    def test_controller_disabled(self, get_info):
        self.assertSQL(
            "SELECT 1; /*framework='flask',route='/'*/",
            with_controller=False,
        )

    @mock.patch('google.cloud.sqlcommenter.psycopg2.extension.get_flask_info', return_value=flask_info)
    def test_route_disabled(self, get_info):
        self.assertSQL(
            "SELECT 1; /*controller='c',framework='flask'*/",
            with_route=False,
        )
