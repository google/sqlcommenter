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

package io.orijtech.integrations.sqlcommenter.threadlocalstorage;

import io.opencensus.trace.SpanContext;
import io.opencensus.trace.SpanId;
import io.opencensus.trace.TraceId;
import io.opencensus.trace.TraceOptions;
import io.opencensus.trace.Tracestate;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.Map;
import java.util.SortedMap;
import java.util.TreeMap;
import java.util.logging.Level;
import java.util.logging.Logger;
import javax.annotation.Nullable;

public final class State {

  private static final Logger logger = Logger.getLogger(State.class.getName());

  private State() {}

  private final String UTF8 = StandardCharsets.UTF_8.toString();
  public static final String W3C_CONTEXT_VERSION = "00"; // This encodes version 1.

  public static class Builder {

    @Nullable private String actionName;
    @Nullable private String controllerName;
    @Nullable private String framework;
    @Nullable private SpanContext spanContext;

    private Builder() {}

    public Builder withActionName(String actionName) {
      this.actionName = actionName;
      return this;
    }

    public Builder withControllerName(String controllerName) {
      this.controllerName = controllerName;
      return this;
    }

    public Builder withSpanContext(SpanContext spanContext) {
      this.spanContext = spanContext;
      return this;
    }

    public Builder withFramework(String framework) {
      this.framework = framework;
      return this;
    }

    public State build() {
      State state = new State();

      state.actionName = this.actionName;
      state.controllerName = this.controllerName;
      state.spanContext = this.spanContext;
      state.framework = this.framework;

      return state;
    }
  }

  @Nullable private String controllerName;
  @Nullable private String actionName;
  @Nullable private String framework;
  @Nullable private SpanContext spanContext;

  public static Builder newBuilder() {
    return new Builder();
  }

  public static Builder newBuilder(State copy) {
    if (copy == null) { // They are requesting a reset
      return new Builder();
    }

    return new Builder()
        .withActionName(copy.actionName)
        .withControllerName(copy.controllerName)
        .withSpanContext(copy.spanContext)
        .withFramework(copy.framework);
  }

  private Boolean hasSQLComment(String stmt) {
    if (stmt == null || stmt.isEmpty()) {
      return false;
    }

    // Perhaps we now have the closing comment or not but that doesn't matter
    // as the SQL should be reported as having an opening comment regardless
    // of it is in properly escaped or not.
    return stmt.contains("--") || stmt.contains("/*");
  }

  public String formatAndAppendToSQL(String sql) {
    if (sql == null || sql.isEmpty()) {
      return sql;
    }

    // If the SQL already has a comment, just return it.
    if (hasSQLComment(sql)) {
      return sql;
    }

    String commentStr = this.toString();
    // In some cases all the fields might be blank and
    // that produces a blank comment, so return the original
    // SQL statement as it was.
    if (commentStr.isEmpty()) {
      return sql;
    }

    // Otherwise, now insert the fields and format.
    return String.format("%s /*%s*/", sql, commentStr);
  }

  private String urlEncode(String s) throws Exception {
    return URLEncoder.encode(s, UTF8);
  }

  private SortedMap<String, Object> sortedKeyValuePairs() {
    SortedMap<String, Object> sortedMap = new TreeMap<String, Object>();
    sortedMap.put("action", this.actionName);
    sortedMap.put("traceparent", this.traceParent());
    sortedMap.put("tracestate", this.traceState());
    sortedMap.put("controller", this.controllerName);
    sortedMap.put("framework", this.framework);

    return sortedMap;
  }

  @Nullable
  private String traceParent() {
    if (spanContext == null || !spanContext.isValid()) {
      // According to the W3C TraceContext specification, a blank spanId
      // is invalid and should not create a TraceParent.
      return null;
    }
    TraceId traceId = spanContext.getTraceId();
    SpanId spanId = spanContext.getSpanId();

    TraceOptions traceOptions = spanContext.getTraceOptions();
    // A sample:
    //    traceparent='00-a22901f654b534675439f71fbe43783d-7fde95452aa72253-01'
    return String.format(
        "%s-%s-%s-%02X",
        W3C_CONTEXT_VERSION,
        traceId.toLowerBase16(),
        spanId.toLowerBase16(),
        traceOptions.getByte());
  }

  @Nullable
  private String traceState() {
    if (this.spanContext == null || !this.spanContext.isValid()) {
      return null;
    }
    Tracestate traceState = this.spanContext.getTracestate();
    if (traceState.getEntries().isEmpty()) {
      return null;
    }

    // Tracestate needs to be serialized in the order of the entries.
    ArrayList<String> pairsList = new ArrayList<String>();
    for (Tracestate.Entry entry : traceState.getEntries()) {
      try {
        String key = entry.getKey();
        // Only don't insert if the key is empty.
        if (key.isEmpty()) {
          continue;
        }

        String value = entry.getValue();
        String encoded = urlEncode(String.format("%s=%s", key, value));
        pairsList.add(encoded);
      } catch (Exception e) {
        logger.log(Level.WARNING, "Exception when encoding Tracestate", e);
      }
    }

    return String.join(",", pairsList);
  }

  public String toString() {
    SortedMap<String, Object> skvp = this.sortedKeyValuePairs();
    ArrayList<String> keyValuePairsList = new ArrayList<String>();
    // Given that the keys are sorted.
    for (Map.Entry<String, Object> entry : skvp.entrySet()) {
      Object value = entry.getValue();
      if (isBlank(value)) {
        continue;
      }

      // Otherwise, now insert the: url_encode(key)='url_encoded(<value>)'
      try {
        String valueStr = String.format("%s", value);
        String keyValuePairString =
            String.format("%s='%s'", urlEncode(entry.getKey()), urlEncode(valueStr));
        keyValuePairsList.add(keyValuePairString);
      } catch (Exception e) {
        logger.log(Level.WARNING, "Exception when encoding State", e);
      }
    }

    return String.join(",", keyValuePairsList);
  }

  // VisibleForTesting
  private Boolean isBlank(Object obj) {
    if (obj == null) {
      return true;
    }
    if (obj instanceof String) {
      return obj == "";
    }
    if (obj instanceof Number) {
      Number number = (Number) obj;
      return number.doubleValue() == 0.0;
    }
    return false;
  }

  public static class Holder {

    private static final ThreadLocal<State> threadState = new ThreadLocal<State>();

    public static void set(State state) {
      threadState.set(state);
    }

    public static State get() {
      return threadState.get();
    }

    public static void remove() {
      threadState.remove();
    }
  }
}
