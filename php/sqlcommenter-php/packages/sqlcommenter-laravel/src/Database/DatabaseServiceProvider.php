<?php
namespace Google\GoogleSqlCommenterLaravel\Database;

use Illuminate\Support\ServiceProvider;

class DatabaseServiceProvider extends ServiceProvider
{
    /**
    * Register the service provider.
    *
    * @return void
    */
    public function boot() {

        Connection::resolverFor('mysql', function($connection, $database, $prefix, $config) {
            return new MySqlConnection($connection, $database, $prefix, $config);
        });

        Connection::resolverFor('pgsql', function($connection, $database, $prefix, $config) {
            return new PostgresConnection($connection, $database, $prefix, $config);
        });

        Connection::resolverFor('sqlite', function($connection, $database, $prefix, $config) {
            return new SQLiteConnection($connection, $database, $prefix, $config);
        });

        Connection::resolverFor('sqlsrv', function($connection, $database, $prefix, $config) {
            return new SqlServerConnection($connection, $database, $prefix, $config);
        });
    }

}
