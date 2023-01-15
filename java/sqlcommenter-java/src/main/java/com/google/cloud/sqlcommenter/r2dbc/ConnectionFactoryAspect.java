package com.google.cloud.sqlcommenter.r2dbc;



import org.aspectj.lang.ProceedingJoinPoint;
import org.aspectj.lang.annotation.Around;
import org.aspectj.lang.annotation.Aspect;
import org.reactivestreams.Publisher;

import reactor.core.publisher.Mono;
import reactor.core.publisher.Operators;

@Aspect
public class ConnectionFactoryAspect {

    @Around("execution(* io.r2dbc.spi.ConnectionFactory.create(..)) ")
    public Object beforeSampleCreation(ProceedingJoinPoint joinPoint) throws Throwable {
        Object object = joinPoint.proceed();

        @SuppressWarnings("unchecked") Publisher<Object> publisher = (Publisher<Object>) object;

        return Mono.from(publisher)
                   .transform(Operators.liftPublisher((publisher1, coreSubscriber) -> new ConnectionDecorator(coreSubscriber)));
    }

}
