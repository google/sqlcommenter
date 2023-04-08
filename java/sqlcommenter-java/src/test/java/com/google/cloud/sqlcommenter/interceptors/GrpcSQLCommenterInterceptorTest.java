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
import static io.grpc.MethodDescriptor.generateFullMethodName;
import static org.mockito.Mockito.when;

import com.google.cloud.sqlcommenter.threadlocalstorage.State;
import io.grpc.*;
import io.grpc.MethodDescriptor.Marshaller;
import io.grpc.MethodDescriptor.MethodType;
import org.junit.After;
import org.junit.Before;
import org.junit.Rule;
import org.junit.Test;
import org.junit.rules.ExpectedException;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;
import org.mockito.ArgumentMatchers;
import org.mockito.Mock;
import org.mockito.Mockito;
import org.mockito.junit.MockitoJUnit;
import org.mockito.junit.MockitoRule;

@RunWith(JUnit4.class)
public class GrpcSQLCommenterInterceptorTest {

  @Rule public final MockitoRule mocks = MockitoJUnit.rule();

  @SuppressWarnings("deprecation") // https://github.com/grpc/grpc-java/issues/7467
  @Rule
  public final ExpectedException thrown = ExpectedException.none();

  @Mock private Marshaller<String> requestMarshaller;

  @Mock private Marshaller<Integer> responseMarshaller;

  @Mock private ServerCallHandler<String, Integer> handler;

  @Mock private ServerCall.Listener<String> listener;

  private MethodDescriptor<String, Integer> flowMethod;

  private ServerServiceDefinition serviceDefinition;

  private final Metadata headers = new Metadata();

  @Before
  public void setUp() {
    State.Holder.remove();

    flowMethod =
        MethodDescriptor.<String, Integer>newBuilder()
            .setType(MethodType.UNKNOWN)
            .setFullMethodName(
                generateFullMethodName(
                    "com.google.cloud.sqlcommenter.grpc.backend.service.StudentServiceImpl",
                    "getStudentInfo"))
            .setRequestMarshaller(requestMarshaller)
            .setResponseMarshaller(responseMarshaller)
            .build();

    Mockito.when(
            handler.startCall(
                ArgumentMatchers.<ServerCall<String, Integer>>any(),
                ArgumentMatchers.<Metadata>any()))
        .thenReturn(listener);

    serviceDefinition =
        ServerServiceDefinition.builder(
                new ServiceDescriptor(
                    "com.google.cloud.sqlcommenter.grpc.backend.service.StudentServiceImpl",
                    flowMethod))
            .addMethod(flowMethod, handler)
            .build();
  }

  @After
  public void tearDown() {
    State.Holder.remove();
  }

