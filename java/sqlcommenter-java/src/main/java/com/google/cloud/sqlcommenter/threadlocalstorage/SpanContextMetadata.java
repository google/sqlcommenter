package com.google.cloud.sqlcommenter.threadlocalstorage;

import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.Map;
import java.util.logging.Level;
import java.util.logging.Logger;
import javax.annotation.Nullable;

public class SpanContextMetadata {
  private static final Logger logger = Logger.getLogger(SpanContextMetadata.class.getName());
  private static final String UTF8 = StandardCharsets.UTF_8.toString();

  private final String traceId;
  private final String spanId;
  private final byte traceOptions;
  private final String traceState;

  private SpanContextMetadata(
      String traceId, String spanId, byte traceOptions, @Nullable String traceState) {
    this.traceId = traceId;
    this.spanId = spanId;
    this.traceOptions = traceOptions;
    this.traceState = traceState;
  }

  public static SpanContextMetadata fromOpenCensusContext(
      io.opencensus.trace.SpanContext spanContext) {
    if (spanContext == null || !spanContext.isValid()) {
      return null;
    }
    String traceId = spanContext.getTraceId().toLowerBase16();
    String spanId = spanContext.getSpanId().toLowerBase16();
    byte traceOptions = spanContext.getTraceOptions().getByte();
    String traceStateStr = null;

    io.opencensus.trace.Tracestate traceState = spanContext.getTracestate();
    if (!traceState.getEntries().isEmpty()) {
      // Tracestate needs to be serialized in the order of the entries.
      ArrayList<String> pairsList = new ArrayList<>();
      for (io.opencensus.trace.Tracestate.Entry entry : traceState.getEntries()) {
        String key = entry.getKey();
        // Only don't insert if the key is empty.
        if (!key.isEmpty()) {
          try {
            String value = entry.getValue();
            String encoded = URLEncoder.encode((String.format("%s=%s", key, value)), UTF8);
            pairsList.add(encoded);
          } catch (Exception e) {
            logger.log(Level.WARNING, "Exception when encoding Tracestate", e);
          }
        }
      }

      traceStateStr = String.join(",", pairsList);
    }

    return new SpanContextMetadata(traceId, spanId, traceOptions, traceStateStr);
  }

  public static SpanContextMetadata fromOpenTelemetryContext(
      io.opentelemetry.api.trace.SpanContext spanContext) {
    if (spanContext == null || !spanContext.isValid()) {
      return null;
    }

    String traceId = spanContext.getTraceIdAsHexString();
    String spanId = spanContext.getSpanIdAsHexString();
    byte traceOptions = spanContext.getTraceFlags();
    String traceStateStr = null;

    io.opentelemetry.api.trace.TraceState traceState = spanContext.getTraceState();
    if (!traceState.isEmpty()) {
      ArrayList<String> pairsList = new ArrayList<>();

      Map<String, String> map = traceState.asMap();
      for (String key : map.keySet()) {
        if (key.isEmpty()) {
          continue;
        }

        try {
          String value = map.get(key);
          String encoded = URLEncoder.encode((String.format("%s=%s", key, value)), UTF8);
        } catch (Exception e) {
          logger.log(Level.WARNING, "Exception when encoding Tracestate", e);
        }
      }
      traceStateStr = String.join(",", pairsList);
    }

    return new SpanContextMetadata(traceId, spanId, traceOptions, traceStateStr);
  }

  public String getTraceId() {
    return traceId;
  }

  public String getSpanId() {
    return spanId;
  }

  public byte getTraceOptions() {
    return traceOptions;
  }

  public String getTraceState() {
    return traceState;
  }
}
