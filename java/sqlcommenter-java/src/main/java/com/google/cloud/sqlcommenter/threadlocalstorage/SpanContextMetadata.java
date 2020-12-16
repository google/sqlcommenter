package com.google.cloud.sqlcommenter.threadlocalstorage;

import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.logging.Level;
import java.util.logging.Logger;

public class SpanContextMetadata {
    private static final Logger logger = Logger.getLogger(SpanContextMetadata.class.getName());

    private boolean isValid = true;
    private String traceId;
    private String spanId;
    private byte traceOptions;
    private String traceState;

    public SpanContextMetadata(io.opencensus.trace.SpanContext spanContext) {
        if (spanContext == null || !spanContext.isValid()) {
            this.isValid = false;
            return;
        }
        this.traceId = spanContext.getTraceId().toLowerBase16();
        this.spanId = spanContext.getSpanId().toLowerBase16();
        this.traceOptions = spanContext.getTraceOptions().getByte();

        io.opencensus.trace.Tracestate traceState = spanContext.getTracestate();
        if (traceState.getEntries().isEmpty()) {
            this.traceState = null;
            return;
        }

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

        this.traceState = String.join(",", pairsList);
    }

    public SpanContextMetadata(io.opentelemetry.api.trace.SpanContext spanContext) {
        if (spanContext == null || !spanContext.isValid()) {
            this.isValid = false;
            return;
        }
        this.traceId = spanContext.getTraceIdAsHexString();
        this.spanId = spanContext.getSpanIdAsHexString();
        this.traceOptions = spanContext.getTraceFlags();
    }

    public boolean isValid() {
        return isValid;
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
