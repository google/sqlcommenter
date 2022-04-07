<?php

namespace Google\GoogleSqlCommenterLaravel;

use OpenTelemetry\API\Trace\Propagation\TraceContextPropagator;

class Opentelemetry
{
    public static function get_opentelemetry_values()
    {
        $carrier = [];

        $trace = TraceContextPropagator::getInstance();
        $trace->inject($carrier);
        return $carrier;
    }
}