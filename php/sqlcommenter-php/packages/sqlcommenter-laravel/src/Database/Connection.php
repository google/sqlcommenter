<?php
namespace Google\GoogleSqlCommenterLaravel\Database;

use Illuminate\Database\Connection as BaseConnection;
use OpenTelemetry\API\Trace\Propagation\TraceContextPropagator;
use Google\GoogleSqlCommenterLaravel\Opentelemetry;


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
        $query .= $this->get_sqlcomments();
        $records = parent::select($query, $bindings, $useReadPdo);

        if (count($records) > 0)
        {
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
        $query .= $this->get_sqlcomments();
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
        $query .= $this->get_sqlcomments();

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
        $query .= $this->get_sqlcomments();

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
        $query .= $this->get_sqlcomments();

        return $this->affectingStatement($query, $bindings);
    }

    public function get_sqlcomments()
    {
        $configurationKey = 'google_sqlcommenter.include.';
        $comment = [];
        if (config($configurationKey.'framework', true)){
            $comment['framework'] = "laravel-" .app()->version();
        }
        if (config($configurationKey.'controller', true)){
            $action = app('request')->route()->getAction();
            $comment['controller'] = class_basename($action['controller']);
        }
        if (config($configurationKey.'route', true)){
            $comment['route'] = request()->getRequestUri();
        }
        if (config($configurationKey.'db_driver', true)){
            $connection = config('database.default');
            $comment['db_driver'] = config("database.connections.{$connection}.driver");
        }
        if (config($configurationKey.'opentelemetry', true)){
            $carrier = Opentelemetry::get_opentelemetry_values();
            $comment = array_merge($comment,$carrier);
        }
        return $this->format_comments($comment);
    }
    public function format_comments($comment)
    {
        if (empty($comment)) {
            return "";
        }
        $lastElement = array_key_last($comment);
        $sql_comment = "/*";
        foreach($comment as $key=>$value)
        {
            if ($key == $lastElement){
                $sql_comment .=$key ."=". "'".$value."'*/";
            }
            else{
                $sql_comment .=$key ."=". "'".$value."',";
            }
        }
        return $sql_comment;
    }

}
