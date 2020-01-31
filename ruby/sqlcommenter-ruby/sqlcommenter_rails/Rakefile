# frozen_string_literal: true

require 'bundler/gem_tasks'
require 'rspec/core/rake_task'

APP_RAKEFILE = File.expand_path('spec/internal/Rakefile', __dir__)
load 'rails/tasks/engine.rake'

RSpec::Core::RakeTask.new(:spec)

task default: :spec
