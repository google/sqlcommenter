# sqlcommenter_rails

[sqlcommenter] for [Ruby on Rails].

Powered by [marginalia] and [marginalia-opencensus].

[sqlcommenter]: #todo
[Ruby on Rails]: https://rubyonrails.org/
[marginalia]: https://github.com/basecamp/marginalia/
[marginalia-opencensus]: https://github.com/google/sqlcommenter/tree/master/ruby/sqlcommenter-ruby/marginalia-opencensus

## Installation

Currently, this gem is not released on rubygems.
But can be installed from source.

The gem requires functionality provided by an [open PR](https://github.com/basecamp/marginalia/pull/130) to [marginalia](https://github.com/basecamp/marginalia). Install the PR by cloning [modulitos's fork of marginalia](https://github.com/modulitos/marginalia) one directory above this folder.

```bash
git clone https://github.com/modulitos/marginalia.git ../marginalia
```

Add the following lines to your application's Gemfile:

```ruby
gem 'sqlcommenter_rails', path: '../sqlcommenter_rails'
gem 'marginalia', path: '../marginalia'
gem 'marginalia-opencensus', path: '../marginalia-opencensus'
```

Install dependencies:

```bash
bin/setup
```

Please look at the [sqlcommenter_rails_demo](https://github.com/google/sqlcommenter/tree/master/ruby/sqlcommenter-ruby/sqlcommenter_rails_demo#sqlcommenter_rails-demo) for an example on how to install and use this gem in your project.


## Usage

This gem registers an `opencensus` component and appends it to the list of default Marginalia components.

In the default configuration, OpenCensus trace will be automatically added to the end of Marginalia comments.

By the default the trace contains OpenCensus Span names from the current one to root, joined with `~`.

See `Marginalia::OpenCensus` documentation for more configuration options.

## Development

After checking out the repo, run `bin/setup` to install dependencies.
Then, run `bundle exec rake` to run the tests (more on testing below).

You can also run `bin/console` for an interactive prompt that will allow you to experiment.

To install this gem onto your local machine, run `bundle exec rake install`.
To release a new version, update the version number in `version.rb`, and then run `bundle exec rake release`,
which will create a git tag for the version, push git commits and tags,
and push the `.gem` file to [rubygems.org](https://rubygems.org).

## Testing

sqlcommenter_rails is tested with multiple Rails versions.

We use the following gems for testing sqlcommenter_rails:

1. [RSpec](https://github.com/rspec/rspec) + [RSpec Rails](https://github.com/rspec/rspec-rails) as the testing framework.
2. [combustion](https://github.com/pat/combustion) for integration tests with a Rails application.
3. [wwtd](https://github.com/grosser/wwtd) for emulating Travis CI locally.

To run the test suite with the latest release of Rails, run:

```bash
bundle exec rake
```

To run the entire test suite (all supported Rails version, rubocop, etc), run:

```bash
bundle exec wwtd
```

To start a web server with the embedded test application, run:

```bash
bin/rails s
```

## Contributing

Bug reports and pull requests are welcome on GitHub at **TODO: REPO URL**.

This project is intended to be a safe, welcoming space for collaboration, and contributors are expected to adhere to the
[Contributor Covenant](http://contributor-covenant.org) code of conduct.

## License

TODO: license.

## Code of Conduct

TODO: code of conduct.
