# SQLCommenter in Ruby

## SQLCommenter support in Rails


Support for [SQLCommenter](https://google.github.io/sqlcommenter/) in [Ruby on Rails](https://rubyonrails.org/) varies, depending on your versions of Rails.

If you are on Rails version:

 - **7.1 and above:** SQLCommenter is supported by default. Enable [query_log_tags](https://guides.rubyonrails.org/configuring.html#config-active-record-query-log-tags-enabled), and SQLCommenter formatting will be [enabled by default](https://edgeguides.rubyonrails.org/configuring.html#config-active-record-query-log-tags-format).
 - **7.0:** Enable [query_log_tags](https://guides.rubyonrails.org/configuring.html#config-active-record-query-log-tags-enabled) and install the [PlanetScale SQLCommenter gem](https://github.com/planetscale/activerecord-sql_commenter#installation) for SQLCommenter support.
 - **Below 7.0:** Refer to the [sqlcommenter_rails gem](https://github.com/google/sqlcommenter/tree/master/ruby/sqlcommenter-ruby/sqlcommenter_rails) in this directory for adding SQLCommenter support. Note that this requires additional work because you will have to install a fork of the [marginalia](https://github.com/basecamp/marginalia/) gem, which has since been consolidated into Rails 7.0 and up.

## Tracing support in Rails

Tracing support has been implemented in the [marginalia-opencensus gem]: https://github.com/google/sqlcommenter/tree/master/ruby/sqlcommenter-ruby/marginalia-opencensus. Note that this only works for Rails versions below 7.0, before the [marginalia](https://github.com/basecamp/marginalia/) gem was consolidated into Rails.

Re-purposing that gem for Rails versions >=7.0 should only require minor modifications (contributions are welcome!).

