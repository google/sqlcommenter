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

import os
import sys
import unittest

# Django, Flask, psycopg2, and sqlalchemy are only installed in their
# respective testing environments. See tox.ini for details.

try:
    import django
    from django.test.runner import DiscoverRunner
except ImportError:
    django = None

try:
    import flask
    import pytest
except ImportError:
    flask = None

try:
    import psycopg2
except ImportError:
    psycopg2 = None

try:
    import sqlalchemy
except ImportError:
    sqlalchemy = None


def run_unittests(module):
    testsuite = unittest.TestLoader().discover(module)
    result = unittest.TextTestRunner(verbosity=1).run(testsuite)
    sys.exit(any(result.failures or result.errors))


def run_django_tests():
    os.environ.setdefault('DJANGO_SETTINGS_MODULE', 'tests.django.settings')
    django.setup()
    runner = DiscoverRunner()
    failures = runner.run_tests(['tests.django'])
    sys.exit(failures)


def run_flask_tests():
    failures = pytest.main(['tests/flask/tests.py'])
    sys.exit(failures)


def main():
    if django:
        run_django_tests()
    elif flask:
        run_flask_tests()
    elif psycopg2:
        run_unittests('tests.psycopg2')
    elif sqlalchemy:
        run_unittests('tests.sqlalchemy')
    else:
        run_unittests('tests.generic')


if __name__ == '__main__':
    main()
