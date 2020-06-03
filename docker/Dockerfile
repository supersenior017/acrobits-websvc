FROM golang

WORKDIR /build/
COPY ./src/ ./src/
COPY ./go.mod ./go.sum ./
RUN go mod download
RUN CGO_ENABLED=0 go build -o ./main ./src/

FROM scratch

COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /etc/passwd /etc/group /etc/
COPY --from=0 /build/main /

USER nobody:nogroup
STOPSIGNAL SIGINT

ENV port 8080
ENV ACROBITS_WEBSVC_ADDR :$port

ENTRYPOINT ["/main"]
EXPOSE $port
