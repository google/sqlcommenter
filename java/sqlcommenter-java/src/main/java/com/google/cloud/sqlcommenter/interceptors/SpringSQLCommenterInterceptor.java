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
import javax.servlet.http.HttpServletRequest;
import javax.servlet.http.HttpServletResponse;
import org.springframework.web.method.HandlerMethod;
import org.springframework.web.servlet.handler.HandlerInterceptorAdapter;

public class SpringSQLCommenterInterceptor extends HandlerInterceptorAdapter {

  @Override
  public boolean preHandle(HttpServletRequest request, HttpServletResponse response, Object handler)
      throws Exception {
    // This method MUST always return true since we are
    // only grabbing information about the request.

    boolean isHandlerMethod = handler instanceof HandlerMethod;
    if (!isHandlerMethod) {
      // In this case, return promptly since we
      // can't extract details from the handler.
      return true;
    }

    HandlerMethod handlerMethod = (HandlerMethod) handler;
    String actionName = handlerMethod.getMethod().getName();
    String controllerName = handlerMethod.getBeanType().getSimpleName().replace("Controller", "");

    State.Holder.set(
        State.newBuilder()
            .withControllerName(controllerName)
            .withActionName(actionName)
            .withFramework("spring")
            .build());

    return true;
  }
}
