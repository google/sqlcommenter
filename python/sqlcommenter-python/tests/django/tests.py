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

import django
from django.db import connection
from django.http import HttpRequest
from django.test import TestCase, override_settings
from django.urls import resolve, reverse
from google.cloud.sqlcommenter.django.middleware import SqlCommenter

from ..opencensus_mock import mock_opencensus_tracer
from . import views


# Query log only active if DEBUG=True.
@override_settings(DEBUG=True)
class Tests(TestCase):

    @staticmethod
    def get_request(path):
        request = HttpRequest()
        request.resolver_match = resolve(path)
        return request

    def get_query(self, path='/'):
        SqlCommenter(views.home)(self.get_request(path))
        # Query with comment added by QueryWrapper and unaltered query added
        # by Django's CursorDebugWrapper.
        self.assertEqual(len(connection.queries), 2)
        return connection.queries[0]

    def assertRoute(self, route, query):
        # route available in Django 2.2 and later.
        if django.VERSION < (2, 2):
            self.assertNotIn('route=', query)
        else:
            self.assertIn("route='%s'" % route, query)

    def test_basic(self):
        query = self.get_query()
        self.assertIn("/*controller='home'", query)
        # Expecting url_quoted("framework='django:'")
        self.assertIn("framework='django%%3A" + django.get_version(), query)
        self.assertRoute('', query)

    def test_basic_disabled(self):
        with self.settings(
                SQLCOMMENTER_WITH_CONTROLLER=False, SQLCOMMENTER_WITH_ROUTE=False,
                SQLCOMMENTER_WITH_FRAMEWORK=False):
            query = self.get_query('/path/')
            self.assertNotIn('controller=', query)
            self.assertNotIn('framework=', query)
            self.assertNotIn('route=', query)

    def test_non_root_path(self):
        query = self.get_query(path='/path/')
        self.assertIn("/*controller='some-path'", query)
        self.assertRoute('path/', query)

    def test_app_path(self):
        with self.settings(SQLCOMMENTER_WITH_APP_NAME=True):
            query = self.get_query(path=reverse('app_urls:app-path'))
            self.assertIn("/*app_name='app_urls'", query)
            self.assertIn("controller='app_urls%%3Aapp-path'", query)
            self.assertRoute('app-urls/app-path/', query)

    def test_app_name_disabled(self):
        query = self.get_query(path=reverse('app_urls:app-path'))
        self.assertNotIn('app_name=', query)

    def test_empty_app_name(self):
        """An empty app_name is omitted."""
        with self.settings(SQLCOMMENTER_WITH_APP_NAME=True):
            query = self.get_query()
            self.assertNotIn("app_name=", query)

    def test_db_driver(self):
        with self.settings(SQLCOMMENTER_WITH_DB_DRIVER=True):
            query = self.get_query()
            self.assertIn("db_driver='django.db.backends.sqlite3'", query)

    def test_db_driver_disabled(self):
        query = self.get_query()
        self.assertNotIn('db_driver=', query)

    def test_opencensus_disabled(self):
        """Opencensus fields are omitted by default."""
        with mock_opencensus_tracer():
            query = self.get_query()
            self.assertNotIn('traceparent', query)
            self.assertNotIn('tracestate', query)

    def test_opencensus_enabled(self):
        with mock_opencensus_tracer(), self.settings(SQLCOMMENTER_WITH_OPENCENSUS=True):
            query = self.get_query()
            self.assertIn(
                "traceparent='00-trace%%20id-span%%20id-00',"
                "tracestate='congo%%3Dt61rcWkgMzE%%2Crojo%%3D00f067aa0ba902b7'",
                query
            )
