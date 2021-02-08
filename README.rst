acrobits-websvc
===============

`Acrobits Web Services <https://doc.acrobits.net/api/client/index.html>`_.
You have to overwrite functions in `main.py` file.

.. code:: python

    async def get_balance(params: websvc.Params) -> balance.Balance:
        raise HTTPException(status.HTTP_501_NOT_IMPLEMENTED)

.. code:: python

    async def get_contacts(params: contacts.Params) -> contacts.Contacts:
        raise HTTPException(status.HTTP_501_NOT_IMPLEMENTED)

.. code:: python

    async def get_rate(params: rate.Params) -> rate.Rate:
        raise HTTPException(status.HTTP_501_NOT_IMPLEMENTED)

Installation
------------

.. code:: sh

    $ make

or

.. code:: sh

    $ make docker

Example
-------

To run localhost:

.. code:: sh

    $ ./run.sh

To run in docker:

.. code:: sh

    $ docker-compose up -d
