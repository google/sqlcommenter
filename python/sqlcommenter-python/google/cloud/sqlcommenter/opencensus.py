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

from __future__ import absolute_import

try:
    from opencensus.trace import execution_context
    from opencensus.trace.propagation import trace_context_http_header_format
except ImportError:
    execution_context = None


def get_opencensus_values():
    """
    Return the OpenCensus Trace and Span IDs if Span ID is set in the
    OpenCensus execution context.
    """
    if execution_context:
        span_ctx = execution_context.get_opencensus_tracer().span_context
        span_id = span_ctx.span_id
        if span_id:
            # Insert the W3C TraceContext generated
            w3C_trace_ctx = trace_context_http_header_format.TraceContextPropagator()
            return w3C_trace_ctx.to_headers(span_ctx)
    else:
        raise ImportError('opencensus is not installed.')
    return {}
