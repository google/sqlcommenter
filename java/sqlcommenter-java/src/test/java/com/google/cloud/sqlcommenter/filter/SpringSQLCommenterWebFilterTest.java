package com.google.cloud.sqlcommenter.filter;

import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.JUnit4;
import org.springframework.mock.http.server.reactive.MockServerHttpRequest;
import org.springframework.mock.web.server.MockServerWebExchange;
import org.springframework.web.reactive.result.method.RequestMappingInfo;
import org.springframework.web.reactive.result.method.annotation.RequestMappingHandlerMapping;
import org.springframework.web.server.WebFilterChain;

import reactor.core.publisher.Mono;
import reactor.test.StepVerifier;


@RunWith(JUnit4.class)
public class SpringSQLCommenterWebFilterTest {

    @Test
    public void testPreHandlePlacesStateInContextView() throws NoSuchMethodException {
        RequestMappingHandlerMapping requestMappingHandlerMapping = new RequestMappingHandlerMapping();
        RequestMappingInfo info = RequestMappingInfo.paths("/test").build();
        requestMappingHandlerMapping.registerMapping(info, this, SpringSQLCommenterWebFilterTest.class.getMethod("testPreHandlePlacesStateInContextView"));

        SpringSQLCommenterWebFilter springSQLCommenterWebFilter = new SpringSQLCommenterWebFilter(requestMappingHandlerMapping);

        WebFilterChain filterChain = filterExchange -> Mono.empty();

        MockServerWebExchange exchange = MockServerWebExchange.from(
                MockServerHttpRequest
                        .get("/test"));

        StepVerifier.create(springSQLCommenterWebFilter.filter(exchange, filterChain))
                .expectAccessibleContext()
                .hasKey("state")
                .then()
                .verifyComplete();
    }

}