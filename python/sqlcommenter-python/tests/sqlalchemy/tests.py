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

import sqlalchemy
from google.cloud.sqlcommenter.sqlalchemy.executor import BeforeExecuteFactory

from ..compat import mock
from ..opencensus_mock import mock_opencensus_tracer


class MockConnection:
    @property
    def engine(self):
        class Engine:
            @property
            def driver(self):
                return 'driver'
        return Engine()


class SQLAlchemyTestCase(TestCase):

    def assertSQL(self, expected_sql, **kwargs):
        before_cursor_execute = BeforeExecuteFactory(**kwargs)
        sql, params = before_cursor_execute(
            MockConnection(),  None, 'SELECT 1;', ('param,'), None, None,
        )
        self.assertEqual(sql, expected_sql)
        self.assertEqual(params, ('param,'))


class Tests(SQLAlchemyTestCase):

    def test_no_args(self):
        self.assertSQL('SELECT 1;')

    def test_db_driver(self):
        self.assertSQL(
            "SELECT 1; /*db_driver='driver'*/",
            with_db_driver=True,
        )

    def test_db_framework(self):
        self.assertSQL(
            "SELECT 1; /*db_framework='sqlalchemy%%3A{}'*/".format(sqlalchemy.__version__),
            with_db_framework=True,
        )

    def test_opencensus(self):
        with mock_opencensus_tracer():
            self.assertSQL(
                "SELECT 1; /*traceparent='00-trace%%20id-span%%20id-00',"
                "tracestate='congo%%3Dt61rcWkgMzE%%2Crojo%%3D00f067aa0ba902b7'*/",
                with_opencensus=True,
            )


class FlaskTests(SQLAlchemyTestCase):
    flask_info = {
        'framework': 'flask',
        'controller': 'c',
        'route': '/',
    }

    @mock.patch('google.cloud.sqlcommenter.sqlalchemy.executor.get_flask_info', return_value=flask_info)
    def test_all_data(self, get_info):
        self.assertSQL(
            "SELECT 1; /*controller='c',framework='flask',route='/'*/",
        )

    @mock.patch('google.cloud.sqlcommenter.sqlalchemy.executor.get_flask_info', return_value=flask_info)
    def test_framework_disabled(self, get_info):
        self.assertSQL(
            "SELECT 1; /*controller='c',route='/'*/",
            with_framework=False,
        )

    @mock.patch('google.cloud.sqlcommenter.sqlalchemy.executor.get_flask_info', return_value=flask_info)
    def test_controller_disabled(self, get_info):
        self.assertSQL(
            "SELECT 1; /*framework='flask',route='/'*/",
            with_controller=False,
        )

    @mock.patch('google.cloud.sqlcommenter.sqlalchemy.executor.get_flask_info', return_value=flask_info)
    def test_route_disabled(self, get_info):
        self.assertSQL(
            "SELECT 1; /*controller='c',framework='flask'*/",
            with_route=False,
        )
