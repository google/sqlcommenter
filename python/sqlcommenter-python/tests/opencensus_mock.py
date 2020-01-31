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

from collections import OrderedDict
from contextlib import contextmanager

from .compat import mock


def get_opencensus_tracer_mock():
    class Tracer:
        @property
        def span_context(self):
            class Context:
                @property
                def span_id(self):
                    return 'span id'

                @property
                def trace_id(self):
                    return 'trace id'

                @property
                def trace_options(self):
                    return TraceOptions()

                @property
                def tracestate(self):
                    od = OrderedDict()
                    od['congo'] = 't61rcWkgMzE'
                    od['rojo'] = '00f067aa0ba902b7'
                    return od

            return Context()
    return Tracer()


class TraceOptions():
    def __init__(self, enabled=False):
        self.enabled = enabled


@contextmanager
def mock_opencensus_tracer():
    path = 'opencensus.trace.execution_context.get_opencensus_tracer'
    with mock.patch(path, get_opencensus_tracer_mock):
        yield
