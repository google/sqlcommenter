// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

const {trace, TraceFlags} = require('@opentelemetry/api');
const {DEFAULT_INSTRUMENTATION_PLUGINS} = require('@opentelemetry/node/build/src/config');

exports.OpenTelemetryProvider = class OpenTelemetryProvider {
    getW3CTraceContext() {
        let spanContext = null;
        // Collect a list of the default loaded nodejs tracers.
        // The tracer names are required to retrieve the latest span.
        let validPluginNames = Object.keys(DEFAULT_INSTRUMENTATION_PLUGINS)
            .map((pluginName) => DEFAULT_INSTRUMENTATION_PLUGINS[pluginName])
            .filter((pluginConfig) => pluginConfig && pluginConfig.enabled)
            .map((pluginConfig) => pluginConfig.path);
        validPluginNames.push('default');
        validPluginNames.forEach((pluginPath) => {
            let span = trace.getTracerProvider().getTracer(pluginPath).getCurrentSpan();
            if (span) {
                spanContext = span.context();
            }
        });

        const traceParent = `00-${spanContext.traceId}-${spanContext.spanId}-0${
            Number(spanContext.traceFlags || TraceFlags.NONE).toString(16)}`;

        return { traceparent: traceParent };
    }
};
