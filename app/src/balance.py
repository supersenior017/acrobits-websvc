"""
https://doc.acrobits.net/api/client/balance_checker.html
"""

from __future__ import annotations

from typing import (
    Awaitable,
    Callable,
)

from fastapi import (
    Depends,
    FastAPI,
)
from pydantic import (
    BaseModel,
    BaseSettings,
)

import utils
import websvc


class Settings(BaseSettings):
    currency: str = 'USD'
    enable: bool = True

    class Config:
        env_prefix = 'balance_'


settings = Settings()


class Balance(BaseModel):
    balance: float
    currency: str = ''


class Response(Balance):
    """
    https://doc.acrobits.net/api/client/balance_checker.html#response
    """
    balance_string: str

    @classmethod
    def make(cls, balance: Balance) -> Response:
        currency = balance.currency or settings.currency
        return cls(
            balance_string=f'{currency} {balance.balance:.2f}'.lstrip(),
            balance=balance.balance,
            currency=currency,
        )

    class Config:
        allow_population_by_field_name = True
        alias_generator = utils.to_lower_camel_case


def add_handler(
    app: FastAPI,
    fn: Callable[[websvc.Params], Awaitable[Balance]],
) -> bool:
    if not settings.enable:
        return False

    @app.get('/balance', response_model=Response)
    async def balance(params: websvc.Params = Depends()) -> Response:
        return Response.make(await fn(params))

    return True
