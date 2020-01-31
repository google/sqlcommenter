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

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertNotNull;

import com.google.cloud.sqlcommenter.spring.backend.domain.Post;
import com.google.cloud.sqlcommenter.spring.backend.domain.Tag;
import com.google.cloud.sqlcommenter.spring.backend.service.ForumService;
import com.google.cloud.sqlcommenter.threadlocalstorage.State;
import com.google.cloud.sqlcommenter.util.SCHibernateWrapper;
import java.util.List;
import javax.persistence.EntityManager;
import javax.persistence.PersistenceContext;
import org.junit.Before;
import org.junit.Test;
import org.junit.runner.RunWith;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.test.annotation.DirtiesContext;
import org.springframework.test.context.ContextConfiguration;
import org.springframework.test.context.junit4.SpringJUnit4ClassRunner;
import org.springframework.transaction.TransactionException;
import org.springframework.transaction.support.TransactionCallback;
import org.springframework.transaction.support.TransactionTemplate;

@RunWith(SpringJUnit4ClassRunner.class)
@ContextConfiguration(classes = JpaTransactionManagerConfiguration.class)
@DirtiesContext(classMode = DirtiesContext.ClassMode.AFTER_EACH_TEST_METHOD)
public class JpaTransactionManagerTest {

  protected final Logger LOGGER = LoggerFactory.getLogger(getClass());

  @Autowired private TransactionTemplate transactionTemplate;

  @PersistenceContext private EntityManager entityManager;

  @Autowired private ForumService forumService;

  @Before
  public void init() {
    try {
      transactionTemplate.execute(
          (TransactionCallback<Void>)
              transactionStatus -> {
                Tag hibernate = new Tag();
                hibernate.setName("hibernate");
                entityManager.persist(hibernate);

                Tag jpa = new Tag();
                jpa.setName("jpa");
                entityManager.persist(jpa);
                return null;
              });
    } catch (TransactionException e) {
      LOGGER.error("Failure", e);
    }
  }

  @Test
  public void test() {
    State.Holder.set(
        State.newBuilder()
            .withControllerName("ForumController")
            .withActionName("SavePost")
            .withFramework("spring")
            .build());

    SCHibernateWrapper.reset();

    Post newPost = forumService.newPost("High-Performance Java Persistence", "hibernate", "jpa");
    assertNotNull(newPost.getId());

    List<String> sqlStatements = SCHibernateWrapper.getAfterSqlStatements();
    assertEquals(5, sqlStatements.size());
    assertEquals(
        5,
        sqlStatements
            .stream()
            .filter(
                sql ->
                    sql.contains(
                        "/*action='SavePost',controller='ForumController',framework='spring'*/"))
            .count());

    SCHibernateWrapper.reset();

    List<Post> posts = forumService.findAllByTitle("High-Performance Java Persistence");
    assertEquals(1, posts.size());

    sqlStatements = SCHibernateWrapper.getAfterSqlStatements();
    assertEquals(1, sqlStatements.size());
    assertEquals(
        1,
        sqlStatements
            .stream()
            .filter(
                sql ->
                    sql.contains(
                        "/*action='SavePost',controller='ForumController',framework='spring'*/"))
            .count());

    SCHibernateWrapper.reset();

    Post post = forumService.findById(newPost.getId());
    assertEquals("High-Performance Java Persistence", post.getTitle());

    sqlStatements = SCHibernateWrapper.getAfterSqlStatements();
    assertEquals(1, sqlStatements.size());
    assertEquals(
        1,
        sqlStatements
            .stream()
            .filter(
                sql ->
                    sql.contains(
                        "/*action='SavePost',controller='ForumController',framework='spring'*/"))
            .count());
  }
}
