from typing import Optional

from fastapi import FastAPI, status
from fastapi.responses import JSONResponse
from google.cloud.sqlcommenter.fastapi import (
    SQLCommenterMiddleware, get_fastapi_info,
)
from starlette.applications import Starlette
from starlette.exceptions import HTTPException as StarletteHTTPException
from starlette.routing import Route

app = FastAPI(title="SQLCommenter")

app.add_middleware(SQLCommenterMiddleware)


@app.get("/fastapi-info")
def fastapi_info():
    return get_fastapi_info()


@app.get("/items/{item_id}")
def read_item(item_id: int, q: Optional[str] = None):
    return get_fastapi_info()


@app.exception_handler(StarletteHTTPException)
async def custom_http_exception_handler(request, exc):
    return JSONResponse(
        status_code=status.HTTP_404_NOT_FOUND,
        content=get_fastapi_info(),
    )


def starlette_endpoint(_):
    return JSONResponse({"from": "starlette"})


starlette_subapi = Starlette(routes=[
    Route("/", starlette_endpoint),
])


app.mount("/starlette", starlette_subapi)
