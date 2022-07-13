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

DATABASES = {
    'default': {
        'ENGINE': 'django.db.backends.sqlite3',
    },

    'other': {
        'ENGINE': 'django.db.backends.sqlite3',
    },
}

INSTALLED_APPS = ['tests.django']

ROOT_URLCONF = 'tests.django.urls'

SECRET_KEY = 'commenter_tests_secret_key'
