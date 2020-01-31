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

from setuptools import find_packages, setup

setup(
    name='sqlcommenter',
    version='0.1',
    author='Google Developers',
    author_email='sqlcommenter@googlegroups.com',
    description=('Augment SQL statements with meta information about frameworks and the running environment.'),
    license='BSD',
    packages=find_packages(exclude=['tests']),
    extras_require={
        'django': ['django >= 1.11'],
        'flask': ['flask'],
        'psycopg2': ['psycopg2'],
        'sqlalchemy': ['sqlalchemy'],
    },
    classifiers=[
        'Development Status :: 4 - Beta',
        'Environment :: Web Environment',
        'Intended Audience :: Developers',
        'License :: OSI Approved :: BSD License',
        'Operating System :: OS Independent',
        'Programming Language :: Python :: 3',
        'Programming Language :: Python :: 3.5',
        'Programming Language :: Python :: 3.6',
        'Programming Language :: Python :: 3.7',
        'Topic :: Utilities',
        'Framework :: Django',
        'Framework :: Django :: 2.1',
        'Framework :: Django :: 2.2',
        'Framework :: SQLAlchemy',
        'Framework :: Psycopg2',
    ],
)
