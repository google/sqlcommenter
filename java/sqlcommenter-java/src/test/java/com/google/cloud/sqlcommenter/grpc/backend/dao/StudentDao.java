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

package com.google.cloud.sqlcommenter.grpc.backend.dao;

import com.google.cloud.sqlcommenter.grpc.backend.domain.Student;
import java.util.NoSuchElementException;
import javax.persistence.EntityManager;
import javax.persistence.EntityManagerFactory;
import javax.persistence.Persistence;

public class StudentDao {
  public Student findById(String studentId) {

    // We use entity managers to manage our two entities.
    // We use the factory design pattern to get the entity manager.
    // Here we should provide the name of the persistence unit that we provided in the
    // persistence.xml file.
    EntityManagerFactory emf = Persistence.createEntityManagerFactory("student-management-system");
    EntityManager em = emf.createEntityManager();

    // We can find a record in the database for a given id using the find method.
    // for the find method we have to provide our entity class and the id.
    Student student = em.find(Student.class, studentId);

    // If there is no record found with the provided student id, then we throw a NoSuchElement
    // exception.
    if (student == null) {
      throw new NoSuchElementException("NO DATA FOUND WITH THE ID " + studentId);
    }

    // If everything worked fine, return the result.
    return student;
  }
}
