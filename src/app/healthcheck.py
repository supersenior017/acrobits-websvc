from fastapi import FastAPI
from fastapi.responses import PlainTextResponse
from pydantic import BaseSettings


class Settings(BaseSettings):
    path: str = 'ping'
    enabled: bool = False

    class Config:
        env_prefix = 'healthcheck_'


def add_handler(app: FastAPI) -> bool:
    settings = Settings()
    if not settings.enabled:
        return False

    @app.get(f'/{settings.path}', include_in_schema=False)
    async def handler() -> PlainTextResponse:
        return PlainTextResponse('OK')

    return True
