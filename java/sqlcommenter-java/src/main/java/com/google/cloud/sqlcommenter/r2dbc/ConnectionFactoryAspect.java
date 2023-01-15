package com.google.cloud.sqlcommenter.r2dbc;



import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.reactivestreams.Publisher;

import io.r2dbc.spi.Connection;
import reactor.core.publisher.Mono;
import reactor.core.publisher.Operators;

@Aspect
public class ConnectionFactoryAspect {

    @Around("execution(* io.r2dbc.spi.ConnectionFactory.create(..)) ")
    public Object beforeSampleCreation(ProceedingJoinPoint joinPoint) throws Throwable {
        Publisher<? extends Connection> publisher = (Publisher<? extends Connection>) joinPoint.proceed();


        return Mono.from(publisher)
                   .cast(Object.class)
                   .transform(Operators.liftPublisher((publisher1, coreSubscriber) -> new ConnectionDecorator(coreSubscriber)));
    }

}
