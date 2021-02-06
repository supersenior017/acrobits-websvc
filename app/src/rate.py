"""
https://doc.acrobits.net/api/client/rate_checker.html
"""

from __future__ import annotations

from typing import (
    Awaitable,
    Callable,
    Optional,
    Union,
)

from fastapi import (
    Depends,
    FastAPI,
    HTTPException,
    Query,
    status,
)
from pydantic import (
    BaseModel,
    BaseSettings,
)

import utils
import websvc


class Settings(BaseSettings):
    currency: str = '¢'
    specification: str = 'min.'
    enable: bool = True

    class Config:
        env_prefix = 'rate_'


settings = Settings()


class Params(websvc.Params):
    """
    https://doc.acrobits.net/api/client/rate_checker.html#parameters
    """
    target_number: Optional[str] = Query(None, alias='targetNumber')
    smart_uri: Optional[str] = Query(None, alias='smartUri')


class Call(BaseModel):
    price: float
    specification: str = ''


class Rate(BaseModel):
    call: Optional[Call] = None
    message: Optional[float] = None
    currency: str = ''


def _format_response(rate: Rate) -> (str, str):
    call = rate.call
    message = rate.message
    currency = rate.currency or settings.currency
    return (
        '' if call is None else (
            f'{call.price:g}{currency} '
            f'{call.specification or settings.specification}'
        ).rstrip(),
        '' if message is None else f'{message:g}{currency}',
    )


class ResponseNumber(BaseModel):
    """
    https://doc.acrobits.net/api/client/rate_checker.html#response
    """
    call_rate_string: str
    message_rate_string: str

    @classmethod
    def make(cls, rate: Rate) -> ResponseNumber:
        call, message = _format_response(rate)
        return cls(call_rate_string=call, message_rate_string=message)

    class Config:
        allow_population_by_field_name = True
        alias_generator = utils.to_lower_camel_case


class ResponseUri(BaseModel):
    """
    https://doc.acrobits.net/api/client/rate_checker.html#response
    """
    smart_call_rate_string: str
    smart_message_rate_string: str

    @classmethod
    def make(cls, rate: Rate) -> ResponseUri:
        call, message = _format_response(rate)
        return cls(
            smart_call_rate_string=call,
            smart_message_rate_string=message,
        )

    class Config:
        allow_population_by_field_name = True
        alias_generator = utils.to_lower_camel_case


def add_handler(app: FastAPI, fn: Callable[[Params], Awaitable[Rate]]) -> bool:
    if not settings.enable:
        return False

    @app.get('/rate', response_model=Union[ResponseNumber, ResponseUri])
    async def rate(
        params: Params = Depends(),
    ) -> Union[ResponseNumber, ResponseUri]:
        if params.target_number:
            resp = ResponseNumber
        elif params.smart_uri:
            resp = ResponseUri
        else:
            raise HTTPException(status.HTTP_400_BAD_REQUEST)
        return resp.make(await fn(params))

    return True
