# syntax=docker/dockerfile:1
FROM golang:1-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && \
    go install -v ./... && \
    go build && \
    go clean -modcache

FROM alpine:latest  
WORKDIR /go/src/app
COPY --from=builder /app ./
CMD ["./getit"]  
