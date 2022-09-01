# Sqlcommenter-laravel

Sqlcommenter is a plugin/middleware/wrapper to augment SQL statements from laravel
with comments that can be used later to correlate user code with SQL statements.

## Installation

```shell
composer require "google/sqlcommenter-laravel"
```
## Usage

Publish the config file from library to into laravel app using below command

```shell
php artisan vendor:publish --provider="Google\GoogleSqlCommenterLaravel\GoogleSqlCommenterServiceProvider"
```

Add the following class above `Illuminate\Database\DatabaseServiceProvider::class`,
 in `config/app.php`
```php
'providers' => [
    ...
    Google\GoogleSqlCommenterLaravel\Database\DatabaseServiceProvider::class,
    Illuminate\Database\DatabaseServiceProvider::class,
    ...
]
```

## Options

With Laravel SqlCommenter, we have configuration to choose which tags to be appended to the comment. It is configurable in `config/google_sqlcommenter.php`

| Field         | Included <br /> by default?                    | Description                                                                                                                 |
| ------------- | ---------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------- |
| `action`      | <div style="text-align: center">&#10004;</div> | The [application namespace](https://laravel.com/docs/9.x/controllers) of the matching URL pattern in your routes/api.php    |
| `controller`  | <div style="text-align: center">&#10004;</div> | The [name](https://laravel.com/docs/9.x/controllers) of the matching URL pattern as described in your routes/api.php        |
| `db_driver`   | <div style="text-align: center">&#10004;</div> | The name of the php [database engine](https://laravel.com/docs/9.x/database)                                                |
| `framework`   | <div style="text-align: center">&#10004;</div> | The word "laravel" and the version of laravel being used                                                                    |
| `route`       | <div style="text-align: center">&#10004;</div> | The [route](https://laravel.com/docs/9.x/routing) of the matching URL pattern as described in your routes/api.php           |
| `traceparent` | <div style="text-align: center">&#10004;</div> | The [W3C TraceContext.Traceparent field](https://www.w3.org/TR/trace-context/#traceparent-field) of the OpenTelemetry trace |
