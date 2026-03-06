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


from contextvars import ContextVar
from typing import Any, Dict, Optional


try:
    import celery
    from celery import signals
except Exception:
    celery = None
    signals = None


_context: ContextVar[Dict[str, Any]] = ContextVar("sqlcommenter_celery_context", default={})


def get_celery_info() -> Dict[str, Any]:
    info = _context.get() or {}
    return dict(info) if info else {}


def install_signals() -> None:
    if celery is None or signals is None:
        raise ImportError("celery is not installed.")

    def _on_task_prerun(sender=None, task_id: Optional[str]=None, task=None, args=None, kwargs=None, **kw):
        fw = f"celery:{getattr(celery, '__version__', 'unknown')}"
        info: Dict[str, Any] = {"framework": fw}
        t = task if task is not None else sender
        try:
            if t is not None:
                name = getattr(t, "name", None)
                if not name:
                    req = getattr(t, "request", None)
                    name = getattr(req, "task", None)
                if name:
                    info["task"] = name
                req = getattr(t, "request", None)
                if req is not None:
                    rk = getattr(req, "routing_key", None)
                    if not rk:
                        di = getattr(req, "delivery_info", None) or {}
                        rk = di.get("routing_key") if isinstance(di, dict) else None
                    if rk:
                        info["route"] = rk
        except Exception:
            pass
        token = _context.set(info)
        if t is not None:
            try:
                setattr(t, "_sqlcommenter_token", token)
            except Exception:
                pass

    def _on_task_postrun(sender=None, task_id: Optional[str]=None, task=None, args=None, kwargs=None, **kw):
        t = task if task is not None else sender
        token = None
        if t is not None:
            token = getattr(t, "_sqlcommenter_token", None)
        try:
            if token is not None:
                _context.reset(token)
                try:
                    delattr(t, "_sqlcommenter_token")
                except Exception:
                    pass
            else:
                _context.set({})
        except Exception: 
            pass

    signals.task_prerun.connect(_on_task_prerun, weak=False)
    signals.task_postrun.connect(_on_task_postrun, weak=False)
