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

package com.google.cloud.sqlcommenter.interceptors;

import com.google.cloud.sqlcommenter.threadlocalstorage.State;
import io.grpc.Metadata;
import io.grpc.ServerCall;
import io.grpc.ServerCallHandler;
import io.grpc.ServerInterceptor;

public class GrpcSQLCommenterInterceptor implements ServerInterceptor {
  @Override
  public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(
      ServerCall<ReqT, RespT> call, Metadata requestHeaders, ServerCallHandler<ReqT, RespT> next) {

    if (call == null || next == null) return null;

    if (call != null) {
      String actionName = call.getMethodDescriptor().getBareMethodName();
      String serviceName = call.getMethodDescriptor().getServiceName();
      State.Holder.set(
          State.newBuilder()
              .withControllerName(serviceName)
              .withActionName(actionName)
              .withFramework("spring-grpc")
              .build());
    }

    return next.startCall(call, requestHeaders);
  }
}
