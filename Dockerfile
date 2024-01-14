FROM golang:1.21 AS builder
ENV CGO_ENABLED 0
ADD . /app
WORKDIR /app
RUN go build -ldflags "-s -w" -v -o beerium .

FROM alpine:3
RUN apk update && \
    apk add openssl tzdata && \
    rm -rf /var/cache/apk/* \
    && mkdir /app

WORKDIR /app

ADD Dockerfile /Dockerfile

COPY --from=builder /app/beerium /app/beerium
ADD pkg/app/templates/homepage.gohtml /app/pkg/app/templates/homepage.gohtml

RUN chown nobody /app/beerium \
    && chmod 500 /app/beerium

USER nobody

EXPOSE 8080

ENTRYPOINT ["/app/beerium"]
