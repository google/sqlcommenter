<?php

namespace Google\GoogleSqlCommenterLaravel;

use Illuminate\Support\ServiceProvider;

class GoogleSqlCommenterServiceProvider extends ServiceProvider
{
    /**
     * Register services.
     *
     * @return void
     */
    public function register()
    {
        //
    }

    /**
     * Publishes configuration file.
     *
     * @return void
     */
    public function boot()
    {
        $this->publishes([
            __DIR__.'/config/google_sqlcommenter.php' => config_path('google_sqlcommenter.php')
        ]);

        $this->mergeConfigFrom(
            __DIR__.'/config/google_sqlcommenter.php',
            'google_sqlcommenter'
        );
    }

}