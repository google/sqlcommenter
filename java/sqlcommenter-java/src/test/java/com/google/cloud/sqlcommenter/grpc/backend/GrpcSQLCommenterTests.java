package com.google.cloud.sqlcommenter.grpc.backend;

import static org.junit.Assert.assertEquals;
import static org.mockito.Mockito.mock;

import com.google.cloud.sqlcommenter.grpc.backend.service.StudentServiceImpl;
import com.google.cloud.sqlcommenter.grpc.stubs.StudentRequest;
import com.google.cloud.sqlcommenter.grpc.stubs.StudentResponse;
import com.google.cloud.sqlcommenter.threadlocalstorage.State;
import com.google.cloud.sqlcommenter.util.SCHibernateWrapper;
import io.grpc.stub.StreamObserver;
import java.util.List;
import org.junit.Test;

public class GrpcSQLCommenterTests {

  @Test
  public void givenStudentService_whenGetStudentInfoIsCalled_thenSQLStatementsShouldBeTagged() {
    // given
    StudentServiceImpl studentService = new StudentServiceImpl();
    StreamObserver<StudentResponse> observer = mock(StreamObserver.class);
    State.Holder.set(
        State.newBuilder()
            .withControllerName("StudentServiceImpl")
            .withActionName("getStudentInfo")
            .withFramework("grpc")
            .build());

    SCHibernateWrapper.reset();
    StudentRequest studentRequest = StudentRequest.newBuilder().setStudentId("st1").build();

    // when
    studentService.getStudentInfo(studentRequest, observer);
    List<String> sqlStatements = SCHibernateWrapper.getAfterSqlStatements();

    // then
    assertEquals(1, sqlStatements.size());
    assertEquals(
        1,
        sqlStatements
            .stream()
            .filter(
                sql ->
                    sql.contains(
                        "/*action='getStudentInfo',controller='StudentServiceImpl',framework='grpc'*/"))
            .count());

    SCHibernateWrapper.reset();
  }
}
