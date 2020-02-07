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

import six
from google.cloud.sqlcommenter.opencensus import get_opencensus_values

from ..compat import mock
from ..opencensus_mock import mock_opencensus_tracer


class OpenCensusTests(TestCase):
    def test_not_installed(self):
        with mock.patch('google.cloud.sqlcommenter.opencensus.execution_context', new=None):
            with six.assertRaisesRegex(self, ImportError, 'opencensus is not installed'):
                get_opencensus_values()

    def test_no_values(self):
        self.assertEqual(get_opencensus_values(), {})

    def test_values(self):
        with mock_opencensus_tracer():
            self.assertEqual(get_opencensus_values(), {
                'traceparent': '00-trace id-span id-00',
                'tracestate': 'congo=t61rcWkgMzE,rojo=00f067aa0ba902b7',
            })
