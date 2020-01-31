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

package com.google.cloud.sqlcommenter.spring.backend.service;

import com.google.cloud.sqlcommenter.spring.backend.dao.PostRepository;
import com.google.cloud.sqlcommenter.spring.backend.dao.TagRepository;
import com.google.cloud.sqlcommenter.spring.backend.domain.Post;
import java.util.Arrays;
import java.util.List;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class ForumServiceImpl implements ForumService {

  @Autowired private PostRepository postRepository;

  @Autowired private TagRepository tagRepository;

  @Override
  @Transactional
  public Post newPost(String title, String... tags) {
    Post post = new Post();
    post.setTitle(title);
    post.getTags().addAll(tagRepository.findByNameIn(Arrays.asList(tags)));
    return postRepository.save(post);
  }

  @Override
  @Transactional(readOnly = true)
  public List<Post> findAllByTitle(String title) {
    return postRepository.findByTitle(title);
  }

  @Override
  @Transactional
  public Post findById(Long id) {
    return postRepository.findById(id).orElse(null);
  }
}
