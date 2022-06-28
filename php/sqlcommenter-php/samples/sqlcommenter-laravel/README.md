 # sqlcommenter-laravel-sample 

sqlcommenter is a plugin/middleware/wrapper to augment SQL statements from laravel
with comments that can be used later to correlate user code with SQL statements.

## Setup

### Install Dependencies

Install the required packages
```shell
composer install
```
### Publish the config file

Publish the config file from library to into laravel app using below command

```shell
php artisan vendor:publish --provider="Google\GoogleSqlCommenterLaravel\GoogleSqlCommenterServiceProvider"
```
### Add .env file 
Rename .env.example to .env and change to appropriate configs

## Apply migrations
```shell
php artisan migrate
```
### Run the server
```shell
php artisan serve
```
### Run the tests

Run the tests using below command
```shell
./vendor/bin/phpunit tests
```
