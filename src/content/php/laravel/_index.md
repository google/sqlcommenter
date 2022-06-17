---
title: "Laravel"
date: 2022-05-04T18:30:11-06:00
draft: false
weight: 1
tags: ["php", "laravel"]
---

![](/images/laravel-logo.png)

- [Introduction](#introduction)
- [Requirements](#requirements)
- [Installation](#installation)
    - [composer install](#composer)
    - [source](#source)
- [Fields](#fields)
    - [Sample log entry](#sample-log-entry)
    - [Expected fields](#expected-fields)
- [References](#references)

### Introduction

This package is in the form of [Illuminate Database Connector Wrapper](https://github.com/illuminate/database) whose purpose is to augment a SQL statement right before execution, with information about the controller and user code to help with later making database optimization decisions, after examining the statements.

### Requirements

It requires [php 8](https://www.php.net) & above.

### Installation
At present, we can install `sqlcommenter-laravel` from source.

This middleware can be installed by one of the following methods:
#### Composer
```shell
composer require google/sqlcommenter-laravel
```
#### Source
```shell
git clone https://github.com/google/sqlcommenter.git
```
```php
# Add the following to composer.json 

"repositories": [
{
"type": "path",
"url": "/full/or/relative/path/to/sqlcommenter/php/sqlcommenter-php/packages/sqlcommenter-laravel/"
}
]
```
```shell
composer require "google/sqlcommenter-laravel"
```

### Enabling it

Publish the config file from library to into laravel app using below command

```shell
php artisan vendor:publish --provider="Google\GoogleSqlCommenterLaravel\GoogleSqlCommenterServiceProvider"

```

Add the following class above ``Illuminate\Database\DatabaseServiceProvider::class,
`` in config/app.php
```php
'providers' => [
    ...
    Google\GoogleSqlCommenterLaravel\Database\DatabaseServiceProvider::class,
    Illuminate\Database\DatabaseServiceProvider::class,
    ...
]
```


### Fields

SQL Statements generated are appended with a comment having fields:

* comma separated key-value pairs e.g. `controller='index'`.
* values are SQL escaped i.e. `key='value'`.
* URL-quoted except for the equals(`=`) sign e.g `route='%5Epolls/%24'`. So, should be URL-unquoted when being consumed.

### Sample log entry

After making requests to the sample middleware-enabled `polls` web-app, we can see logs like:

```shell
2022-04-29 13:59:39.922 IST [27935] LOG:  duration: 0.012 ms  execute pdo_stmt_00000003: Select * from users
/*framework='laravel-9.7.0',controller='UserController',action='index',route='%%2Fapi%%2Ftest',db_driver='pgsql',traceparent='00-1cd60708968a61e942b5dacc2d4a5473-7534abe7ed36ce35-01'*/
```

### Expected Fields

Field| Included <br /> by default?                    |Description
---|------------------------------------------------|---
`action`| <div style="text-align: center">&#10004;</div> |The [application namespace](https://laravel.com/docs/9.x/controllers) of the matching URL pattern in your routes/api.php
`controller`| <div style="text-align: center">&#10004;</div> |The [name](https://laravel.com/docs/9.x/controllers) of the matching URL pattern as described in your routes/api.php
`db_driver`| <div style="text-align: center">&#10004;</div> |The name of the php [database engine](https://laravel.com/docs/9.x/database)
`framework`| <div style="text-align: center">&#10004;</div> |The word "laravel" and the version of laravel being used
`route`| <div style="text-align: center">&#10004;</div> |The [route](https://laravel.com/docs/9.x/routing) of the matching URL pattern as described in your routes/api.php
`traceparent`| <div style="text-align: center">&#10004;</div> |The [W3C TraceContext.Traceparent field](https://www.w3.org/TR/trace-context/#traceparent-field) of the OpenTelemetry trace

### End to end examples

Examples are based upon the [sample app](https://github.com/google/sqlcommenter/tree/master/php/sqlcommenter-php/samples/sqlcommenter-laravel).

#### Source code
```php
# config/google_sqlcommenter.php
<?php
return [

    /*
    |
    | These parameters enables/disable whether the specified info can be
    | appended to the query
    */
    'include' => [
        'framework' => true,
        'controller' => true,
        'route' => true,
        'db_driver' => true,
        'opentelemetry' => true,
        'action' => true,
    ]

];
```
From the command line, we run the laravel development server in one terminal:
```shell
php artisan serve
```
And we use [curl](https://curl.haxx.se/) to make an HTTP request in another:
```shell
curl http://127.0.0.1:8000/user/select
```
#### Results

Examining our Postgresql server logs, with the various options

```shell
2022-04-29 13:59:39.922 IST [27935] LOG:  duration: 0.012 ms  execute pdo_stmt_00000003: Select 1/*framework='laravel-9.7.0',controller='UserController',action='index',route='%%2Fapi%%2Ftest',db_driver='pgsql',
traceparent='00-1cd60708968a61e942b5dacc2d4a5473-7534abe7ed36ce35-01'*/
```

### References

Resource|URL
---|---
laravel|https://laravel.com/docs/5.1/quickstart
OpenTelemetry|https://opentelemetry.io
opentelemetry-php|https://github.com/open-telemetry/opentelemetry-php
sqlcommenter on Github|https://github.com/google/sqlcommenter
