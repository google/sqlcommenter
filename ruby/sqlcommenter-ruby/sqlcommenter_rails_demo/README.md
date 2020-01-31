# sqlcommenter_rails demo

This is a demo [Rails API] application to demonstrate [sqlcommenter_rails] integration.

[Rails API]: https://guides.rubyonrails.org/api_app.html
[sqlcommenter_rails]: https://github.com/google/sqlcommenter/ruby/sqlcommenter-ruby/sqlcommenter_rails

## Setup

Install [Ruby v2.6.3](https://www.ruby-lang.org/en/news/2019/04/17/ruby-2-6-3-released/) if you don't already have it installed.

This demo requires functionality provided by an [open PR](https://github.com/basecamp/marginalia/pull/89) to [marginalia](https://github.com/basecamp/marginalia). Install the PR by cloning [glebm's fork of marginalia](https://github.com/glebm/marginalia) one directory above this demo. Starting from the root directory of this demo:

```bash
git clone https://github.com/glebm/marginalia.git ../marginalia
git -C ../marginalia checkout formatting
```

Install dependencies and prepare the database:

```bash
bin/setup
```

Start the server:

```bash
bin/rails s
```

Run this command in a separate terminal to monitor SQL queries:

```bash
tail -f log/development.log | grep 'Post '
```

## Usage

The demo app has 2 endpoints: `GET /posts` and `POST /posts`.

### GET /posts

```bash
curl localhost:3000/posts
```

<blockquote>
Post Load (0.1ms)  SELECT "posts".* FROM "posts" /*action='index',application='SqlcommenterRailsDemo',controller='posts',db_driver='ActiveRecord::ConnectionAdapters::SQLite3Adapter',framework='rails_v6.0.0.rc1',route='/posts',traceparent='00-828f28f7fb3df3dd07ee6478b2016b2a-a52cad0a8d1425ab-01'*/
</blockquote>

### POST /posts

```bash
curl -X POST localhost:3000/posts -d 'title=my-post'
```

<blockquote>
Post Create (0.2ms)  INSERT INTO "posts" ("title", "created_at", "updated_at") VALUES (?, ?, ?) /*action='create',application='SqlcommenterRailsDemo',controller='posts',db_driver='ActiveRecord::ConnectionAdapters::SQLite3Adapter',framework='rails_v6.0.0.rc1',route='/posts',traceparent='00-ff19308b1f17fedc5864e929bed1f44e-6ddace73a9debf63-01'*/  [["title", "my-post"], ["created_at", "2019-06-08 15:47:59.089692"], ["updated_at", "2019-06-08 15:47:59.089692"]]
</blockquote>
