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

package com.google.cloud.sqlcommenter.schibernate;

import static com.google.common.truth.Truth.assertThat;

import com.google.cloud.sqlcommenter.threadlocalstorage.State;

import io.opencensus.trace.samplers.Samplers;
import io.opentelemetry.sdk.trace.TracerSdkProvider;
import org.junit.*;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;

/**
 * Tests for {@link SCHibernateTest}.
 */
@RunWith(JUnit4.class)
public class SCHibernateTest {
  private final String stmt1 = "SELECT * from FOO";
  private final SCHibernate sch = new SCHibernate();
  private final State state1 =
          State.newBuilder()
                  .withFramework("jetty")
                  .withControllerName("baz")
                  .withActionName("may")
                  .build();

  @BeforeClass
  public static void setUp() {
    State.Holder.remove();
  }

  @After
  public void tearDown() {
    State.Holder.remove();
  }

  @Test
  public void testWithoutState() {
    String got1 = sch.inspect(stmt1);

    // 1. Since we don't have any data nor state inside
    // the current thread local storage, we should get back
    // the original statement as is.
    assertThat(got1).isEqualTo(stmt1);
  }

  @Test
  public void testWithBasicState() {
    State.Holder.set(state1);

    String got2 = sch.inspect(stmt1);
    assertThat(got2)
            .isEqualTo("SELECT * from FOO /*action='may',controller='baz',framework='jetty'*/");
  }

  @Test
  public void testWithOpenCensusContext() {
    State.Holder.set(state1);

    io.opencensus.trace.Tracer tracer = io.opencensus.trace.Tracing.getTracer();
    // 2. Now insert a span and assert that the SQL has that OpenCensus Trace information.
    try (io.opencensus.common.Scope ss = tracer.spanBuilder("TestSpan").setSampler(Samplers.alwaysSample()).startScopedSpan()) {
      io.opencensus.trace.SpanContext spanContext = tracer.getCurrentSpan().getContext();

      // With that span, now try generating the SQL again.
      String got3 = sch.inspect(stmt1);
      String want3 =
              String.format(
                      "SELECT * from FOO /*action='may',controller='baz',framework='jetty',traceparent='%s-%s-%s-%02X'*/",
                      State.W3C_CONTEXT_VERSION,
                      spanContext.getTraceId().toLowerBase16(),
                      spanContext.getSpanId().toLowerBase16(),
                      spanContext.getTraceOptions().getByte());

      assertThat(got3).isEqualTo(want3);
    }

    // Now that that that scope has ended with resources,
    // assert that the SQL is the same as before.
    State.Holder.set(state1);
    String got3 = sch.inspect(stmt1);
    assertThat(got3)
            .isEqualTo("SELECT * from FOO /*action='may',controller='baz',framework='jetty'*/");
  }

  @Test
  public void testWithOpenTelemetryContext() {
    State.Holder.set(state1);

    io.opentelemetry.api.trace.Tracer tracer = TracerSdkProvider.builder().build().get("SCHibernateTest");
    // 2. Now insert a span and assert that the SQL has that OpenTelemetry Trace information.
    io.opentelemetry.api.trace.Span span = tracer.spanBuilder("TestSpan").startSpan();
    try (io.opentelemetry.context.Scope ss = span.makeCurrent()) {
      io.opentelemetry.api.trace.SpanContext spanContext = span.getSpanContext();

      // With that span, now try generating the SQL again.
      String got = sch.inspect(stmt1);
      String want3 = String.format("SELECT * from FOO /*action='may',controller='baz',framework='jetty',traceparent='%s-%s-%s-%02X'*/",
              State.W3C_CONTEXT_VERSION,
              spanContext.getTraceIdAsHexString(),
              spanContext.getSpanIdAsHexString(),
              spanContext.getTraceFlags());

      assertThat(got).isEqualTo(want3);
    } finally {
      span.end();
    }

    // Now that that that scope has ended with resources,
    // assert that the SQL is the same as before.
    State.Holder.set(state1);
    String got3 = sch.inspect(stmt1);
    assertThat(got3)
            .isEqualTo("SELECT * from FOO /*action='may',controller='baz',framework='jetty'*/");
  }
}
