FROM python:3.8

WORKDIR /usr/src/

COPY ./src/ ./app/
COPY ./requirements.txt ./
RUN pip install --no-cache-dir -r ./requirements.txt

USER nobody:nogroup
STOPSIGNAL SIGINT

ENV HEALTHCHECK_ENABLE 1
HEALTHCHECK CMD curl -f http://localhost:8000/ping || exit 1

ENTRYPOINT [                    \
    "uvicorn",     "main:app",  \
    "--app-dir",   "./app/",    \
    "--log-level", "warning",   \
    "--host",      "0.0.0.0"    \
]
EXPOSE 8000
