#!/usr/bin/python
#
# Copyright 2022 Google LLC
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
    from typing import Optional
    from asgiref.compatibility import guarantee_single_callable
    from contextvars import ContextVar
    import fastapi
    from fastapi import FastAPI
    from starlette.routing import Match, Route
except ImportError:
    fastapi = None

context = ContextVar("context", default={})


def get_fastapi_info():
    """
    Get available info from the current FastAPI request, if we're in a
    FastAPI request-response cycle. Else, return an empty dict.
    """
    info = {}
    if fastapi and context:
        info = context.get()
    return info


class SQLCommenterMiddleware:
    """The ASGI application middleware.
    This class is an ASGI middleware that augment SQL statements before execution,
     with comments containing information about the code that caused its execution.

    Args:
        app: The ASGI application callable to forward requests to.
    """

    def __init__(self, app):
        self.app = guarantee_single_callable(app)

    async def __call__(self, scope, receive, send):
        """The ASGI application
        Args:
            scope: An ASGI environment.
            receive: An awaitable callable yielding dictionaries
            send: An awaitable callable taking a single dictionary as argument.
        """
        if scope["type"] not in ("http", "websocket"):
            return await self.app(scope, receive, send)

        if not isinstance(scope["app"], FastAPI):
            return await self.app(scope, receive, send)

        fastapi_app = scope["app"]
        info = _get_fastapi_info(fastapi_app, scope)
        token = context.set(info)

        try:
            await self.app(scope, receive, send)
        finally:
            context.reset(token)


def _get_fastapi_info(fastapi_app: FastAPI, scope) -> dict:
    info = {
        "framework": 'fastapi:%s' % fastapi.__version__,
        "app_name": fastapi_app.title,
    }

    route = _get_fastapi_route(fastapi_app, scope)
    if route:
        info["controller"] = route.name
        info["route"] = route.path

    return info


def _get_fastapi_route(fastapi_app: FastAPI, scope) -> Optional[Route]:
    for route in fastapi_app.router.routes:
        # Determine if any route matches the incoming scope,
        # and return the route name if found.
        match, child_scope = route.matches(scope)
        if match == Match.FULL:
            return child_scope["route"]
    return None
