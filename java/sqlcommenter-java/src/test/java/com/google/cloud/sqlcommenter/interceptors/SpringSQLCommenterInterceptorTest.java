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

import static com.google.common.truth.Truth.assertThat;

import com.google.cloud.sqlcommenter.threadlocalstorage.State;
import org.junit.After;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;
import org.springframework.web.method.HandlerMethod;

/**
 * Tests for {@link
 * com.google.cloud.sqlcommenter.interceptors.SpringSQLCommenterInterceptor}.
 */
@RunWith(JUnit4.class)
public class SpringSQLCommenterInterceptorTest {

  @Before
  public void setUp() {
    State.Holder.remove();
  }

  @After
  public void tearDown() {
    State.Holder.remove();
  }

  private static class fakeBean {
    public boolean methodUno() {
      return true;
    }

    public String methodDos() {
      return "dos";
    }

    public String methodThreaded() {
      return "threaded";
    }
  }

  @Test
  public void testPreHandlePlacesStateInThreadLocal() throws Exception {
    // 0. Precursor: no state should have been set.
    assertThat(State.Holder.get()).isEqualTo(null);

    SpringSQLCommenterInterceptor springsci = new SpringSQLCommenterInterceptor();

    // 1. A handler that isn't an instance of HandlerMethod should ALWAYS return true,
    // but not set the state.
    boolean ok = springsci.preHandle(null, null, true);
    assertThat(ok).isEqualTo(true);

    // 2. A handler that is an instance of HandlerMethod should ALWAYS return true,
    // and also set the thread local state.
    fakeBean bean = new fakeBean();
    HandlerMethod hm = new HandlerMethod(bean, "methodUno");
    ok = springsci.preHandle(null, null, hm);
    assertThat(ok).isEqualTo(true);

    // 2.2. Ensure that we can retrieve the newly inserted state from threadlocal storage.
    State state = State.Holder.get();
    assertThat(state).isNotEqualTo(null);
    assertThat(state.toString())
        .isEqualTo("action='methodUno',controller='fakeBean',framework='spring'");
    // 2.3. Now with SQL that it is formatted alright.
    assertThat(state.formatAndAppendToSQL("SELECT * from FOO"))
        .isEqualTo(
            "SELECT * from FOO /*action='methodUno',controller='fakeBean',framework='spring'*/");

    // 3.0. On a subsequent call, the state should be over-written in the same thread.
    hm = new HandlerMethod(bean, "methodDos");
    ok = springsci.preHandle(null, null, hm);
    assertThat(ok).isEqualTo(true);

    state = State.Holder.get();
    assertThat(state).isNotEqualTo(null);
    assertThat(state.toString())
        .isEqualTo("action='methodDos',controller='fakeBean',framework='spring'");
    // 3.1. Now with SQL that it is formatted alright.
    assertThat(state.formatAndAppendToSQL("SELECT * from FOO"))
        .isEqualTo(
            "SELECT * from FOO /*action='methodDos',controller='fakeBean',framework='spring'*/");

    // 4.0 Ensure that in a separate thread the state doesn't pre-exist
    Thread th2 =
        new Thread(
            () -> {
              State state2 = State.Holder.get();
              // 4.1. Check that in a separate thread we start with a null State.
              assertThat(state2).isEqualTo(null);

              try {
                HandlerMethod hm2 = new HandlerMethod(bean, "methodThreaded");
                boolean ok2 = springsci.preHandle(null, null, hm2);
                assertThat(ok2).isEqualTo(true);

                state2 = State.Holder.get();
                assertThat(state2.toString())
                    .isEqualTo(
                        "action='methodThreaded',controller='fakeBean',framework='spring'");
                // 3.1. Now with SQL that it is formatted alright.
                assertThat(state2.formatAndAppendToSQL("SELECT * from FOO"))
                    .isEqualTo(
                        "SELECT * from FOO /*action='methodThreaded',controller='fakeBean',framework='spring'*/");
              } catch (Exception e) {
              }
            });
    th2.start();
    th2.join();

    // 3.2. Ensure that the previous state didn't interleave with the current one,
    // thus they are separated by being in different thread local storage.
    // So that in 3.2. our state is as it was in 3.1.
    assertThat(state.formatAndAppendToSQL("SELECT * from FOO"))
        .isEqualTo(
            "SELECT * from FOO /*action='methodDos',controller='fakeBean',framework='spring'*/");
  }
}
