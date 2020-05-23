FROM golang

WORKDIR /build
COPY ./src/ .
RUN CGO_ENABLED=0 go build -o ./main

FROM scratch

COPY --from=0 /etc/ssl/certs/ /etc/ssl/certs/
COPY --from=0 /etc/passwd /etc/passwd
COPY --from=0 /etc/group /etc/group
USER nobody:nogroup
STOPSIGNAL SIGINT

COPY --from=0 /build/main /
COPY ./config/acrobits-balance.json /config.json

ENV port 8080
ENV ACROBITS_BALANCE_PATH /acrobits/balance
ENV ACROBITS_BALANCE_ADDR :$port
ENV ACROBITS_BALANCE_CURRENCY USD

ENTRYPOINT ["/main"]
EXPOSE $port
