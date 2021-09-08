---
title: "FAQ"
date: 2018-07-16T14:46:09-07:00
draft: false
weight: 90
class: "resized-logo"
---

![](/images/sqlcommenter_logo.png)

- [Why sqlcommenter](#why-sqlcommenter)
- [How does sqlcommenter benefit me](#how-does-sqlcommenter-benefit-me)
- [What ORMs does sqlcommenter support](#what-orms-does-sqlcommenter-support)
- [What databases does sqlcommenter support](#what-databases-does-sqlcommenter-support)
- [How do I use sqlcommenter in my application](#how-do-i-use-sqlcommenter-in-my-application)
- [How do I examine the augmented SQL statements](#how-do-i-examine-the-augmented-sql-statements)
- [Where is the specification](#where-is-the-specification)
- [What are the default fields](#what-are-the-default-fields)

### Why sqlcommenter?

Most applications require persistent data yet when database performance goes awry, it is next to impossible to
easily correlate slow queries with source code.


###  How does sqlcommenter benefit me?

* It helps provide observability and can help correlate your source code with slow queries thus guiding you in performance optimization


### What ORMs does sqlcommenter support?

See [the root of this project](/)


### What databases does sqlcommenter support?

When developing sqlcommenter, we've tested it with a couple of databases. Please see [/databases](/databases) for an authoritative list but here are some:


{{<card-vendor href="/databases/postgresql" src="/images/postgresql-logo.png">}}
{{<card-vendor href="/databases/mysql" src="/images/mysql-logo.png">}}
{{<card-vendor href="/databases/mariadb" src="/images/mariadb-logo.png">}}
{{<card-vendor href="https://sqlite.org/cli.html" src="/images/sqlite-logo.png">}}
{{<card-vendor href="https://cloud.google.com/sql/" src="/images/cloud-sql-card.png">}}


#### How do I use sqlcommenter in my application?
If you are using a supported ORM/framework, it shouldn't be a hassle at all to use. Just pick any of the ORMs in your favorite language

{{<card-vendor href="/python/django" src="/images/django-logo.png">}}
{{<card-vendor href="/python/psycopg2" src="/images/psycopg2-logo.png">}}
{{<card-vendor href="/python/sqlalchemy" src="/images/sqlalchemy-logo.png">}}
{{<card-vendor href="/python/flask" src="/images/flask-logo.png">}}
{{<card-vendor href="/ruby/rails" src="/images/activerecord_marginalia-logo.png">}}
{{<card-vendor href="/java/hibernate" src="/images/hibernate-logo.svg">}}
{{<card-vendor href="/java/spring" src="/images/spring-logo.png">}}
{{<card-vendor href="/node/knex" src="/images/knex-logo.png">}}
{{<card-vendor href="/node/sequelize" src="/images/sequelize-logo.png">}}

### How do I examine the augmented SQL statements?

If you manage your databases or have access to database server logs, the statements will be logged there. Examine [databases](/databases) for more information how.

### Where is the specification?

The specification is available [here](/spec).

### What are the default fields

Integration|action|controller|framework|route
---|---|---|---|---
<a href="/ruby/activerecord">ActiveRecord<br />/Marginalia</a>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10060;</div>
<a href="/java/hibernate">Hibernate<br />+ Spring</a>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10060;</div>
<a href="/node/knex">knex.js<br />+ express.js</a>|<div style="text-align: center">&#10060;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10060;</div>|<div style="text-align: center">&#10060;</div>|<div style="text-align: center">&#10060;</div>
[Django](/python/django)|<div style="text-align: center">&#10060;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10060;</div>
<a href="/python/psycopg2">psycopg2 <br />+ Flask</a>|<div style="text-align: center">&#10060;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10060;</div>
<a href="/python/sqlalchemy">sqlalchemy<br />+ Flask</a>|<div style="text-align: center">&#10060;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10060;</div>
<a href="/node/sequelize">sequelize.js<br />+ express.js</a>|<div style="text-align: center">&#10060;</div>|<div style="text-align: center">&#10004;</div>|<div style="text-align: center">&#10060;</div>|<div style="text-align: center">&#10060;</div>|<div style="text-align: center">&#10060;</div>

### What ORMs support OpenTelemetry ?
Node.js ORMs(Knex, Sequelize, Express) support OpenTelemetry.
