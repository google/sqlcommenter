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

require 'marginalia'
require 'opencensus'

require_relative 'opencensus/marginalia_components'
require_relative 'opencensus/version'

module Marginalia
  module OpenCensus
    class << self
      # @return [String] The string to join span names with.
      attr_accessor :join_span_names_with

      # Decides which items to include in the `span_names` component.
      # @return A callable that accepts a single {::OpenCensus::Trace::SpanContext} and returns a {Boolean}.
      #   The span is included in the output only if this callable returns `true`.
      attr_accessor :span_names_filter

      # Formats an individual span name.
      # @return A callable that accepts a single {::OpenCensus::Trace::SpanContext} and returns a {String}.
      attr_accessor :span_name_formatter
    end

    self.join_span_names_with = '~'
    self.span_names_filter = ->(_span_context) { true }
    self.span_name_formatter = ->(span_context) { span_context.this_span.name }

    # Enumerates all span contexts from the given one until the root one.
    # The root context itself is not included.
    #
    # @param [::OpenCensus::Trace::SpanContext] span_context
    # @return [Enumerable<::OpenCensus::Trace::SpanContext>]
    def self.path_to_root(span_context)
      return to_enum(:path_to_root, span_context) unless block_given?

      cur = span_context
      while (parent = cur.parent)
        yield cur
        cur = parent
      end
    end
  end
end
