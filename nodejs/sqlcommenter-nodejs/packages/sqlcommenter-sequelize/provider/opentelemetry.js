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

const {context, defaultSetter} = require('@opentelemetry/api');
const {HttpTraceContext} = require('@opentelemetry/core')

exports.OpenTelemetry = class OpenTelemetry {
    static addW3CTraceContext(comments) {
        let propagator = new HttpTraceContext();
        propagator.inject(context.active(), comments, defaultSetter);
    }
};
