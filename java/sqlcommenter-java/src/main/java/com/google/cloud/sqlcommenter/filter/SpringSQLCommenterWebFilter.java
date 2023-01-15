package com.google.cloud.sqlcommenter.filter;

import org.springframework.web.method.HandlerMethod;
import org.springframework.web.reactive.result.method.annotation.RequestMappingHandlerMapping;
import org.springframework.web.server.ServerWebExchange;
import org.springframework.web.server.WebFilter;
import org.springframework.web.server.WebFilterChain;

import com.google.cloud.sqlcommenter.threadlocalstorage.State;

import reactor.core.publisher.Mono;

public class SpringSQLCommenterWebFilter implements WebFilter {

    private final RequestMappingHandlerMapping handlerMapping;

    public SpringSQLCommenterWebFilter(RequestMappingHandlerMapping handlerMapping) {
        this.handlerMapping = handlerMapping;
    }

    @Override
    public Mono<Void> filter(ServerWebExchange exchange, WebFilterChain chain) {
        HandlerMethod handlerMethod = (HandlerMethod) this.handlerMapping.getHandler(exchange).toFuture().getNow(null);

        State state = State.newBuilder()
                           .withActionName(handlerMethod.getMethod().getName())
                           .withFramework("spring")
                           .withControllerName(handlerMethod.getBeanType().getSimpleName().replace("Controller",""))
                           .build();

        return chain
                .filter(exchange)
                .contextWrite(ctx -> ctx.put("state", state));
    }

}
