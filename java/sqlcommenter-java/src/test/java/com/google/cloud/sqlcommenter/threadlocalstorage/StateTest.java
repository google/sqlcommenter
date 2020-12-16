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

package com.google.cloud.sqlcommenter.threadlocalstorage;

import static com.google.common.truth.Truth.assertThat;

import io.opencensus.trace.SpanContext;
import io.opencensus.trace.SpanId;
import io.opencensus.trace.TraceId;
import io.opencensus.trace.TraceOptions;
import io.opencensus.trace.Tracestate;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;

/**
 * Tests for {@link State}.
 */
@RunWith(JUnit4.class)
public class StateTest {

  private static final Byte byteSampled = (byte) 0x01;
  private static final Byte byteNotSampled = (byte) 0x00;

  @Before
  public void setUp() {
    State.Holder.remove();
  }

  @After
  public void tearDown() {
    State.Holder.remove();
  }

  @Test
  public void testFormatAndAppendToSQL() {
    State state =
            State.newBuilder()
                    .withFramework("spring")
                    .withControllerName("foo;DROP TABLE BAR")
                    .withActionName("run this & that")
                    .withSpanContextMetadata(
                            new SpanContextMetadata(SpanContext.create(
                                    TraceId.fromLowerBase16("9a4589fe88dd0fc9ffee11228888ff11"),
                                    SpanId.fromLowerBase16("11fa8b009a4589fe"),
                                    TraceOptions.fromByte(byteSampled))))
                    .build();

    // 1. Assert that proper comments are generated.
    assertThat(state.formatAndAppendToSQL("SELECT * from USERS"))
            .isEqualTo(
                    // Java's url_encode encodes " " as "+" instead of "%20".
                    "SELECT * from USERS /*action='run+this+%26+that',controller='foo%3BDROP+TABLE+BAR',framework='spring',traceparent='00-9a4589fe88dd0fc9ffee11228888ff11-11fa8b009a4589fe-01'*/");

    // 2. Assert that passing in null returns null.
    assertThat(state.formatAndAppendToSQL(null)).isEqualTo(null);

    // 3.1. Assert that passing in a SQL statement that has comments /*...*/ will return the
    // original, untampered.
    assertThat(state.formatAndAppendToSQL("SELECT * from USERS /*this is a pre-existing comment*/"))
            .isEqualTo("SELECT * from USERS /*this is a pre-existing comment*/");

    // 3.2. Assert that passing in a SQL statement that has comments -- will return the original,
    // untampered.
    assertThat(state.formatAndAppendToSQL("SELECT * from USERS -- this is a pre-existing comment"))
            .isEqualTo("SELECT * from USERS -- this is a pre-existing comment");

    // 3.3. Assert that passing in a SQL statement that has comments-- will return the original,
    // untampered.
    assertThat(
            state.formatAndAppendToSQL(
                    "SELECT * from USERS-- this is a pre-existing but malformed comment"))
            .isEqualTo("SELECT * from USERS-- this is a pre-existing but malformed comment");

    // 4.1. Assert that passing in the empty string, returns as is.
    assertThat(state.formatAndAppendToSQL("")).isEqualTo("");

    // 4.2. Assert that passing in null, returns null as is.
    assertThat(state.formatAndAppendToSQL(null)).isEqualTo(null);
  }

  @Test
  public void testFormatAndAppendToSQLNotSampled() {
    State state =
            State.newBuilder()
                    .withFramework("spring")
                    .withControllerName("foo;DROP TABLE BAR")
                    .withActionName("run this & that")
                    .withSpanContextMetadata(
                            new SpanContextMetadata(SpanContext.create(
                                    TraceId.fromLowerBase16("9a4589fe88dd0fc911ff2233ffee7899"),
                                    SpanId.fromLowerBase16("11fa8b00dd11eeff"),
                                    TraceOptions.fromByte(byteNotSampled),
                                    Tracestate.builder()
                                            .build()
                                            .toBuilder()
                                            // A new entry will always be added in the front of the list of entries.
                                            .set("congo", "t61rcWkgMzE")
                                            .set("rojo", "00f067aa0ba902b7")
                                            .build())))
                    .build();

    // 1. Assert that proper comments are generated.
    assertThat(state.formatAndAppendToSQL("SELECT * from USERS"))
            .isEqualTo(
                    // Java's url_encode encodes " " as "+" instead of "%20".
                    "SELECT * from USERS /*action='run+this+%26+that',controller='foo%3BDROP+TABLE+BAR',framework='spring',traceparent='00-9a4589fe88dd0fc911ff2233ffee7899-11fa8b00dd11eeff-00',tracestate='rojo%253D00f067aa0ba902b7%2Ccongo%253Dt61rcWkgMzE'*/");

    // 2. Assert that passing in null returns null.
    assertThat(state.formatAndAppendToSQL(null)).isEqualTo(null);

    // 3.1. Assert that passing in a SQL statement that has comments /*...*/ will return the
    // original, untampered.
    assertThat(state.formatAndAppendToSQL("SELECT * from USERS /*this is a pre-existing comment*/"))
            .isEqualTo("SELECT * from USERS /*this is a pre-existing comment*/");

    // 3.2. Assert that passing in a SQL statement that has comments -- will return the original,
    // untampered.
    assertThat(state.formatAndAppendToSQL("SELECT * from USERS -- this is a pre-existing comment"))
            .isEqualTo("SELECT * from USERS -- this is a pre-existing comment");

    // 3.3. Assert that passing in a SQL statement that has comments-- will return the original,
    // untampered.
    assertThat(
            state.formatAndAppendToSQL(
                    "SELECT * from USERS-- this is a pre-existing but malformed comment"))
            .isEqualTo("SELECT * from USERS-- this is a pre-existing but malformed comment");

    // 4.1. Assert that passing in the empty string, returns as is.
    assertThat(state.formatAndAppendToSQL("")).isEqualTo("");

    // 4.2. Assert that passing in null, returns null as is.
    assertThat(state.formatAndAppendToSQL(null)).isEqualTo(null);
  }

  @Test
  public void testBuilderWithNullState() {
    // Ensure that passing in null to newBuilder doesn't crash.
    State sNull = State.newBuilder(null).build();
    assertThat(sNull).isNotEqualTo(null);

    State state =
            State.newBuilder()
                    .withFramework("spring")
                    .withControllerName("foo;DROP TABLE BAR")
                    .withActionName("run this & that")
                    .withSpanContextMetadata(
                            new SpanContextMetadata(
                                    SpanContext.create(
                                            TraceId.fromLowerBase16("9a4589fe88dd0fc9ffdd11eedd2233ff"),
                                            SpanId.fromLowerBase16("11fa8b00cc23114f"),
                                            TraceOptions.fromByte(byteSampled))))
                    .build();
    assertThat(state).isNotEqualTo(null);

    // 1. Reset the state with null, we shouldn't crash.
    state = State.newBuilder(null).build();
    assertThat(state).isNotEqualTo(null);

    // 2. Assert that passing in null returns null.
    assertThat(state.formatAndAppendToSQL(null)).isEqualTo(null);
  }
}
