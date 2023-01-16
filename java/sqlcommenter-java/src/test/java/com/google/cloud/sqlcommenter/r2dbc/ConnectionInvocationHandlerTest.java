package com.google.cloud.sqlcommenter.r2dbc;

import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

import com.google.cloud.sqlcommenter.threadlocalstorage.State;
import io.r2dbc.spi.Connection;
import io.r2dbc.spi.Statement;
import java.lang.reflect.Proxy;
import junit.framework.TestCase;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;
import reactor.util.context.Context;

@RunWith(JUnit4.class)
public class ConnectionInvocationHandlerTest extends TestCase {

  private final String stmt1 = "SELECT * from FOO";

  private Connection connection = mock(Connection.class);
  private Connection proxiedConnection;
  private Context context = Context.empty();
  private ConnectionInvocationHandler connectionInvocationHandler;

  @Override
  @Before
  public void setUp() throws Exception {
    connection = mock(Connection.class);
    when(connection.createStatement(anyString())).thenAnswer(i -> mock(Statement.class));

    connectionInvocationHandler = new ConnectionInvocationHandler(connection, context);
    proxiedConnection = proxyConnection(connectionInvocationHandler);
  }

  @Test
  public void testInvokeCreateStatementEmptyState() {
    proxiedConnection.createStatement(stmt1);
    verify(connection).createStatement(eq(stmt1));
  }

  @Test
  public void testInvokeCreateStatement() {
    State state =
        State.newBuilder()
            .withControllerName("Order")
            .withFramework("spring")
            .withActionName("add")
            .build();
    Context updatedContext = context.put("state", state);

    connectionInvocationHandler = new ConnectionInvocationHandler(connection, updatedContext);
    proxiedConnection = proxyConnection(connectionInvocationHandler);

    proxiedConnection.createStatement(stmt1);
    String appendedState = state.formatAndAppendToSQL(stmt1);
    verify(connection).createStatement(eq(appendedState));
  }

  @Test
  public void testInvokeDelegate() {
    proxiedConnection.beginTransaction();
    proxiedConnection.createStatement(stmt1);
    verify(connection, times(1)).beginTransaction();
  }

  private Connection proxyConnection(ConnectionInvocationHandler connectionInvocationHandler) {
    return (Connection)
        Proxy.newProxyInstance(
            Connection.class.getClassLoader(),
            new Class[] {Connection.class},
            connectionInvocationHandler);
  }
}
