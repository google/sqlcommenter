#!/usr/bin/ruby
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

# frozen_string_literal: true

module Marginalia
  module Comment
    # @return [String, nil]
    def self.traceparent
      return unless (span_context = ::OpenCensus::Trace.span_context)

      # Following W3C Tracecontext.traceparent as per:
      #   https://www.w3.org/TR/trace-context/#traceparent-field
      format(
        '00-%<trace_id>s-%<span_id>s-%<trace_options>.2d',
        trace_id: span_context.trace_id,
        span_id: span_context.span_id,
        trace_options: span_context.trace_options,
      )
    end

    # @return [String, nil]
    def self.tracestate
      # Tracestate is not yet supported in OpenCensus-Ruby.
      # See issue https://github.com/census-instrumentation/opencensus-ruby/issues/88
      # However, when added, please follow:
      #   https://www.w3.org/TR/trace-context/#tracestate-field
      nil
    end
  end
end

Marginalia::Comment.components += %i[traceparent tracestate]
