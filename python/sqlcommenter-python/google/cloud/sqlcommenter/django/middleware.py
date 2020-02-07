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
from django.db.backends.utils import CursorDebugWrapper
from google.cloud.sqlcommenter import generate_sql_comment
from google.cloud.sqlcommenter.opencensus import get_opencensus_values

django_version = django.get_version()


class SqlCommenter:
    """
    Middleware to append a comment to each database query with details about
    the framework and the execution context.
    """
    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        with connection.execute_wrapper(QueryWrapper(request)):
            return self.get_response(request)


class QueryWrapper:
    def __init__(self, request):
        self.request = request

    def __call__(self, execute, sql, params, many, context):
        with_framework = getattr(django.conf.settings, 'SQLCOMMENTER_WITH_FRAMEWORK', True)
        with_controller = getattr(django.conf.settings, 'SQLCOMMENTER_WITH_CONTROLLER', True)
        with_route = getattr(django.conf.settings, 'SQLCOMMENTER_WITH_ROUTE', True)
        with_app_name = getattr(django.conf.settings, 'SQLCOMMENTER_WITH_APP_NAME', False)
        with_opencensus = getattr(django.conf.settings, 'SQLCOMMENTER_WITH_OPENCENSUS', False)
        with_db_driver = getattr(django.conf.settings, 'SQLCOMMENTER_WITH_DB_DRIVER', False)

        db_driver = context['connection'].settings_dict.get('ENGINE', '')
        resolver_match = self.request.resolver_match

        sql_comment = generate_sql_comment(
            # Information about the controller.
            controller=resolver_match.view_name if resolver_match and with_controller else None,
            # route is the pattern that matched a request with a controller i.e. the regex
            # See https://docs.djangoproject.com/en/stable/ref/urlresolvers/#django.urls.ResolverMatch.route
            # getattr() because the attribute doesn't exist in Django < 2.2.
            route=getattr(resolver_match, 'route', None) if resolver_match and with_route else None,
            # app_name is the application namespace for the URL pattern that matches the URL.
            # See https://docs.djangoproject.com/en/stable/ref/urlresolvers/#django.urls.ResolverMatch.app_name
            app_name=(resolver_match.app_name or None) if resolver_match and with_app_name else None,
            # Framework centric information.
            framework=('django:%s' % django_version) if with_framework else None,
            # Information about the database and driver.
            db_driver=db_driver if with_db_driver else None,
            **get_opencensus_values() if with_opencensus else {}
        )

        # TODO: MySQL truncates logs > 1024B so prepend comments
        # instead of statements, if the engine is MySQL.
        # See:
        #  * https://github.com/basecamp/marginalia/issues/61
        #  * https://github.com/basecamp/marginalia/pull/80
        sql += sql_comment

        # Add the query to the query log if debugging.
        if execute.__self__.__class__ is CursorDebugWrapper:
            execute.__self__.db.queries_log.append(sql)

        return execute(sql, params, many, context)
