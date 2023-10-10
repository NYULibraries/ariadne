# syntax=docker/dockerfile:1
FROM golang:1.19.7-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod tidy && \
    go install -v ./... && \
    go build && \
    go clean -modcache

FROM alpine:latest

RUN addgroup --gid 1000 docker && \
    adduser --uid 1000 --ingroup docker --disabled-password --gecos "" docker

WORKDIR /app
COPY --from=builder /app ./
RUN chown -R docker:docker /app

USER docker
EXPOSE 8080
CMD ["./ariadne", "server"]
