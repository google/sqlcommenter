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

package com.google.cloud.sqlcommenter.util;

import com.google.cloud.sqlcommenter.schibernate.SCHibernate;
import java.util.ArrayList;
import java.util.List;

/**
 * The {@link SCHibernateWrapper} class wraps the {@link
 * com.google.cloud.sqlcommenter.schibernate.SCHibernate} class so that we can also get a chance to
 * intercept the executed statements and assert them during testing.
 */
public class SCHibernateWrapper extends SCHibernate {

  private static List<String> beforeSqlStatements = new ArrayList<>();

  private static List<String> afterSqlStatements = new ArrayList<>();

  public static List<String> getBeforeSqlStatements() {
    return beforeSqlStatements;
  }

  public static List<String> getAfterSqlStatements() {
    return afterSqlStatements;
  }

  public static void reset() {
    beforeSqlStatements.clear();
    afterSqlStatements.clear();
  }

  @Override
  public String inspect(String beforeSql) {
    beforeSqlStatements.add(beforeSql);
    String afterSql = super.inspect(beforeSql);
    afterSqlStatements.add(afterSql);
    return afterSql;
  }
}
