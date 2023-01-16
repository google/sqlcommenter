package com.google.cloud.sqlcommenter.r2dbc;

import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.context.annotation.EnableAspectJAutoProxy;

import io.r2dbc.h2.H2ConnectionFactory;
import io.r2dbc.spi.ConnectionFactory;

@Configuration
@EnableAspectJAutoProxy
public class R2DBCConfiguration {

    @Bean
    public ConnectionFactory connectionFactory() {
        return new H2ConnectionFactory(
                io.r2dbc.h2.H2ConnectionConfiguration.builder()
                                                     .url("mem:testdb;MODE=PostgreSQL;DB_CLOSE_DELAY=-1;")
                                                     .build()
        );
    }

    @Bean
    public ConnectionFactoryAspect connectionFactoryAspect() {
        return new ConnectionFactoryAspect();
    }

}
