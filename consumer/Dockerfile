
FROM golang:tip-alpine3.23 AS builder

RUN apk add --no-cache build-base

WORKDIR /app

ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE -D_FILE_OFFSET_BITS=64"

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /app/consumer ./cmd/

FROM alpine:3.23

RUN apk add --no-cache sqlite-libs

WORKDIR /app

COPY --from=builder /app/consumer .

ENTRYPOINT ["./consumer"]