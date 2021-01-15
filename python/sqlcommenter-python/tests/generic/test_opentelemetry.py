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
from google.cloud.sqlcommenter.opentelemetry import get_opentelemetry_values

from ..compat import mock, skipIfPy2
from ..opentelemetry_mock import mock_opentelemetry_context


@skipIfPy2
class OpenTelemetryTests(TestCase):
    def test_not_installed(self):
        with mock.patch("google.cloud.sqlcommenter.opentelemetry.propagator", new=None):
            with six.assertRaisesRegex(
                self, ImportError, "OpenTelemetry is not installed"
            ):
                get_opentelemetry_values()

    def test_no_values(self):
        self.assertEqual(get_opentelemetry_values(), {})

    def test_trace_parent_and_tracestate(self):
        with mock_opentelemetry_context():
            self.assertEqual(
                get_opentelemetry_values(),
                {
                    "traceparent": "00-000000000000000000000000deadbeef-000000000000beef-00",
                    "tracestate": "some_key=some_value",
                },
            )

    def test_trace_flags(self):
        from opentelemetry import trace
        with mock_opentelemetry_context(trace_flags=trace.TraceFlags.SAMPLED):
            self.assertEqual(
                get_opentelemetry_values(),
                {
                    "traceparent": "00-000000000000000000000000deadbeef-000000000000beef-01",
                    "tracestate": "some_key=some_value",
                },
            )

    def test_empty_trace_state(self):
        with mock_opentelemetry_context(trace_state={}):
            self.assertEqual(
                get_opentelemetry_values(),
                {
                    "traceparent": "00-000000000000000000000000deadbeef-000000000000beef-00",
                },
            )
