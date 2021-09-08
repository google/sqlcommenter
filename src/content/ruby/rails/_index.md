---
title: "Ruby on Rails"
date: 2019-06-13T16:04:11-06:00
draft: false
weight: 1
tags: ["ruby", "rubyonrails", "rails", "activerecord", "marginalia"]
---

![](/images/activerecord_marginalia-logo.png)

- [Introduction](#introduction)
- [Installation](#installation)
- [Usage](#usage)
- [Fields](#fields)
- [End-to-end example](#end-to-end-example)
    - [Results](#results)
- [References](#references)

### Introduction

[sqlcommenter_rails] adds comments to your SQL statements.

It is powered by [marginalia] and also adds [OpenCensus] information to the
comments if you use the [opencensus gem].

[sqlcommenter_rails] configures [marginalia] and [marginalia-opencensus] to
match the SQLCommenter format.

### Installation

Add this line to your application's Gemfile

```ruby
gem 'sqlcommenter_rails'
```

Then run `bundle` and restart your Rails server.

To enable [OpenCensus] support, add the [`opencensus`][opencensus gem] gem to
your Gemfile, and add the following line in the beginning of
`config/application.rb`:

```ruby
require 'opencensus/trace/integrations/rails'
```

### Usage

All of the SQL queries initiated by the application will now come with comments!

### Fields

The following fields will be added to your SQL statements as comments:

Field|Included <br /> by default?|Description|Provided by
---|---|---|---
`action` |<div style="text-align: center">&#10004;</div>| Controller action name | [marginalia]
`application` |<div style="text-align: center">&#10004;</div>|Application name | [marginalia]
`controller` |<div style="text-align: center">&#10004;</div>| Controller name | [marginalia]
`controller_with_namespace` |<div style="text-align: center">&#10060;</div>| Full classname (including namespace) of the controller | [marginalia]
`database` |<div style="text-align: center">&#10060;</div>| Database name | [marginalia]
`db_driver` |<div style="text-align: center">&#10004;</div>| Database adapter class name | [sqlcommenter_rails]
`db_host` |<div style="text-align: center">&#10060;</div>| Database hostname | [marginalia]
`framework` |<div style="text-align: center">&#10004;</div>| `rails_v` followed by `Rails::VERSION` | [sqlcommenter_rails]
`hostname` |<div style="text-align: center">&#10060;</div>| Socket.gethostname | [marginalia]
`job` |<div style="text-align: center">&#10060;</div>| Classname of the ActiveJob being performed | [marginalia]
`line`|<div style="text-align: center">&#10060;</div>| File and line number calling query | [marginalia]
`pid` |<div style="text-align: center">&#10060;</div>| Current process id | [marginalia]
`route` |<div style="text-align: center">&#10004;</div>| Request's full path | [sqlcommenter_rails]
`socket` |<div style="text-align: center">&#10060;</div>| Database socket | [marginalia]
`traceparent`|<div style="text-align: center">&#10060;</div>|The [W3C TraceContext.Traceparent field](https://www.w3.org/TR/trace-context/#traceparent-field) of the OpenCensus trace | [marginalia-opencensus]

To include the `traceparent` field, install the [marginalia-opencensus] gem and it will be automatically included by default.

To change which fields are included, set `Marginalia::Comment.components = [ :field1, :field2, ... ]` in `config/initializers/marginalia.rb` as described in the [marginalia documentation](https://github.com/basecamp/marginalia#components).

### End to end example

A Rails 6 [sqlcommenter_rails] demo is available at:
https://github.com/google/sqlcommenter/tree/master/ruby/sqlcommenter-ruby/sqlcommenter_rails_demo

The demo is a vanilla Rails API application with [sqlcommenter_rails] and
OpenCensus enabled.

First, we create a vanilla Rails application with the following command:

```shell
gem install rails -v 6.0.0.rc1
rails _6.0.0.rc1_ new sqlcommenter_rails_demo --api
```

Then, we add and implement a basic `Post` model and controller:

```shell
bin/rails g model Post title:text
```

```shell
bin/rails g controller Posts index create
```

Implement the controller:

```ruby
# app/controllers/posts_controller.rb
class PostsController < ApplicationController
  def index
    render json: Post.all
  end

  def create
    title = params[:title].to_s.strip
    head :bad_request if title.empty?
    render json: Post.create!(title: title)
  end
end
```

Configure the routes:

```ruby
# config/routes.rb
Rails.application.routes.draw do
 resources :posts, only: %i[index create]
end
```

Then, we add `sqlcommenter_rails` and OpenCensus:

```ruby
# Gemfile
gem 'opencensus'
gem 'sqlcommenter_rails'
```

```ruby
# config/application.rb
require "opencensus/trace/integrations/rails"
```

Finally, we run `bundle` to install the newly added gems:

```shell
bundle
```

Now, we can start the server:

```shell
bin/rails s
```

In a separate terminal, you can monitor the relevant SQL statements in the server
log with the following command:

```bash
tail -f log/development.log | grep 'Post '
```

#### Results

The demo application has 2 endpoints: `GET /posts` and `POST /posts`.

##### GET /posts

```shell
curl localhost:3000/posts
```

```
Post Load (0.1ms)  SELECT "posts".* FROM "posts" /*
action='index',application='SqlcommenterRailsDemo',controller='posts',
db_driver='ActiveRecord::ConnectionAdapters::SQLite3Adapter',
framework='rails_v6.0.0.rc1',route='/posts',
traceparent='00-ff19308b1f17fedc5864e929bed1f44e-6ddace73a9debf63-01'*/
```

##### POST /posts

```shell
curl -X POST localhost:3000/posts -d 'title=my-post'
```

```
Post Create (0.2ms)  INSERT INTO "posts" ("title", "created_at", "updated_at")
VALUES (?, ?, ?) /*action='create',application='SqlcommenterRailsDemo',
controller='posts',db_driver='ActiveRecord::ConnectionAdapters::SQLite3Adapter',
framework='rails_v6.0.0.rc1',route='/posts',
traceparent='00-ff19308b1f17fedc5864e929bed1f44e-6ddace73a9debf63-01'*/
```

### References

| Resource                | URL                                                                                             |
|-------------------------|-------------------------------------------------------------------------------------------------|
| [sqlcommenter_rails]    | https://github.com/google/sqlcommenter/tree/master/ruby/sqlcommenter-ruby/sqlcommenter_rails    |
| [marginalia]            | https://github.com/basecamp/marginalia                                                          |
| [OpenCensus]            | https://opencensus.io/                                                                          |
| The [opencensus gem]    | https://github.com/census-instrumentation/opencensus-ruby                                       |
| [marginalia-opencensus] | https://github.com/google/sqlcommenter/tree/master/ruby/sqlcommenter-ruby/marginalia-opencensus |

[sqlcommenter_rails]: https://github.com/google/sqlcommenter/tree/master/ruby/sqlcommenter-ruby/sqlcommenter_rails
[marginalia]: https://github.com/basecamp/marginalia
[marginalia-opencensus]: https://github.com/google/sqlcommenter/tree/master/ruby/sqlcommenter-ruby/marginalia-opencensus
[OpenCensus]: https://opencensus.io/
[opencensus gem]: https://github.com/census-instrumentation/opencensus-ruby
