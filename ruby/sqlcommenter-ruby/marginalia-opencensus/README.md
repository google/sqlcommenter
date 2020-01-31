# marginalia-opencensus

Adds [OpenCensus] support to [Marginalia].

[Marginalia]: https://github.com/basecamp/marginalia/
[OpenCensus]: https://opencensus.io/

## Installation

Add this line to your application's Gemfile:

```ruby
gem 'marginalia-opencensus'
```

And then execute:

    $ bundle

Or install it yourself as:

    $ gem install marginalia-opencensus

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

marginalia-opencensus is tested with multiple Rails versions.

We use the following gems for testing marginalia-opencensus:

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
