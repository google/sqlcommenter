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

RSpec.describe Marginalia::Comment do # rubocop:disable RSpec/FilePath
  describe 'traceparent' do
    it 'is set to the correct value inside OpenCensus trace' do
      mocked_span_context = OpenStruct.new(
        trace_id: '123321',
        span_id: '987789',
        trace_options: '42',
      )
      expected_traceparent = '00-123321-987789-42'

      allow(::OpenCensus::Trace).to receive(:span_context).and_return(mocked_span_context)
      ::OpenCensus::Trace.start_request_trace do
        expect(described_class.traceparent).to eq(expected_traceparent)
      end
    end

    it 'is set to nil outside of OpenCensus trace' do
      expect(described_class.traceparent).to eq(nil)
    end
  end
end
