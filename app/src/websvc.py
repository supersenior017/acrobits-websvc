from typing import Optional

from fastapi import (
    FastAPI,
    Query,
)
from fastapi.responses import JSONResponse
from pydantic import BaseModel
from starlette.exceptions import HTTPException as StarletteHTTPException


def wrap_http_exception(app: FastAPI):
    """
    https://doc.acrobits.net/api/client/intro.html#web-service-responses
    """
    @app.exception_handler(StarletteHTTPException)
    async def http_exception_handler(request, exc):
        return JSONResponse({'message': exc.detail}, exc.status_code)


class Account(BaseModel):
    """
    https://doc.acrobits.net/api/client/intro.html#account-parameters
    """
    username: Optional[str] = Query(None)
    password: Optional[str] = Query(None)


class Params(Account):
    nonce: Optional[str] = Query(None)
    user: Optional[str] = Query(None)
