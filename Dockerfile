# syntax=docker/dockerfile:1.7

FROM golang:1.26.2-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG SERVICE_PATH

RUN test -n "$SERVICE_PATH"
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -trimpath -ldflags="-s -w" -o /out/service "$SERVICE_PATH"

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates netcat-openbsd tzdata \
	&& addgroup -S app \
	&& adduser -S -G app app \
	&& mkdir -p /app/logs \
	&& chown -R app:app /app

COPY --from=builder /out/service /app/service

USER app

ENTRYPOINT ["/app/service"]
