package com.google.cloud.sqlcommenter.r2dbc;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.times;
import static org.mockito.Mockito.verify;

import io.r2dbc.spi.Connection;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;
import org.reactivestreams.Subscription;
import reactor.core.CoreSubscriber;

@RunWith(JUnit4.class)
public class ConnectionDecoratorTest {

  private CoreSubscriber<Object> delegate;
  private ConnectionDecorator connectionDecorator;

  @Before
  public void setUp() throws Exception {
    delegate = mock(CoreSubscriber.class);
    connectionDecorator = new ConnectionDecorator(delegate);
  }

  @Test
  public void testOnSubscribe() {
    Subscription mocked = mock(Subscription.class);
    connectionDecorator.onSubscribe(mocked);
    verify(delegate, times(1)).onSubscribe(mocked);
  }

  @Test
  public void testOnNext() {
    Connection connection = mock(Connection.class);
    connectionDecorator.onNext(connection);
    verify(delegate, times(1)).currentContext();
    verify(delegate, times(1)).onNext(any());
  }

  @Test(expected = AssertionError.class)
  public void testOnNextNonConnection() {
    Object object = new Object();
    connectionDecorator.onNext(object);
  }

  @Test
  public void testOnError() {
    Throwable mocked = mock(Throwable.class);
    connectionDecorator.onError(mocked);
    verify(delegate, times(1)).onError(mocked);
  }

  @Test
  public void testOnComplete() {
    connectionDecorator.onComplete();
    verify(delegate, times(1)).onComplete();
  }
}
