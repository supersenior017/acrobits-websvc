"""
https://doc.acrobits.net/api/client/web_contacts.html
"""

from __future__ import annotations

import enum
from email.utils import parsedate_to_datetime
from typing import (
    Awaitable,
    Callable,
    List,
    Optional,
    Union,
)

from fastapi import (
    Depends,
    FastAPI,
    Header,
    Response as FastAPIResponse,
    status,
)
from pydantic import (
    BaseModel,
    BaseSettings,
)

import utils
import websvc


class Settings(BaseSettings):
    enable: bool = True

    class Config:
        env_prefix = 'contacts_'


settings = Settings()


class ContactEntryType(str, enum.Enum):
    TEL = 'tel'
    EMAIL = 'email'
    URL = 'url'


class ContactEntry(BaseModel):
    entry_id: str
    type: ContactEntryType
    label: Optional[str]
    uri: str

    class Config:
        allow_population_by_field_name = True
        alias_generator = utils.to_lower_camel_case


class Contact(BaseModel):
    """
    https://doc.acrobits.net/cloudsoftphone/contacts.html
    """
    contact_id: str
    display_name: str
    checksum: Optional[str]
    fname: Optional[str]
    mname: Optional[str]
    lname: Optional[str]
    fname_phonetic: Optional[str]
    mname_phonetic: Optional[str]
    lname_phonetic: Optional[str]
    nick: Optional[str]
    name_prefix: Optional[str]
    name_suffix: Optional[str]
    company: Optional[str]
    department_name: Optional[str]
    job_title: Optional[str]
    birthday: Optional[str]
    street: Optional[str]
    city: Optional[str]
    state: Optional[str]
    zip: Optional[str]
    country: Optional[str]
    country_code: Optional[str]
    notes: Optional[str]
    contact_entries: List[ContactEntry]

    avatar: Optional[str]
    large_avatar: Optional[str]

    class Config:
        allow_population_by_field_name = True
        alias_generator = utils.to_lower_camel_case


class Contacts(BaseModel):
    contacts: List[Contact]
    last_modified: Optional[str] = None


class Params(websvc.Params):
    if_modified_since: Optional[str] = None


class Response(BaseModel):
    """
    https://doc.acrobits.net/api/client/web_contacts.html#response
    """
    contacts: List[Contact]

    @classmethod
    def make(cls, contacts: Contacts) -> Response:
        return cls(contacts=contacts.contacts)


def add_handler(
    app: FastAPI,
    fn: Callable[[Params], Awaitable[Contacts]],
) -> bool:
    if not settings.enable:
        return False

    @app.get(
        '/contacts',
        response_model=Optional[Response],
        response_model_exclude_none=True,
    )
    async def contacts(
        response: FastAPIResponse,
        params: websvc.Params = Depends(),
        if_modified_since: Optional[str] = Header(None),
    ) -> Union[Response, FastAPIResponse]:
        params1 = Params(**params.__dict__)
        params1.if_modified_since = if_modified_since
        contacts = await fn(params1)
        if contacts.last_modified is not None:
            if (if_modified_since is not None and (
                parsedate_to_datetime(if_modified_since) >=
                parsedate_to_datetime(contacts.last_modified)
            )):
                response.status_code = status.HTTP_304_NOT_MODIFIED
                return response
            response.headers['Last-Modified'] = contacts.last_modified
        return Response.make(contacts)

    return True
