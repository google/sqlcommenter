## sqlcommenter-java

- [Introduction](#introduction)
- [Integrations](#integrations)
    - [ORMs](#orms)
    - [Frameworks](#frameworks)
- [Using them](#using-them)
    - [Spring](#spring)
        - [In XML Configuration](#in-xml-configuration)
        - [In your Java source code](#in-your-java-source-code)
    - [Hibernate](#hibernate)
    - [Spring Hibernate](#spring-hibernate)

### Introduction
Provides integrations to correlate user source code from various
web frameworks with SQL comments from various ORMs.

When results are examined in SQL database logs, they'll look like this:

```shell
SELECT * from USERS /*action='run+this+%26+that',
controller='foo%3BDROP+TABLE+BAR',framework='spring,
traceparent='00-9a4589fe88dd0fc911ff2233ffee7899-11fa8b00dd11eeff-01',
tracestate='rojo%253D00f067aa0ba902b7%2Ccongo%253Dt61rcWkgMzE''*/
```

### Integrations

#### ORMs

- [X] Hibernate
- [ ] JDBC

#### Frameworks

- [X] Spring
- [ ] Jetty
- [ ] Netty
- [ ] Apache Tomcat
- [X] gRPC

### Using it

#### Spring

##### Java-based configuration

If your're using Spring 5, then you can add the `SpringSQLCommenterInterceptor` as follows:

```java
import com.google.cloud.sqlcommenter.interceptors.SpringSQLCommenterInterceptor;

@EnableWebMvc
@Configuration
public class WebConfig implements WebMvcConfigurer {

    @Bean
    public SpringSQLCommenterInterceptor sqlInterceptor() {
         return new SpringSQLCommenterInterceptor();
    }
 
    @Override
    public void addInterceptors(InterceptorRegistry registry) {
        registry.addInterceptor(sqlInterceptor());
    }
}
```

If you're using an older version of Spring, then your `WebConfig` class needs to extend the `WebMvcConfigureAdapter`
class instead:

```java
import com.google.cloud.sqlcommenter.interceptors.SpringSQLCommenterInterceptor;

@EnableWebMvc
@Configuration
public class WebConfig extends WebMvcConfigureAdapter {

    @Bean
    public SpringSQLCommenterInterceptor sqlInterceptor() {
         return new SpringSQLCommenterInterceptor();
    }
 
    @Override
    public void addInterceptors(InterceptorRegistry registry) {
        registry.addInterceptor(sqlInterceptor());
    }
}
```

##### XML-based configuration

You can also add the interceptor as a bean in your XML configuration:

```xml
<mvc:interceptors>
    <bean class="com.google.cloud.sqlcommenter.interceptors.SpringSQLCommenterInterceptor"></bean>
</mvc:interceptors>
```

or alternatively just for a specific method
```xml
<mvc:interceptors>
    <mvc:interceptor>
        <mvc:mapping path="/flights"></mvc:mapping>
        <bean class="com.google.cloud.sqlcommenter.interceptors.SpringSQLCommenterInterceptor"></bean>
    </mvc:interceptor>
</mvc:interceptors>
```

#### Hibernate

##### Hibernate via JPA

If you're using Hibernate via JPA, then you can simply set the `hibernate.session_factory.statement_inspector` configuration property in the `persistence.xml` configuration file:

```xml
<property name="hibernate.session_factory.statement_inspector" value="com.google.cloud.sqlcommenter.schibernate.SCHibernate" />
```

##### Hibernate via Spring Application Properties
If you're using Hibernate via Spring application properties(i.e. you have `application.yml` or `application.properties` in your resources folder)
- `application.properties`

```properties
spring.jpa.properties.hibernate.session_factory.statement_inspector=com.google.cloud.sqlcommenter.schibernate.SCHibernate
```

- `application.yml`

```yml 
spring:
  jpa:
    properties:
      hibernate:
        session_factory:
          statement_inspector: com.google.cloud.sqlcommenter.schibernate.SCHibernate
```

##### Hibernate via Spring

If you're using Hibernate via Spring, then you might not use a `persistence.xml` configuration file, in which case,
you can set up the `hibernate.session_factory.statement_inspector` configuration property as follows:

```java
@Configuration
@EnableTransactionManagement
public class JPAConfig {
 
    @Bean
    public LocalContainerEntityManagerFactoryBean entityManagerFactory() {
        LocalContainerEntityManagerFactoryBean em 
        = new LocalContainerEntityManagerFactoryBean();
        em.setDataSource(dataSource());
        em.setPackagesToScan(new String[] { "you.application.domain.model" });

        em.setJpaVendorAdapter(new HibernateJpaVendorAdapter());
        em.setJpaProperties(additionalProperties());

        return em;
    }
    
    private Properties additionalProperties() {
        Properties properties = new Properties();
        properties.setProperty("hibernate.session_factory.statement_inspector", SCHibernate.class.getName());

        return properties;
    }
}
```

#### Spring Hibernate

1. Please follow the instructions to add the [Spring interceptor](#spring)
2. Please follow the instructions to add the [Hibernate StatementInspector](#hibernate)

#### Spring gRPC Configuration 
Add following configuration into your code
```java
@Configuration
public class gRPCConfig {

    @Bean
    public GrpcSQLCommenterInterceptor grpcSQLCommenterInterceptor() {
        return new GrpcSQLCommenterInterceptor();
    }

}
```