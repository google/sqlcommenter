package com.google.cloud.sqlcommenter.threadlocalstorage;

import javax.annotation.Nullable;
import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.logging.Level;
import java.util.logging.Logger;

public class SpanContextMetadata {
    private static final Logger logger = Logger.getLogger(SpanContextMetadata.class.getName());

    private final String traceId;
    private final String spanId;
    private final byte traceOptions;
    private final String traceState;

    private SpanContextMetadata(String traceId, String spanId, byte traceOptions, @Nullable String traceState) {
        this.traceId = traceId;
        this.spanId = spanId;
        this.traceOptions = traceOptions;
        this.traceState = traceState;
    }

    public static SpanContextMetadata fromOpenCensusContext(io.opencensus.trace.SpanContext spanContext) {
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
                try {
                    String key = entry.getKey();
                    // Only don't insert if the key is empty.
                    if (key.isEmpty()) {
                        continue;
                    }

                    String value = entry.getValue();
                    String encoded = URLEncoder.encode((String.format("%s=%s", key, value)), StandardCharsets.UTF_8.toString());
                    pairsList.add(encoded);
                } catch (Exception e) {
                    logger.log(Level.WARNING, "Exception when encoding Tracestate", e);
                }
            }

            traceStateStr = String.join(",", pairsList);
        }

        return new SpanContextMetadata(traceId, spanId, traceOptions, traceStateStr);
    }

    public static SpanContextMetadata fromOpenTelemetryContext(io.opentelemetry.api.trace.SpanContext spanContext) {
        if (spanContext == null || !spanContext.isValid()) {
            return null;
        }

        String traceId = spanContext.getTraceIdAsHexString();
        String spanId = spanContext.getSpanIdAsHexString();
        byte traceOptions = spanContext.getTraceFlags();

        return new SpanContextMetadata(traceId, spanId, traceOptions, null);
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
