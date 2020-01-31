# frozen_string_literal: true

source 'https://rubygems.org'

gem 'rails', '~> 6.0.0.rc1'
# https://github.com/rails/rails/blob/v6.0.0.rc1/activerecord/lib/active_record/connection_adapters/sqlite3_adapter.rb#L12
gem 'sqlite3', '~> 1.3', '>= 1.3.6'

gemspec

group :debug do
  gem 'byebug'
end

eval_gemfile File.expand_path('rubocop.gemfile', __dir__)
