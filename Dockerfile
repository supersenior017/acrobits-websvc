FROM golang

WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 go build -o ./target/main ./src

FROM scratch

ARG port=8080
ENV ACROBITS_BALANCE_PATH /acrobits/balance
ENV ACROBITS_BALANCE_ADDR :$port
ENV ACROBITS_BALANCE_CURRENCY USD

COPY --from=0 /etc/passwd /etc/passwd
COPY --from=0 /etc/group /etc/group
USER nobody:nogroup

COPY --from=0 /build/target/main /
COPY ./config/acrobits-balance.json /config.json

ENTRYPOINT ["/main"]
EXPOSE $port
