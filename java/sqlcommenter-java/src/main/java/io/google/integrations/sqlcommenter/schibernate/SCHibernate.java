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

package io.orijtech.integrations.sqlcommenter.schibernate;

import io.opencensus.trace.SpanContext;
import io.opencensus.trace.Tracer;
import io.opencensus.trace.Tracing;
import io.orijtech.integrations.sqlcommenter.threadlocalstorage.State;
import org.hibernate.resource.jdbc.spi.StatementInspector;

public class SCHibernate implements StatementInspector {

  private static final Tracer tracer = Tracing.getTracer();

  /**
   * inspect augments SQL with statements about the current code setup if any. It tries to check if
   * the current context contains OpenCensus information and if so, will augment the SQL but also it
   * will check if any web frameworks have inserted information about the currently running code
   * into thread local storage.
   */
  public String inspect(String sql) {
    State state = State.Holder.get();

    // Otherwise it is time to augment the SQL with information about the controller.
    SpanContext spanContext = tracer.getCurrentSpan().getContext();
    if (spanContext.isValid() && spanContext.getTraceOptions().isSampled()) {
      // Since our goal at this point is NOT to mutate the threadlocal storage's state.
      // yet we still need to carefully capture the Span information, we'll make
      // a copy of the thread local storage but set the TraceId and SpanId.

      // Finally replace this mvcState but it'll just be local to this function
      // and not the final ThreadLocalStorage State.
      state = State.newBuilder(state).withSpanContext(spanContext).build();
    }

    if (state == null) return sql;

    return state.formatAndAppendToSQL(sql);
  }
}
