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

from flask import Flask
from google.cloud.sqlcommenter.flask import get_flask_info

app = Flask(__name__)


@app.route('/flask-info')
def flask_info():
    return get_flask_info()


@app.errorhandler(404)
def handle_bad_request(e):
    return get_flask_info(), 404
