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

from __future__ import absolute_import

try:
    import flask
    from flask import request
except ImportError:
    flask = None


def get_flask_info():
    """
    Get available info from the current flask request, if we're in a
    flask request-response cycle. Else, return an empty dict.
    """
    info = {}
    # https://flask.palletsprojects.com/en/1.1.x/api/#flask.has_request_context
    if flask and request:
        info['framework'] = 'flask:%s' % flask.__version__
        if request.endpoint:
            info['controller'] = request.endpoint
        if request.url_rule and request.url_rule.rule:
            info['route'] = request.url_rule.rule
    return info
