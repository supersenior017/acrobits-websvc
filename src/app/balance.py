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
    Field,
)

import websvc


class Settings(BaseSettings):
    path: str = 'balance'
    currency: str = 'USD'
    enabled: bool = True

    class Config:
        env_prefix = 'balance_'


class Balance(BaseModel):
    balance: float
    currency: str = ''


class Response(Balance):
    """
    https://doc.acrobits.net/api/client/balance_checker.html#response
    """
    balance_string: str = Field(..., alias='balanceString')

    @classmethod
    def make(cls, balance: Balance, settings: Settings) -> Response:
        currency = balance.currency or settings.currency
        return cls(
            balance_string=f'{currency} {balance.balance:.2f}'.lstrip(),
            balance=balance.balance,
            currency=currency,
        )

    class Config:
        allow_population_by_field_name = True


def add_handler(
    app: FastAPI,
    fn: Callable[[websvc.Account], Awaitable[Balance]],
) -> bool:
    settings = Settings()
    if not settings.enabled:
        return False

    @app.get(f'/{settings.path}', response_model=Response)
    async def balance(account: websvc.Account = Depends()) -> Response:
        return Response.make(await fn(account), settings)

    return True
