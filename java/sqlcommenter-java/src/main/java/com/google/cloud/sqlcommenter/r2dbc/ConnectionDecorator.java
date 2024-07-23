package com.google.cloud.sqlcommenter.r2dbc;

import io.r2dbc.spi.Connection;
import java.lang.reflect.Proxy;
import org.reactivestreams.Subscription;
import reactor.core.CoreSubscriber;

public class ConnectionDecorator implements CoreSubscriber<Object> {

  private final CoreSubscriber<Object> delegate;

  public ConnectionDecorator(CoreSubscriber<Object> delegate) {
    this.delegate = delegate;
  }

  @Override
  public void onSubscribe(Subscription s) {
    this.delegate.onSubscribe(s);
  }

  @Override
  public void onNext(Object o) {
    assert o instanceof Connection;
    Connection connection = (Connection) o;

    Object proxied =
        Proxy.newProxyInstance(
            Connection.class.getClassLoader(),
            new Class[] {Connection.class},
            new ConnectionInvocationHandler(connection, delegate.currentContext()));
    this.delegate.onNext(proxied);
  }

  @Override
  public void onError(Throwable t) {
    this.delegate.onError(t);
  }

  @Override
  public void onComplete() {
    this.delegate.onComplete();
  }
}
