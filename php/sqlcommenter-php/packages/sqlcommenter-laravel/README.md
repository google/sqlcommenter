# sqlcommenter-laravel [In Development]

sqlcommenter is a plugin/middleware/wrapper to augment SQL statements from laravel
with comments that can be used later to correlate user code with SQL statements.


### Installation

Add this to your `composer.json`
```shell
"repositories": [
    {
        "type": "path",
        "url": "/full/or/relative/path/to/sqlcommenter-laravel/package"
    }
]
```
Install the package
```shell
composer require "google/sqlcommenter-laravel"
```
### Usage

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
### Run unit tests

Run unit tests using below command
```shell
./vendor/bin/phpunit tests
```
