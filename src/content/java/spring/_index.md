---
title: "Spring"
date: 2019-05-31T18:20:05-07:00
draft: false
weight: 2
---

![](/images/spring-logo.png)

- [Introduction](#introduction)
- [Requirements](#requirements)
- [Dependency management](#dependency-management)
    - [Manually](#manually)
    - [Package management](#package-management)
- [Expected fields](#expected-fields)
- [Using it](#using-it)
    - [Spring 5](#spring-5)
    - [Before Spring 5](#before-spring-5)
- [XML based configuration](#xml-based-configuration)
- [Hibernate](#hibernate)
- [References](#references)

### Introduction
We provide an integration for the Spring framework. The integration is an [interceptor](https://docs.spring.io/spring/docs/5.0.4.BUILD-SNAPSHOT/javadoc-api/org/aopalliance/intercept/Interceptor.html) that will record properties about your MVC application per HTTP request, and then
later those properties will be picked up by the actual ORMs and augment your SQL statements. It is best used with the following ORM integrations:

{{<card-vendor href="/java/hibernate" src="/images/hibernate-logo.svg">}}

### Requirements

- Java 8+
- Successfully installed [sqlcommenter-java](/java/#install)

### Dependency management

We can add the integration to our applications in the following ways:

#### Manually

Please read [installing sqlcommenter-java from source](/java/#install-from-source)

#### Package management

Please include this in your dependency management system as follows

{{<tabs Maven Gradle>}}
{{<highlight xml>}}
    <dependency>
        <groupId>com.google.cloud</groupId>
        <artifactId>sqlcommenter-java</artifactId>
        <version>0.0.1</version>
    </dependency>
{{</highlight>}}

{{<highlight gradle>}}
// https://mvnrepository.com/artifact/com.google.cloud/sqlcommenter-java
compile group: 'com.google.cloud', name: 'sqlcommenter-java', version: '0.0.1'
{{</highlight>}}

{{</tabs>}}

### Expected fields
When coupled say with [sqlcommenter for Hibernate](/java/hibernate), the following fields will be added to your SQL statement as comments

Field|Description
---|---
action|The name of the command that execute the logical behavior e.g. `'/fees'`
controller|The name of your controller e.g. `'fees_controller'`
web\_framework|The name of the framework, it will always be `'spring'`

### Using it
There are 2 different flavors of Spring -- Spring 5 and later vs before Spring 5. Please read along to see
how to enable it for the different versions:

#### Spring 5
If using Spring 5, please import the `SpringSQLCommenterInterceptor` class by:

{{<highlight java>}}
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
{{</highlight>}}

#### Before Spring 5

If using a version before Spring 5, your `WebConfig` class needs to extend the [WebMVCConfigureAdapter](https://docs.spring.io/spring/docs/current/javadoc-api/org/springframework/web/servlet/config/annotation/WebMvcConfigurerAdapter.html) class instead like this:

{{<highlight java>}}
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
{{</highlight>}}

### XML based configuration

You can add the interceptor as a bean in your XML configuration

{{<tabs For_Every_Method Method_Specific>}}
{{<highlight xml>}}
<mvc:interceptors>
    <bean class="com.google.cloud.sqlcommenter.interceptors.SpringSQLCommenterInterceptor"></bean>
</mvc:interceptors>
{{</highlight>}}

{{<highlight xml>}}
<mvc:interceptors>
    <mvc:interceptor>
        <mvc:mapping path="/flights"></mvc:mapping>
        <bean class="com.google.cloud.sqlcommenter.interceptors.SpringSQLCommenterInterceptor"></bean>
    </mvc:interceptor>
</mvc:interceptors>
{{</highlight>}}
{{</tabs>}}

### Hibernate

If Spring is using Hibernate, in addtion to the step [XML based configuration](#xml-based-configuration),
since you might not be using a `persistence.xml` file, we can setup in Java code the
`hibernate.session_factory.statement_inspector` configuration property in your `additionalProperties` method as per

{{<highlight java>}}
import com.google.cloud.sqlcommenter.schhibernate.SCHibernate;

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
{{</highlight>}}

### References

Resource|URL
---|---
Spring framework homepage|https://spring.io/
sqlcommenter-java on Github|https://github.com/google/sqlcommenter/tree/master/java/sqlcommenter-java
Spring Interceptor|https://docs.spring.io/spring/docs/5.0.4.BUILD-SNAPSHOT/javadoc-api/org/aopalliance/intercept/Interceptor.html
Hibernate SQLCommenter integration|[/java/hibernate](/java/hibernate)
