// Copyright 2020 Google LLC
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

const OpenTelemetry = require('./opentelemetry');

const providers = {
    'opentelemetry': OpenTelemetry
}

exports.attachComments = function attachComments(providerName, comments) {
    // Verify we have a comments object to modify
    if (!comments || typeof comments !== 'object') return;

    // Lookup the provider by name, or use the default.
    let provider = providers[String(providerName).toLowerCase()] || OpenTelemetry;
    provider.addW3CTraceContext(comments);
}
