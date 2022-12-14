# syntax=docker/dockerfile:1
FROM golang:1.17-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && \
    go install -v ./... && \
    go build && \
    go clean -modcache

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app ./
EXPOSE 8080
CMD ["./ariadne", "server"]
