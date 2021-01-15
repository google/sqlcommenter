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

"""Python 2 compatibility shims."""

import sys
import unittest

try:
    from unittest import mock
except ImportError:  # Python 2
    import mock


def skipIfPy2(testcase):
    return unittest.skipIf(
        sys.version_info.major == 2, "Feature only support in python3+"
    )(testcase)


__all__ = ["mock"]
