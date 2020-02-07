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

import flask
import pytest
from google.cloud.sqlcommenter.flask import get_flask_info

from .app import app

# Following the example provided at
# https://flask.palletsprojects.com/en/1.1.x/testing/#the-testing-skeleton


@pytest.fixture
def client():
    app.config['TESTING'] = True
    with app.test_client() as client:
        yield client


def test_get_flask_info_in_request_context(client):
    expected = {
        'controller': 'flask_info',
        'framework': 'flask:%s' % flask.__version__,
        'route': '/flask-info',
    }
    resp = client.get('/flask-info')
    assert json.loads(resp.data.decode('utf-8')) == expected


def test_get_flask_info_in_404_error_context(client):
    expected = {
        'framework': 'flask:%s' % flask.__version__,
    }
    resp = client.get('/doesnt-exist')
    assert json.loads(resp.data.decode('utf-8')) == expected


def test_get_flask_info_outside_request_context(client):
    assert get_flask_info() == {}
