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
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;

/**
 * Tests for {@link State}.
 */
@RunWith(JUnit4.class)
public class ThreadLocalStorageTest {

  private static final Byte byteSampled = (byte) 0x01;

  @Before
  public void setUp() {
    State.Holder.remove();
  }

  @After
  public void tearDown() {
    State.Holder.remove();
  }

  @Test
  public void testGet() throws Exception {
    // 1. No state in this thread.
    State withNonPreviously = State.Holder.get();

    assertThat(withNonPreviously).isEqualTo(null);

    // 2. Insert some state into thread local storage.
    State inTh1 =
            State.newBuilder()
                    .withControllerName("foo;DROP TABLE BAR")
                    .withActionName("run this & that")
                    .withSpanContextMetadata(
                            SpanContextMetadata.fromOpenCensusContext(SpanContext.create(
                                    TraceId.fromLowerBase16("9a4589fe88dd0fc9ffee11228888ff11"),
                                    SpanId.fromLowerBase16("11fa8b00ff221eec"),
                                    TraceOptions.fromByte(byteSampled))))
                    .build();

    State.Holder.set(inTh1);

    State got = State.Holder.get();
    assertThat(inTh1).isEqualTo(got);

    // 3. Delete it.
    State.Holder.remove();

    // 3.1. Now check that what's stored doesn't exist anymore.
    State gotAfterRemove = State.Holder.get();
    assertThat(gotAfterRemove).isEqualTo(null);

    assertThat(gotAfterRemove).isNotEqualTo(got);
    assertThat(gotAfterRemove).isNotEqualTo(inTh1);

    // 4. Now insert it into this current thread,
    // then try to retrieve it from a different thread.
    State.Holder.set(inTh1);

    State inTh2 =
            State.newBuilder()
                    .withControllerName("foo")
                    .withActionName("action")
                    .withSpanContextMetadata(
                            SpanContextMetadata.fromOpenCensusContext(SpanContext.create(
                                    TraceId.fromLowerBase16("aa4589fe88dd0faae1f2d3c4dd11f344"),
                                    SpanId.fromLowerBase16("91ea8891ff221eec"),
                                    TraceOptions.fromByte(byteSampled))))
                    .build();

    assertThat(inTh2).isNotEqualTo(inTh1);

    State.Holder.set(inTh1);

    Thread th2 =
            new Thread(
                    () -> {
                      State stateInOtherThread = State.Holder.get();
                      // 4.1. Check that in a separate thread we start with a null State.
                      assertThat(stateInOtherThread).isEqualTo(null);

                      State.Holder.set(inTh2);
                      State got1 = State.Holder.get();

                      // 4.2. After retrieval, make that insertion in thread 2.
                      assertThat(got1).isEqualTo(inTh2);
                    });
    th2.start();
    th2.join();

    // 4.3. Immediately after check that inTh2 can't be accessed from a different thread.
    State gotAfter = State.Holder.get();
    assertThat(gotAfter).isNotEqualTo(inTh2);
    // 4.4. Assert that what's in this current thread is what we had previously inserted.
    assertThat(gotAfter).isEqualTo(inTh1);
  }
}
