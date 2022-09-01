# Contributing to php-laravel

The contributions for Sqlcommenter laravel should be made to this [repo](https://github.com/google/sqlcommenter)

## Download Source

```shell
git clone https://github.com/google/sqlcommenter
```

## Install from source

Add this to your `composer.json`
```shell
"repositories": [
    {
        "type": "path",
        "url": "/full/or/relative/path/to/sqlcommenter-laravel/package"
    }
]
```

```shell
composer require "google/sqlcommenter-laravel"
```
## Run Unittests
Run unit tests using below command
```shell
./vendor/bin/phpunit tests
```