package com.google.cloud.sqlcommenter.r2dbc;

import static org.junit.Assert.assertEquals;

import com.google.cloud.sqlcommenter.threadlocalstorage.State;
import io.r2dbc.spi.Connection;
import io.r2dbc.spi.ConnectionFactory;
import io.r2dbc.spi.Statement;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.reactivestreams.Publisher;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.annotation.DirtiesContext;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.junit4.SpringJUnit4ClassRunner;
import org.springframework.test.util.ReflectionTestUtils;
import reactor.core.publisher.Mono;
import reactor.util.context.Context;

@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration(classes = R2DBCConfiguration.class)
@DirtiesContext(classMode = DirtiesContext.ClassMode.AFTER_EACH_TEST_METHOD)
public class ConnectionFactoryAspectTest {

  private final String stmt1 = "SELECT * from FOO";

  @Autowired private ConnectionFactory connectionFactory;

  @Test
  public void name() {
    State state =
        State.newBuilder()
            .withControllerName("Order")
            .withFramework("spring")
            .withActionName("add")
            .build();
    Context updatedContext = Context.of("state", state);

    Publisher<? extends Connection> conn = connectionFactory.create();

    Statement statement =
        Mono.from(conn)
            .contextWrite(updatedContext)
            .map(connection -> connection.createStatement(stmt1))
            .block();

    String value = (String) ReflectionTestUtils.getField(statement, null, "sql");
    assertEquals(value, state.formatAndAppendToSQL(stmt1));
  }
}
