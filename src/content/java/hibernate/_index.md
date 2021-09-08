---
title: "Hibernate"
date: 2019-05-31T18:21:05-07:00
draft: false
weight: 1
---

![](/images/hibernate-logo.svg)

- [Introduction](#introduction)
- [Requirements](#requirements)
- [Using the integration](#using-the-integration)
    - [Imports](#imports)
    - [Maven imports](#imports#0)
- [End to end example](#end-to-end-example)
    - [Directory structure](#directory-structure)
    - [Source code](#source-code)
    - [Question.java](#0)
    - [ListQuestions.java](#1)
    - [hibernate.cfg.xml](#2)
    - [Question.hbm.xml](#3)
- [References](#references)

### Introduction
We provide an integration for Hibernate ORM that will inspect and augment your SQL with information about your
setup. It is best used when coupled with other frameworks such as:

{{<card-vendor href="/java/spring" src="/images/spring-logo.png">}}
<!-- {{<card-vendor href="/java/jetty" src="/images/jetty-logo.png">}}
{{<card-vendor href="/java/grpc" src="/images/grpc-logo.png">}}
{{<card-vendor href="/java/tomcat" src="/images/tomcat-logo.png">}} -->

### Requirements

- Java 8+
- Successfully installed [sqlcommenter-java](/java/#install)

### Using the integration

You can include this integration in your Java programs in 2 ways.

#### Persistance XML file

By simply adding to your persistence XML file the property
`"hibernate.session_factory.statement_inspector"`

e.g. to your `hibernate.cfg.xml` file

{{<highlight xml>}}
<property name="hibernate.session_factory.statement_inspector"
          value="com.google.cloud.sqlcommenter.schibernate.SCHibernate" />
{{</highlight>}}

#### In Java source code

When creating your Hibernate session factory, add our StatementInspector like this:

{{<highlight java>}}
import com.google.cloud.sqlcommenter.schhibernate.SCHibernate;

...
        sessionFactoryBuilder.applyStatementInspector(new SCHibernate());
{{</highlight>}}

### Spring and JPA end-to-end example

First thing you need to do is to download the [sqlcommenter-java-guides-spring-jpa](https://github.com/google/sqlcommenter/tree/master/java/sqlcommenter-java#spring)
Java project.

#### Source code

This project uses the following JPA entities:

{{<tabs Post_java Tag_java>}}
{{<highlight java>}}
// In file Post.java
package com.google.cloud.sqlcommenter.spring.jpa.domain;

import javax.persistence.*;
import java.util.ArrayList;
import java.util.List;

@Entity
@Table(name = "post")
public class Post {

    @Id
    @GeneratedValue
    private Long id;

    private String title;

    @ManyToMany
    @JoinTable(
            name = "post_tag",
            joinColumns = @JoinColumn(name = "post_id"),
            inverseJoinColumns = @JoinColumn(name = "tag_id"))
    private List<Tag> tags = new ArrayList<>();

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public List<Tag> getTags() {
        return tags;
    }
}
{{</highlight>}}

{{<highlight java>}}
// In file Tag.java
package com.google.cloud.sqlcommenter.spring.jpa.domain;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.Id;
import javax.persistence.Table;

@Entity
@Table(name = "tag")
public class Tag {

    @Id
    @GeneratedValue
    private Long id;

    private String name;

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }
}
{{</highlight>}}
{{</tabs>}}

The Repository layer looks as follows:

{{<tabs PostRepository_java TagRepository_java>}}
{{<highlight java>}}
// In file PostRepository.java

package com.google.cloud.sqlcommenter.spring.jpa.dao;

import com.google.cloud.sqlcommenter.spring.jpa.domain.Post;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface PostRepository extends CrudRepository<Post, Long> {

    List<Post> findByTitle(String title);
}
{{</highlight>}}

{{<highlight java>}}
// In file TagRepository.java

package com.google.cloud.sqlcommenter.spring.jpa.dao;

import com.google.cloud.sqlcommenter.spring.jpa.domain.Tag;
import org.springframework.data.repository.CrudRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public interface TagRepository extends CrudRepository<Tag, Long> {

    List<Tag> findByNameIn(List<String> names);
}
{{</highlight>}}
{{</tabs>}}

The Service layer looks as follows:

{{<tabs ForumService_java ForumServiceImpl_java>}}
{{<highlight java>}}
// In file ForumService.java

package com.google.cloud.sqlcommenter.spring.jpa.service;

import com.google.cloud.sqlcommenter.spring.jpa.domain.Post;

import java.util.List;

public interface ForumService {

    Post createPost(String title, String... tags);

    List<Post> findPostsByTitle(String title);

    Post findPostById(Long id);
}
{{</highlight>}}

{{<highlight java>}}
// In file ForumServiceImpl.java

package com.google.cloud.sqlcommenter.spring.jpa.service;

import com.google.cloud.sqlcommenter.spring.jpa.dao.PostRepository;
import com.google.cloud.sqlcommenter.spring.jpa.dao.TagRepository;
import com.google.cloud.sqlcommenter.spring.jpa.domain.Post;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.Arrays;
import java.util.List;

@Service
public class ForumServiceImpl implements ForumService {

    @Autowired
    private PostRepository postRepository;

    @Autowired
    private TagRepository tagRepository;

    @Override
    @Transactional
    public Post createPost(String title, String... tags) {
        Post post = new Post();
        post.setTitle(title);
        post.getTags().addAll(tagRepository.findByNameIn(Arrays.asList(tags)));
        return postRepository.save(post);
    }

    @Override
    @Transactional(readOnly = true)
    public List<Post> findPostsByTitle(String title) {
        return postRepository.findByTitle(title);
    }

    @Override
    @Transactional
    public Post findPostById(Long id) {
        return postRepository.findById(id).orElse(null);
    }
}
{{</highlight>}}
{{</tabs>}}

The Spring JPA configuration looks as follows:

````java
// In file JpaTransactionManagerConfiguration.java

package com.google.cloud.sqlcommenter.spring.jpa;

import com.zaxxer.hikari.HikariConfig;
import com.zaxxer.hikari.HikariDataSource;
import com.google.cloud.sqlcommenter.spring.util.SCHibernateWrapper;
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

import javax.persistence.EntityManagerFactory;
import javax.sql.DataSource;
import java.util.Properties;

@Configuration
@PropertySource({"/META-INF/jdbc-hsqldb.properties"})
@ComponentScan(basePackages = "com.google.cloud.sqlcommenter.spring.jpa")
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
    public static PropertySourcesPlaceholderConfigurer properties() {
        return new PropertySourcesPlaceholderConfigurer();
    }

    @Bean
    public DataSource dataSource() {
        return actualDataSource();
    }

    @Bean
    public LocalContainerEntityManagerFactoryBean entityManagerFactory() {
        LocalContainerEntityManagerFactoryBean localContainerEntityManagerFactoryBean = new LocalContainerEntityManagerFactoryBean();
        localContainerEntityManagerFactoryBean.setPersistenceUnitName(getClass().getSimpleName());
        localContainerEntityManagerFactoryBean.setPersistenceProvider(new HibernatePersistenceProvider());
        localContainerEntityManagerFactoryBean.setDataSource(dataSource());
        localContainerEntityManagerFactoryBean.setPackagesToScan(packagesToScan());

        JpaVendorAdapter vendorAdapter = new HibernateJpaVendorAdapter();
        localContainerEntityManagerFactoryBean.setJpaVendorAdapter(vendorAdapter);
        localContainerEntityManagerFactoryBean.setJpaProperties(additionalProperties());
        return localContainerEntityManagerFactoryBean;
    }

    @Bean
    public JpaTransactionManager transactionManager(EntityManagerFactory entityManagerFactory){
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
            "hibernate.session_factory.statement_inspector", SCHibernateWrapper.class.getName()
        );
        return properties;
    }

    protected String[] packagesToScan() {
        return new String[]{
            "com.google.cloud.sqlcommenter.spring.jpa.domain"
        };
    }
}
````

Now, the test looks as follows:

````java
// In file JpaTransactionManagerTest.java

package com.google.cloud.sqlcommenter.spring.jpa;

import com.google.cloud.sqlcommenter.spring.jpa.dao.TagRepository;
import com.google.cloud.sqlcommenter.spring.jpa.domain.Post;
import com.google.cloud.sqlcommenter.spring.jpa.domain.Tag;
import com.google.cloud.sqlcommenter.spring.jpa.service.ForumService;
import com.google.cloud.sqlcommenter.spring.util.SCHibernateWrapper;
import com.google.cloud.sqlcommenter.threadlocalstorage.State;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.annotation.DirtiesContext;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.junit4.SpringJUnit4ClassRunner;

import java.util.List;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration(classes = JpaTransactionManagerConfiguration.class)
@DirtiesContext(classMode = DirtiesContext.ClassMode.AFTER_EACH_TEST_METHOD)
public class JpaTransactionManagerTest {

    @Autowired
    private TagRepository tagRepository;

    @Autowired
    private ForumService forumService;

    @Before
    public void init() {
        Tag hibernate = new Tag();
        hibernate.setName("hibernate");
        tagRepository.save(hibernate);

        Tag jpa = new Tag();
        jpa.setName("jpa");
        tagRepository.save(jpa);
    }

    @Test
    public void test() {
        State.Holder.set(
            State.newBuilder()
                .withControllerName("ForumController")
                .withActionName("CreatePost")
                .withWebFramework("spring")
                .build()
        );

        SCHibernateWrapper.reset();

        Post newPost = forumService.createPost("High-Performance Java Persistence", "hibernate", "jpa");
        assertNotNull(newPost.getId());

        List<String> sqlStatements = SCHibernateWrapper.getAfterSqlStatements();
        assertEquals(5, sqlStatements.size());
        assertEquals(
                5,
                sqlStatements
                    .stream()
                    .filter(
                        sql -> sql.contains(
                            "/*action='CreatePost',controller='ForumController',framework='spring'*/"
                        )
                    )
                    .count()
        );

        SCHibernateWrapper.reset();

        State.Holder.set(
                State.newBuilder()
                        .withControllerName("ForumController")
                        .withActionName("FindPostsByTitle")
                        .withWebFramework("spring")
                        .build()
        );

        List<Post> posts = forumService.findPostsByTitle("High-Performance Java Persistence");
        assertEquals(1, posts.size());

        sqlStatements = SCHibernateWrapper.getAfterSqlStatements();
        assertEquals(1, sqlStatements.size());
        assertEquals(
                1,
                sqlStatements
                    .stream()
                    .filter(
                        sql -> sql.contains(
                            "/*action='FindPostsByTitle',controller='ForumController',framework='spring'*/"
                        )
                    )
                    .count()
        );

        State.Holder.set(
                State.newBuilder()
                        .withControllerName("ForumController")
                        .withActionName("FindPostById")
                        .withWebFramework("spring")
                        .build()
        );

        SCHibernateWrapper.reset();

        Post post = forumService.findPostById(newPost.getId());
        assertEquals("High-Performance Java Persistence", post.getTitle());

        sqlStatements = SCHibernateWrapper.getAfterSqlStatements();
        assertEquals(1, sqlStatements.size());
        assertEquals(
                1,
                sqlStatements
                    .stream()
                    .filter(
                        sql -> sql.contains(
                            "/*action='FindPostById',controller='ForumController',framework='spring'*/"
                        )
                    )
                    .count()
        );
    }
}
````

When running the unit test above, we can see that the SQL statements include the comments as well:

````sql
select tag0_.id as id1_2_, tag0_.name as name2_2_ from tag tag0_ where tag0_.name in (? , ?) /*action='CreatePost',controller='ForumController',framework='spring'*/

call next value for hibernate_sequence /*action='CreatePost',controller='ForumController',framework='spring'*/

insert into post (title, id) values (?, ?) /*action='CreatePost',controller='ForumController',framework='spring'*/

insert into post_tag (post_id, tag_id) values (?, ?) /*action='CreatePost',controller='ForumController',framework='spring'*/

insert into post_tag (post_id, tag_id) values (?, ?) /*action='CreatePost',controller='ForumController',framework='spring'*/

select post0_.id as id1_0_, post0_.title as title2_0_ from post post0_ where post0_.title=? /*action='FindPostsByTitle',controller='ForumController',framework='spring'*/

select post0_.id as id1_0_0_, post0_.title as title2_0_0_ from post post0_ where post0_.id=? /*action='FindPostById',controller='ForumController',framework='spring'*/
````

### Spring and Hibernate end-to-end example

#### Source code

This project uses the following JPA entities:

{{<tabs Post_java Tag_java>}}
{{<highlight java>}}
// In file Post.java
package com.google.cloud.sqlcommenter.spring.hibernate.domain;

import javax.persistence.*;
import java.util.ArrayList;
import java.util.List;

@Entity
@Table(name = "post")
public class Post {

    @Id
    @GeneratedValue
    private Long id;

    private String title;

    @ManyToMany
    @JoinTable(
            name = "post_tag",
            joinColumns = @JoinColumn(name = "post_id"),
            inverseJoinColumns = @JoinColumn(name = "tag_id"))
    private List<Tag> tags = new ArrayList<>();

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getTitle() {
        return title;
    }

    public void setTitle(String title) {
        this.title = title;
    }

    public List<Tag> getTags() {
        return tags;
    }
}
{{</highlight>}}

{{<highlight java>}}
// In file Tag.java
package com.google.cloud.sqlcommenter.spring.hibernate.domain;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.Id;
import javax.persistence.Table;

@Entity
@Table(name = "tag")
public class Tag {

    @Id
    @GeneratedValue
    private Long id;

    private String name;

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name;
    }
}
{{</highlight>}}
{{</tabs>}}

The DAO (Data Access Object) layer looks as follows:

{{<tabs GenericDAO_java GenericDAOImpl_java PostDAO_java PostDAOImpl_java TagDAO_java TagDAOImpl_java>}}
{{<highlight java>}}
// In file GenericDAO.java

package com.google.cloud.sqlcommenter.spring.hibernate.dao;

import java.io.Serializable;

public interface GenericDAO<T, ID extends Serializable> {

    T findById(ID id);

    T save(T entity);
}
{{</highlight>}}

{{<highlight java>}}
// In file GenericDAOImpl.java

package com.google.cloud.sqlcommenter.spring.hibernate.dao;

import org.hibernate.Session;
import org.hibernate.SessionFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Repository;
import org.springframework.transaction.annotation.Transactional;

import java.io.Serializable;

@Repository
@Transactional
public abstract class GenericDAOImpl<T, ID extends Serializable> implements GenericDAO<T, ID> {

    @Autowired
    private SessionFactory sessionFactory;

    private final Class<T> entityClass;

    protected SessionFactory getSessionFactory() {
        return sessionFactory;
    }

    protected Session getSession() {
        return sessionFactory.getCurrentSession();
    }

    protected GenericDAOImpl(Class<T> entityClass) {
        this.entityClass = entityClass;
    }

    public Class<T> getEntityClass() {
        return entityClass;
    }

    @Override
    public T findById(ID id) {
        return getSession().get(entityClass, id);
    }

    @Override
    public T save(T entity) {
        getSession().persist(entity);
        return entity;
    }
}
{{</highlight>}}

{{<highlight java>}}
// In file PostDAO.java

package com.google.cloud.sqlcommenter.spring.hibernate.dao;

import com.google.cloud.sqlcommenter.spring.hibernate.domain.Post;

import java.util.List;

public interface PostDAO extends GenericDAO<Post, Long> {

    List<Post> findByTitle(String title);
}
{{</highlight>}}

{{<highlight java>}}
// In file PostDAOImpl.java

package com.google.cloud.sqlcommenter.spring.hibernate.dao;

import com.google.cloud.sqlcommenter.spring.hibernate.domain.Post;
import org.springframework.stereotype.Repository;

import java.util.List;

@Repository
public class PostDAOImpl extends GenericDAOImpl<Post, Long> implements PostDAO {

    protected PostDAOImpl() {
        super(Post.class);
    }

    @Override
    public List<Post> findByTitle(String title) {
        return getSession().createQuery(
                "select p from Post p where p.title = :title", Post.class)
                .setParameter("title", title)
                .getResultList();
    }
}
{{</highlight>}}

{{<highlight java>}}
// In file TagDAO.java

package com.google.cloud.sqlcommenter.spring.hibernate.dao;

import com.google.cloud.sqlcommenter.spring.hibernate.domain.Tag;

import java.util.List;

public interface TagDAO extends GenericDAO<Tag, Long> {

    List<Tag> findByName(String... tags);
}
{{</highlight>}}

{{<highlight java>}}
// In file TagDAOImpl.java

package com.google.cloud.sqlcommenter.spring.hibernate.dao;

import com.google.cloud.sqlcommenter.spring.hibernate.domain.Tag;
import org.springframework.stereotype.Repository;

import java.util.Arrays;
import java.util.List;

@Repository
public class TagDAOImpl extends GenericDAOImpl<Tag, Long> implements TagDAO {

    protected TagDAOImpl() {
        super(Tag.class);
    }

    @Override
    public List<Tag> findByName(String... tags) {
        if (tags.length == 0) {
            throw new IllegalArgumentException("There's no tag name to search for!");
        }
        return getSession()
            .createQuery(
                "select t " +
                "from Tag t " +
                "where t.name in :tags")
            .setParameterList("tags", Arrays.asList(tags))
            .list();
    }
}
{{</highlight>}}
{{</tabs>}}

The Service layer looks as follows:

{{<tabs ForumService_java ForumServiceImpl_java>}}
{{<highlight java>}}
// In file ForumService.java

package com.google.cloud.sqlcommenter.spring.hibernate.service;

import com.google.cloud.sqlcommenter.spring.hibernate.domain.Post;

import java.util.List;

public interface ForumService {

    Post createPost(String title, String... tags);

    List<Post> findPostsByTitle(String title);

    Post findPostById(Long id);
}
{{</highlight>}}

{{<highlight java>}}
// In file ForumServiceImpl.java

package com.google.cloud.sqlcommenter.spring.hibernate.service;

import com.google.cloud.sqlcommenter.spring.hibernate.dao.PostDAO;
import com.google.cloud.sqlcommenter.spring.hibernate.dao.TagDAO;
import com.google.cloud.sqlcommenter.spring.hibernate.domain.Post;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

@Service
public class ForumServiceImpl implements ForumService {

    @Autowired
    private PostDAO postDAO;

    @Autowired
    private TagDAO tagDAO;

    @Override
    @Transactional
    public Post createPost(String title, String... tags) {
        Post post = new Post();
        post.setTitle(title);
        post.getTags().addAll(tagDAO.findByName(tags));
        return postDAO.save(post);
    }

    @Override
    @Transactional(readOnly = true)
    public List<Post> findPostsByTitle(String title) {
        return postDAO.findByTitle(title);
    }

    @Override
    @Transactional
    public Post findPostById(Long id) {
        return postDAO.findById(id);
    }
}
{{</highlight>}}
{{</tabs>}}

The Spring Hibernate configuration looks as follows:

````java
// In file HibernateTransactionManagerConfiguration.java

package com.google.cloud.sqlcommenter.spring.hibernate;

import com.zaxxer.hikari.HikariConfig;
import com.zaxxer.hikari.HikariDataSource;
import com.google.cloud.sqlcommenter.spring.util.SCHibernateWrapper;
import org.hibernate.SessionFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.*;
import org.springframework.context.support.PropertySourcesPlaceholderConfigurer;
import org.springframework.orm.hibernate5.HibernateTransactionManager;
import org.springframework.orm.hibernate5.LocalSessionFactoryBean;
import org.springframework.transaction.annotation.EnableTransactionManagement;
import org.springframework.transaction.support.TransactionTemplate;

import javax.sql.DataSource;
import java.util.Properties;

@Configuration
@PropertySource({"/META-INF/jdbc-hsqldb.properties"})
@ComponentScan(basePackages = "com.google.cloud.sqlcommenter.spring.hibernate")
@EnableTransactionManagement
@EnableAspectJAutoProxy
public class HibernateTransactionManagerConfiguration {

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
    public static PropertySourcesPlaceholderConfigurer properties() {
        return new PropertySourcesPlaceholderConfigurer();
    }

    @Bean
    public DataSource dataSource() {
        return actualDataSource();
    }

    @Bean
    public LocalSessionFactoryBean sessionFactory() {
        LocalSessionFactoryBean localSessionFactoryBean = new LocalSessionFactoryBean();
        localSessionFactoryBean.setDataSource(dataSource());
        localSessionFactoryBean.setPackagesToScan(packagesToScan());
        localSessionFactoryBean.setHibernateProperties(additionalProperties());
        return localSessionFactoryBean;
    }

    @Bean
    public HibernateTransactionManager transactionManager(SessionFactory sessionFactory) {
        HibernateTransactionManager transactionManager = new HibernateTransactionManager();
        transactionManager.setSessionFactory(sessionFactory);
        return transactionManager;
    }

    @Bean
    public TransactionTemplate transactionTemplate(SessionFactory sessionFactory) {
        return new TransactionTemplate(transactionManager(sessionFactory));
    }

    protected Properties additionalProperties() {
        Properties properties = new Properties();
        properties.setProperty("hibernate.dialect", hibernateDialect);
        properties.setProperty("hibernate.hbm2ddl.auto", "create-drop");
        properties.put(
            "hibernate.session_factory.statement_inspector", SCHibernateWrapper.class.getName()
        );
        return properties;
    }

    protected String[] packagesToScan() {
        return new String[]{
                "com.google.cloud.sqlcommenter.spring.hibernate.domain"
        };
    }
}
````

Now, the test looks as follows:

````java
// In file HibernateTransactionManagerTest.java

package com.google.cloud.sqlcommenter.spring.hibernate;

import com.google.cloud.sqlcommenter.spring.hibernate.domain.Post;
import com.google.cloud.sqlcommenter.spring.hibernate.domain.Tag;
import com.google.cloud.sqlcommenter.spring.hibernate.service.ForumService;
import com.google.cloud.sqlcommenter.spring.util.SCHibernateWrapper;
import com.google.cloud.sqlcommenter.threadlocalstorage.State;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.annotation.DirtiesContext;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.junit4.SpringJUnit4ClassRunner;
import org.springframework.transaction.TransactionException;
import org.springframework.transaction.support.TransactionCallback;
import org.springframework.transaction.support.TransactionTemplate;

import javax.persistence.EntityManager;
import javax.persistence.PersistenceContext;
import java.util.List;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.fail;

@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration(classes = HibernateTransactionManagerConfiguration.class)
@DirtiesContext(classMode = DirtiesContext.ClassMode.AFTER_EACH_TEST_METHOD)
public class HibernateTransactionManagerTest {

    @Autowired
    private TransactionTemplate transactionTemplate;

    @PersistenceContext
    private EntityManager entityManager;

    @Autowired
    private ForumService forumService;

    @Before
    public void init() {
        try {
            transactionTemplate.execute((TransactionCallback<Void>) transactionStatus -> {
                Tag hibernate = new Tag();
                hibernate.setName("hibernate");
                entityManager.persist(hibernate);

                Tag jpa = new Tag();
                jpa.setName("jpa");
                entityManager.persist(jpa);

                return null;
            });
        } catch (TransactionException e) {
            fail(e.getMessage());
        }
    }

    @Test
    public void test() {
        State.Holder.set(
            State.newBuilder()
                .withControllerName("ForumController")
                .withActionName("CreatePost")
                .withWebFramework("spring")
                .build()
        );

        SCHibernateWrapper.reset();

        Post newPost = forumService.createPost("High-Performance Java Persistence", "hibernate", "jpa");
        assertNotNull(newPost.getId());

        List<String> sqlStatements = SCHibernateWrapper.getAfterSqlStatements();
        assertEquals(5, sqlStatements.size());
        assertEquals(
                5,
                sqlStatements
                .stream()
                .filter(
                    sql -> sql.contains(
                        "/*action='CreatePost',controller='ForumController',framework='spring'*/"
                    )
                )
                .count()
        );

        SCHibernateWrapper.reset();

        State.Holder.set(
            State.newBuilder()
                .withControllerName("ForumController")
                .withActionName("FindPostsByTitle")
                .withWebFramework("spring")
                .build()
        );

        List<Post> posts = forumService.findPostsByTitle("High-Performance Java Persistence");
        assertEquals(1, posts.size());

        sqlStatements = SCHibernateWrapper.getAfterSqlStatements();
        assertEquals(1, sqlStatements.size());
        assertEquals(
                1,
                sqlStatements
                    .stream()
                    .filter(
                        sql -> sql.contains(
                            "/*action='FindPostsByTitle',controller='ForumController',framework='spring'*/"
                        )
                    )
                    .count()
        );

        State.Holder.set(
            State.newBuilder()
                    .withControllerName("ForumController")
                    .withActionName("FindPostById")
                    .withWebFramework("spring")
                    .build()
        );

        SCHibernateWrapper.reset();

        Post post = forumService.findPostById(newPost.getId());
        assertEquals("High-Performance Java Persistence", post.getTitle());

        sqlStatements = SCHibernateWrapper.getAfterSqlStatements();
        assertEquals(1, sqlStatements.size());
        assertEquals(
                1,
                sqlStatements
                    .stream()
                    .filter(
                        sql -> sql.contains(
                            "/*action='FindPostById',controller='ForumController',framework='spring'*/"
                        )
                    )
                    .count()
        );
    }
}
````

When running the unit test above, we can see that the SQL statements include the comments as well:

````sql
select tag0_.id as id1_2_, tag0_.name as name2_2_ from tag tag0_ where tag0_.name in (? , ?) /*action='CreatePost',controller='ForumController',framework='spring'*/

call next value for hibernate_sequence /*action='CreatePost',controller='ForumController',framework='spring'*/

insert into post (title, id) values (?, ?) /*action='CreatePost',controller='ForumController',framework='spring'*/

insert into post_tag (post_id, tag_id) values (?, ?) /*action='CreatePost',controller='ForumController',framework='spring'*/

insert into post_tag (post_id, tag_id) values (?, ?) /*action='CreatePost',controller='ForumController',framework='spring'*/

select post0_.id as id1_0_, post0_.title as title2_0_ from post post0_ where post0_.title=? /*action='FindPostsByTitle',controller='ForumController',framework='spring'*/

select post0_.id as id1_0_0_, post0_.title as title2_0_0_ from post post0_ where post0_.id=? /*action='FindPostById',controller='ForumController',framework='spring'*/
````

### References

Resource|URL
---|---
Hibernate ORM project|https://hibernate.org/orm/
sqlcommenter-java on Github|https://github.com/google/sqlcommenter/tree/master/java/sqlcommenter-java
