acrobits-websvc
===============

`Acrobits Web Services <https://doc.acrobits.net/api/client/index.html>`_.
You should overwrite functions in `main.py` file.

.. code:: python

    async def get_balance(account: websvc.Account) -> balance.Balance:
        raise HTTPException(status.HTTP_503_SERVICE_UNAVAILABLE)

.. code:: python

    async def get_rate(params: rate.Params) -> rate.Rate:
        raise HTTPException(status.HTTP_503_SERVICE_UNAVAILABLE)

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

To run with docker:

.. code:: sh

    $ docker run --rm -d -p 8000:8000 acrobits-websvc

To run with docker-compose:

.. code:: sh

    $ docker-compose up -d