  @Test
  public void testPreHandlePlacesStateInThreadLocal() throws Exception {
    // 0. Precursor: no state should have been set.
    assertThat(State.Holder.get()).isEqualTo(null);

    GrpcSQLCommenterInterceptor grpcsci = new GrpcSQLCommenterInterceptor();

    // 1. A handler that isn't an instance of HandlerMethod should ALWAYS return true,
    // but not set the state.
    grpcsci.interceptCall(null, null, null);

    ServerCall<String, Integer> serverCall = Mockito.mock(ServerCall.class);
    when(serverCall.getMethodDescriptor()).thenReturn(flowMethod);

    // 2. A handler that is an instance of HandlerMethod should ALWAYS return true,
    // and also set the thread local state.
    grpcsci.interceptCall(serverCall, null, null);
    // assertThat(ok).isEqualTo(null);

    // 2.2. Ensure that we can retrieve the newly inserted state from threadlocal storage.
    ServerCallHandler<String, Integer> serverCallHandler = Mockito.mock(ServerCallHandler.class);
    // when(serverCallHandler.startCall(any(), any())).thenReturn(ok);
    grpcsci.interceptCall(serverCall, null, serverCallHandler);
    State state = State.Holder.get();
    assertThat(state).isNotEqualTo(null);
    assertThat(state.toString())
        .isEqualTo(
            "action='getStudentInfo',controller='com.google.cloud.sqlcommenter.grpc.backend.service.StudentServiceImpl',framework='spring-grpc'");
    // 2.3. Now with SQL that it is formatted alright.
    assertThat(state.formatAndAppendToSQL("SELECT * from FOO"))
        .isEqualTo(
            "SELECT * from FOO /*action='getStudentInfo',controller='com.google.cloud.sqlcommenter.grpc.backend.service.StudentServiceImpl',framework='spring-grpc'*/");

    // 3.0. On a subsequent call, the state should be over-written in the same thread.
    flowMethod =
        MethodDescriptor.<String, Integer>newBuilder()
            .setType(MethodType.UNKNOWN)
            .setFullMethodName(
                generateFullMethodName(
                    "com.google.cloud.sqlcommenter.grpc.backend.service.TeacherServiceImpl",
                    "getTeacherInfo"))
            .setRequestMarshaller(requestMarshaller)
            .setResponseMarshaller(responseMarshaller)
            .build();
    when(serverCall.getMethodDescriptor()).thenReturn(flowMethod);
    grpcsci.interceptCall(serverCall, null, serverCallHandler);
    // assertThat(ok).isEqualTo(true);

    state = State.Holder.get();
    assertThat(state).isNotEqualTo(null);
    assertThat(state.toString())
        .isEqualTo(
            "action='getTeacherInfo',controller='com.google.cloud.sqlcommenter.grpc.backend.service.TeacherServiceImpl',framework='spring-grpc'");
    // 3.1. Now with SQL that it is formatted alright.
    assertThat(state.formatAndAppendToSQL("SELECT * from FOO"))
        .isEqualTo(
            "SELECT * from FOO /*action='getTeacherInfo',controller='com.google.cloud.sqlcommenter.grpc.backend.service.TeacherServiceImpl',framework='spring-grpc'*/");

    // 4.0 Ensure that in a separate thread the state doesn't pre-exist
    Thread th2 =
        new Thread(
            () -> {
              State state2 = State.Holder.get();
              // 4.1. Check that in a separate thread we start with a null State.
              assertThat(state2).isEqualTo(null);

              try {
                MethodDescriptor<String, Integer> flowMethod2 =
                    MethodDescriptor.<String, Integer>newBuilder()
                        .setType(MethodType.UNKNOWN)
                        .setFullMethodName(
                            generateFullMethodName(
                                "com.google.cloud.sqlcommenter.grpc.backend.service.TeacherServiceImplThreaded",
                                "getTeacherInfoThreaded"))
                        .setRequestMarshaller(requestMarshaller)
                        .setResponseMarshaller(responseMarshaller)
                        .build();
                ServerCallHandler<String, Integer> serverCallHandler2 =
                    Mockito.mock(ServerCallHandler.class);

                ServerCall<String, Integer> serverCall2 = Mockito.mock(ServerCall.class);
                when(serverCall2.getMethodDescriptor()).thenReturn(flowMethod2);

                ServerCall.Listener<String> ok2 =
                    grpcsci.interceptCall(serverCall2, null, serverCallHandler2);

                state2 = State.Holder.get();
                assertThat(state2.toString())
                    .isEqualTo(
                        "action='getTeacherInfoThreaded',controller='com.google.cloud.sqlcommenter.grpc.backend.service.TeacherServiceImplThreaded',framework='spring-grpc'");
                // 3.1. Now with SQL that it is formatted alright.
                assertThat(state2.formatAndAppendToSQL("SELECT * from FOO"))
                    .isEqualTo(
                        "SELECT * from FOO /*action='getTeacherInfoThreaded',controller='com.google.cloud.sqlcommenter.grpc.backend.service.TeacherServiceImplThreaded',framework='spring-grpc'*/");
              } catch (Exception e) {
                e.printStackTrace();
              }
            });
    th2.start();
    th2.join();

    // 3.2. Ensure that the previous state didn't interleave with the current one,
    // thus they are separated by being in different thread local storage.
    // So that in 3.2. our state is as it was in 3.1.
    assertThat(state.formatAndAppendToSQL("SELECT * from FOO"))
        .isEqualTo(
            "SELECT * from FOO /*action='getTeacherInfo',controller='com.google.cloud.sqlcommenter.grpc.backend.service.TeacherServiceImpl',framework='spring-grpc'*/");
  }
}
