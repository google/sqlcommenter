<?php
/*
 * Copyright 2022 Google LLC

 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at

 * http:*www.apache.org/licenses/LICENSE-2.0

 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

namespace Google\GoogleSqlCommenterLaravel\Database;

use Illuminate\Database\Connection as BaseConnection;
use OpenTelemetry\API\Trace\Propagation\TraceContextPropagator;
use Google\GoogleSqlCommenterLaravel\Opentelemetry;
use Google\GoogleSqlCommenterLaravel\Utils;


class Connection extends BaseConnection
{
    /**
     * @inheritDoc
     */
    protected function run($query, $bindings, \Closure $callback)
    {
        return parent::run(
            $this->appendSqlComments($query),
            $bindings,
            $callback
        );
    }

    /**
     * Append SQL comments to the underlying query.
     *
     * @param string $query
     * @return string
     */
    private function appendSqlComments(string $query): string
    {
        static $configurationKey = 'google_sqlcommenter.include';
        $comments = [];
        $action = null;

        if (!empty(request()->route())) {
            $action = request()->route()->getAction();
        }
        if (config("{$configurationKey}.framework", true)) {
            $comments['framework'] = "laravel-" . app()->version();
        }
        if (config("{$configurationKey}.controller", true) && !empty($action['controller'])) {
            $comments['controller'] = explode("@", class_basename($action['controller']))[0];
        }
        if (config("{$configurationKey}.action", true) && !empty($action) && !empty($action['controller']) && str_contains($action['controller'], '@')) {
            $comments['action'] = explode("@", class_basename($action['controller']))[1];
        }
        if (config("{$configurationKey}.route", true)) {
            $comments['route'] = request()->getRequestUri();
        }
        if (config("{$configurationKey}.db_driver", true)) {
            $connection = config('database.default');
            $comments['db_driver'] = config("database.connections.{$connection}.driver");
        }
        if (config("{$configurationKey}.opentelemetry", true)) {
            $carrier = Opentelemetry::getOpentelemetryValues();
            $comments = $comments + $carrier;
        }

        $query = trim($query);
        $hasSemicolon = $query[-1] === ';';
        $query = rtrim($query, ';');

        return $query . Utils::formatComments(array_filter($comments)) . ($hasSemicolon ? ';' : '');
    }
}
