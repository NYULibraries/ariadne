# syntax=docker/dockerfile:1
FROM golang:1.19.7-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy && \
    go install -v ./... && \
    go build && \
    go clean -modcache

FROM alpine:latest

RUN addgroup --gid 1001 docker && \
    adduser --uid 1001 --ingroup docker --disabled-password --gecos "" docker

WORKDIR /app
COPY --from=builder /app ./

USER docker
EXPOSE 8080
CMD ["./ariadne", "server"]
