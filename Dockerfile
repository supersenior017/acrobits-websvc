FROM golang

WORKDIR /build
COPY ./src/ .
RUN CGO_ENABLED=0 go build -o ./main

FROM scratch

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /etc/passwd /etc/group /etc/
COPY --from=0 /build/main /

USER nobody:nogroup
STOPSIGNAL SIGINT

ENV port 8080
ENV ACROBITS_BALANCE_PATH /acrobits/balance
ENV ACROBITS_BALANCE_ADDR :$port
ENV ACROBITS_BALANCE_CURRENCY USD

ENTRYPOINT ["/main"]
EXPOSE $port
