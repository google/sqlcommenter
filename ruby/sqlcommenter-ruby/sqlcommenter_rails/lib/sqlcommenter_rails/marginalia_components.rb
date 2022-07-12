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
    def self.framework
      "rails_v#{Rails::VERSION::STRING}"
    end

    # @return [String, nil]
    def self.db_driver
      marginalia_adapter&.class&.name
    end

    def self.route
      marginalia_controller&.request&.fullpath
    end
  end
end

Marginalia::Comment.update_formatter!(:sqlcommenter)
Marginalia::Comment.components += %i[framework db_driver route]
Marginalia::Comment.components.sort!
