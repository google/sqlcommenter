# frozen_string_literal: true

lib = File.expand_path('lib', __dir__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require 'marginalia/opencensus/version'

Gem::Specification.new do |spec| # rubocop:disable Metrics/BlockLength
  spec.name = 'marginalia-opencensus'
  spec.version = Marginalia::OpenCensus::VERSION
  spec.authors = ['Google developers']
  spec.email = ['sqlcommenter@googlegroups.com']

  spec.summary = 'Adds OpenCensus support to Marginalia.'
  # spec.homepage = "TODO: Put your gem's website or public repo URL here."
  # spec.license = "TODO
  spec.required_ruby_version = '>= 2.3.0'

  # Prevent pushing this gem to RubyGems.org. To allow pushes either set the 'allowed_push_host'
  # to allow pushing to a single host or delete this section to allow pushing to any host.
  if spec.respond_to?(:metadata)
    # spec.metadata["allowed_push_host"] = "TODO: Set to 'http://mygemserver.com'"

    # spec.metadata["homepage_uri"] = spec.homepage
    # spec.metadata["source_code_uri"] = "TODO: Put your gem's public repo URL here."
    # spec.metadata["changelog_uri"] = "TODO: Put your gem's CHANGELOG.md URL here."
  else
    fail 'RubyGems 2.0 or newer is required to protect against ' \
      'public gem pushes.'
  end

  # Specify which files should be added to the gem when it is released.
  # The `git ls-files -z` loads the files in the RubyGem that have been added into git.
  spec.files = Dir.chdir(File.expand_path(__dir__)) do
    `git ls-files -z`.split("\x0").reject { |f| f.match(%r{^(test|spec|features)/}) }
  end
  spec.bindir = 'exe'
  spec.executables = spec.files.grep(%r{^exe/}) { |f| File.basename(f) }
  spec.require_paths = ['lib']

  spec.add_dependency 'actionpack', '>= 5.2.0'
  spec.add_dependency 'activerecord', '>= 5.2.0'
  spec.add_dependency 'marginalia', '>= 1.8.0', '< 2'
  spec.add_dependency 'opencensus', '~> 0.4.0'

  spec.add_development_dependency 'bundler', '~> 2.0'
  spec.add_development_dependency 'rake', '~> 12.3'

  # Testing:
  spec.add_development_dependency 'capybara', '~> 3.20.2'
  spec.add_development_dependency 'combustion', '~> 1.1'
  spec.add_development_dependency 'rspec', '~> 3.0'
  spec.add_development_dependency 'rspec-rails', '~> 4.0.0.beta2'
  spec.add_development_dependency 'wwtd', '~> 1.3.0'
end
