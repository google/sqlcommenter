// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package com.google.cloud.sqlcommenter.spring.backend;

import com.zaxxer.hikari.HikariConfig;
import com.zaxxer.hikari.HikariDataSource;
import com.google.cloud.sqlcommenter.util.SCHibernateWrapper;
import java.util.Properties;
import javax.persistence.EntityManagerFactory;
import javax.sql.DataSource;
import org.hibernate.jpa.HibernatePersistenceProvider;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.*;
import org.springframework.context.support.PropertySourcesPlaceholderConfigurer;
import org.springframework.data.jpa.repository.config.EnableJpaRepositories;
import org.springframework.orm.jpa.JpaTransactionManager;
import org.springframework.orm.jpa.JpaVendorAdapter;
import org.springframework.orm.jpa.LocalContainerEntityManagerFactoryBean;
import org.springframework.orm.jpa.vendor.HibernateJpaVendorAdapter;
import org.springframework.transaction.annotation.EnableTransactionManagement;
import org.springframework.transaction.support.TransactionTemplate;

@Configuration
@PropertySource({"/META-INF/jdbc-hsqldb.properties"})
@ComponentScan(basePackages = "com.google.cloud.sqlcommenter.spring.backend")
@EnableTransactionManagement
@EnableAspectJAutoProxy
@EnableJpaRepositories
public class JpaTransactionManagerConfiguration {

  @Value("${jdbc.dataSourceClassName}")
  private String dataSourceClassName;

  @Value("${jdbc.url}")
  private String jdbcUrl;

  @Value("${jdbc.username}")
  private String jdbcUser;

  @Value("${jdbc.password}")
  private String jdbcPassword;

  @Value("${hibernate.dialect}")
  private String hibernateDialect;

  @Bean
  public static PropertySourcesPlaceholderConfigurer properties() {
    return new PropertySourcesPlaceholderConfigurer();
  }

  @Bean(destroyMethod = "close")
  public DataSource actualDataSource() {
    Properties driverProperties = new Properties();
    driverProperties.setProperty("url", jdbcUrl);
    driverProperties.setProperty("user", jdbcUser);
    driverProperties.setProperty("password", jdbcPassword);

    Properties properties = new Properties();
    properties.put("dataSourceClassName", dataSourceClassName);
    properties.put("dataSourceProperties", driverProperties);
    properties.setProperty("maximumPoolSize", String.valueOf(3));
    return new HikariDataSource(new HikariConfig(properties));
  }

  @Bean
  public DataSource dataSource() {
    return actualDataSource();
  }

  @Bean
  public LocalContainerEntityManagerFactoryBean entityManagerFactory() {
    LocalContainerEntityManagerFactoryBean localContainerEntityManagerFactoryBean =
        new LocalContainerEntityManagerFactoryBean();
    localContainerEntityManagerFactoryBean.setPersistenceUnitName(getClass().getSimpleName());
    localContainerEntityManagerFactoryBean.setPersistenceProvider(
        new HibernatePersistenceProvider());
    localContainerEntityManagerFactoryBean.setDataSource(dataSource());
    localContainerEntityManagerFactoryBean.setPackagesToScan(packagesToScan());

    JpaVendorAdapter vendorAdapter = new HibernateJpaVendorAdapter();
    localContainerEntityManagerFactoryBean.setJpaVendorAdapter(vendorAdapter);
    localContainerEntityManagerFactoryBean.setJpaProperties(additionalProperties());
    return localContainerEntityManagerFactoryBean;
  }

  @Bean
  public JpaTransactionManager transactionManager(EntityManagerFactory entityManagerFactory) {
    JpaTransactionManager transactionManager = new JpaTransactionManager();
    transactionManager.setEntityManagerFactory(entityManagerFactory);
    return transactionManager;
  }

  @Bean
  public TransactionTemplate transactionTemplate(EntityManagerFactory entityManagerFactory) {
    return new TransactionTemplate(transactionManager(entityManagerFactory));
  }

  protected Properties additionalProperties() {
    Properties properties = new Properties();
    properties.setProperty("hibernate.dialect", hibernateDialect);
    properties.setProperty("hibernate.hbm2ddl.auto", "create-drop");
    properties.put(
        "hibernate.session_factory.statement_inspector", SCHibernateWrapper.class.getName());
    return properties;
  }

  protected String[] packagesToScan() {
    return new String[] {"com.google.cloud.sqlcommenter.spring.backend.domain"};
  }
}
