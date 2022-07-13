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
     * Run a select statement and return a single result.
     *
     * @param  string  $query
     * @param  array  $bindings
     * @param  bool  $useReadPdo
     * @return mixed
     */
    public function selectOne($query, $bindings = [], $useReadPdo = true)
    {
        $query = $this->getSqlComments($query);
        $records = parent::select($query, $bindings, $useReadPdo);

        if (count($records) > 0) {
            return $records;
        }
        return null;
    }

    /**
     * Run a select statement against the database.
     *
     * @param  string  $query
     * @param  array  $bindings
     * @param  bool  $useReadPdo
     * @return array
     */
    public function select($query, $bindings = [], $useReadPdo = true)
    {
        $query = $this->getSqlComments($query);
        $records = parent::select($query, $bindings, $useReadPdo);
        return $records;
    }

    /**
     * Run an insert statement against the database.
     *
     * @param  string  $query
     * @param  array  $bindings
     * @return bool
     */
    public function insert($query, $bindings = [])
    {
        $query = $this->getSqlComments($query);
        $records = parent::insert($query, $bindings);

        return $records;
    }

    /**
     * Run an update statement against the database.
     *
     * @param  string  $query
     * @param  array  $bindings
     * @return int
     */
    public function update($query, $bindings = [])
    {
        $query = $this->getSqlComments($query);
      
        return $this->affectingStatement($query, $bindings);
    }

    /**
     * Run a delete statement against the database.
     *
     * @param  string  $query
     * @param  array  $bindings
     * @return int
     */
    public function delete($query, $bindings = [])
    {

        $query = $this->getSqlComments($query);

        return $this->affectingStatement($query, $bindings);
    }

    private function getSqlComments($query)
    {
        $configurationKey = 'google_sqlcommenter.include.';
        $comment = [];
        $action = null;

        if (!empty(app('request')->route())) {
            $action = app('request')->route()->getAction();
        }
        if (config($configurationKey . 'framework', true)) {
            $comment['framework'] = "laravel-" . app()->version();
        }
        if (config($configurationKey . 'controller', true) and !empty($action['controller'])) {
            $comment['controller'] = explode("@", class_basename($action['controller']))[0];
        }
        if (config($configurationKey . 'action', true) and !empty($action and $action['controller'] and str_contains($action['controller'], '@'))) {
            $comment['action'] = explode("@", class_basename($action['controller']))[1];
        }
        if (config($configurationKey . 'route', true)) {
            $comment['route'] = request()->getRequestUri();
        }
        if (config($configurationKey . 'db_driver', true)) {
            $connection = config('database.default');
            $comment['db_driver'] = config("database.connections.{$connection}.driver");
        }
        if (config($configurationKey . 'opentelemetry', true)) {
            $carrier = Opentelemetry::getOpentelemetryValues();
            $comment = array_merge($comment, $carrier);
        }

        $query=trim($query);

        if ($query[-1] == ';'){
            return rtrim($query ,";"). Utils::formatComments(array_filter(($comment))). ';';
        }
        return $query . Utils::formatComments(array_filter(($comment)));

    }
}
