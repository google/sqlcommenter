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

from sqlcommenter import generate_sql_comment

from .opencensus_mock import mock_opencensus_tracer


class GenerateSqlCommentTests(TestCase):
    def test_no_args(self):
        # As per internal issue #28, we've commented out
        # the old defaults such as os.uname() information.
        self.assertEqual(generate_sql_comment(), '')

    def test_end_comment_escaping(self):
        self.assertIn(r'%%2A/abc', generate_sql_comment(a='*/abc'))

    def test_url_quoting(self):
        self.assertIn("foo='bar%%2Cbaz'", generate_sql_comment(foo='bar,baz'))

    def test_opencensus_not_available(self):
        self.assertEqual(generate_sql_comment(with_opencensus=True), '')
        self.assertEqual(generate_sql_comment(with_opencensus=False), '')

    def test_opencensus_disabled(self):
        with mock_opencensus_tracer():
            self.assertEqual(generate_sql_comment(with_opencensus=False), '')

    def test_opencensus_enabled(self):
        with mock_opencensus_tracer():
            self.assertEqual(
                generate_sql_comment(with_opencensus=True),
                " /*traceparent='00-trace%%20id-span%%20id-00',"
                "tracestate='congo%%3Dt61rcWkgMzE%%2Crojo%%3D00f067aa0ba902b7'*/"
            )
