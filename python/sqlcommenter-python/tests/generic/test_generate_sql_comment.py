#!/usr/bin/python
#
# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from unittest import TestCase

from google.cloud.sqlcommenter import generate_sql_comment


class GenerateSqlCommentTests(TestCase):
    def test_no_args(self):
        self.assertEqual(generate_sql_comment(), '')

    def test_end_comment_escaping(self):
        self.assertIn("a='%%2A/abc'", generate_sql_comment(a='*/abc'))

    def test_bytes(self):
        self.assertIn("a='%%2A/abc'", generate_sql_comment(a=b'*/abc'))

    def test_url_quoting(self):
        self.assertIn("foo='bar%%2Cbaz'", generate_sql_comment(foo='bar,baz'))
