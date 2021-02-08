"""
Acrobits Web Services
"""

__version__ = '0.1.4'

from fastapi import (
    FastAPI,
    HTTPException,
    status,
)

import balance
import contacts
import healthcheck
import rate
import websvc


async def get_balance(params: websvc.Params) -> balance.Balance:
    raise HTTPException(status.HTTP_501_NOT_IMPLEMENTED)


async def get_contacts(params: contacts.Params) -> contacts.Contacts:
    raise HTTPException(status.HTTP_501_NOT_IMPLEMENTED)


async def get_rate(params: rate.Params) -> rate.Rate:
    raise HTTPException(status.HTTP_501_NOT_IMPLEMENTED)


def add_handlers(app: FastAPI):
    if not any((
        balance.add_handler(app, get_balance),
        contacts.add_handler(app, get_contacts),
        rate.add_handler(app, get_rate),
    )):
        raise RuntimeError('no enabled modules')

    healthcheck.add_handler(app)


app = FastAPI(title=__doc__.strip(), version=__version__)
websvc.wrap_http_exception(app)
add_handlers(app)
