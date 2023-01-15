package com.google.cloud.sqlcommenter.r2dbc;

import com.google.cloud.sqlcommenter.threadlocalstorage.State;
import io.r2dbc.spi.Connection;
import java.lang.reflect.InvocationHandler;
import java.lang.reflect.Method;
import reactor.util.context.ContextView;

public class ConnectionInvocationHandler implements InvocationHandler {

  private final Connection connection;
  private final ContextView contextView;

  public ConnectionInvocationHandler(Connection connection, ContextView contextView) {
    this.connection = connection;
    this.contextView = contextView;
  }

  @Override
  public Object invoke(Object proxy, Method method, Object[] args) throws Throwable {
    String methodName = method.getName();

    if ("createStatement".equals(methodName)) {
      String query = (String) args[0];
      if (contextView != null) {
        State state = contextView.get("state");
        query = state.formatAndAppendToSQL(query);
      }

      return method.invoke(connection, query);
    } else {
      return method.invoke(connection, args);
    }
  }
}
