"""
Acrobits Web Services
"""

__version__ = '0.1.2'

from fastapi import (
    FastAPI,
    HTTPException,
    status,
)

import balance
import healthcheck
import rate
import websvc


async def get_balance(account: websvc.Account) -> balance.Balance:
    raise HTTPException(status.HTTP_503_SERVICE_UNAVAILABLE)


async def get_rate(params: rate.Params) -> rate.Rate:
    raise HTTPException(status.HTTP_503_SERVICE_UNAVAILABLE)


def add_handlers(app: FastAPI):
    if not any((
        balance.add_handler(app, get_balance),
        rate.add_handler(app, get_rate),
    )):
        raise RuntimeError('No enabled modules')

    healthcheck.add_handler(app)


app = FastAPI(title=__doc__.strip(), version=__version__)
websvc.wrap_http_exception(app)
add_handlers(app)
