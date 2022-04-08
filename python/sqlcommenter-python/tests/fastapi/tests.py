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

import json

import fastapi
import pytest
from starlette.testclient import TestClient

from google.cloud.sqlcommenter.fastapi import get_fastapi_info

from .app import app


@pytest.fixture
def client():
    client = TestClient(app)
    yield client


def test_get_fastapi_info_in_request_context(client):
    expected = {
        'app_name': 'SQLCommenter',
        'controller': 'fastapi_info',
        'framework': 'fastapi:%s' % fastapi.__version__,
        'route': '/fastapi-info',
    }
    resp = client.get('/fastapi-info')
    assert json.loads(resp.content.decode('utf-8')) == expected


def test_get_fastapi_info_in_404_error_context(client):
    expected = {
        'app_name': 'SQLCommenter',
        'framework': 'fastapi:%s' % fastapi.__version__,
    }
    resp = client.get('/doesnt-exist')
    assert json.loads(resp.content.decode('utf-8')) == expected


def test_get_fastapi_info_outside_request_context(client):
    assert get_fastapi_info() == {}
