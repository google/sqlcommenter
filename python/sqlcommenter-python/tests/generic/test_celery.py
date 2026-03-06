#!/usr/bin/python
#
# Copyright 2026
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

import unittest

import google.cloud.sqlcommenter.celery as mod


class CeleryModuleTests(unittest.TestCase):
    def tearDown(self) -> None:
        try:
            mod._context.set({})
        except Exception:
            pass

    def test_get_celery_info_empty(self):
        self.assertEqual(mod.get_celery_info(), {})

    def test_get_celery_info_after_context_set(self):
        token = mod._context.set({
            'framework': 'celery:5.3.0',
            'task': 'tasks.add',
            'route': 'celery',
        })
        try:
            self.assertEqual(
                mod.get_celery_info(),
                {'framework': 'celery:5.3.0', 'task': 'tasks.add', 'route': 'celery'}
            )
        finally:
            mod._context.reset(token)

    def test_install_signals_without_celery_raises(self):
        original_celery, original_signals = mod.celery, mod.signals
        try:
            mod.celery = None
            mod.signals = None
            with self.assertRaises(ImportError):
                mod.install_signals()
        finally:
            mod.celery, mod.signals = original_celery, original_signals

    def test_signal_flow_sets_and_clears_context(self):
        class _SigList:
            def __init__(self):
                self._cbs = []
            def connect(self, cb, weak=False):  # noqa: ARG002 - weak unused
                self._cbs.append(cb)
                return cb

        class FakeSignals:
            def __init__(self):
                self.task_prerun = _SigList()
                self.task_postrun = _SigList()

        class FakeCelery:
            __version__ = '9.9.9'

        original_celery, original_signals = mod.celery, mod.signals
        try:
            mod.celery = FakeCelery()
            mod.signals = FakeSignals()
            mod.install_signals()

            self.assertEqual(len(mod.signals.task_prerun._cbs), 1)
            self.assertEqual(len(mod.signals.task_postrun._cbs), 1)

            prerun_cb = mod.signals.task_prerun._cbs[0]
            postrun_cb = mod.signals.task_postrun._cbs[0]

            class FakeReq:
                routing_key = 'celery'
            class FakeTask:
                name = 'tasks.add'
                request = FakeReq()

            t = FakeTask()
            prerun_cb(task=t)
            self.assertEqual(
                mod.get_celery_info(),
                {'framework': 'celery:9.9.9', 'task': 'tasks.add', 'route': 'celery'}
            )

            postrun_cb(task=t)
            self.assertEqual(mod.get_celery_info(), {})
        finally:
            mod.celery, mod.signals = original_celery, original_signals