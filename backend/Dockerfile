# syntax=docker/dockerfile:1
FROM golang:1.19.7-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy && \
    go install -v ./... && \
    go build && \
    go clean -modcache

FROM alpine:latest

RUN addgroup -g 1001 docker && \
    adduser -D -u 1001 -G docker docker


WORKDIR /app
COPY --from=builder /app ./
RUN chown -R docker:docker /app

USER docker
EXPOSE 8080
CMD ["./ariadne", "server"]
