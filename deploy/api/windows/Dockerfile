# Build stage
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY ../../.. .
RUN apk update && \
    apk add --no-cache gcc && \
    rm -rf /var/lib/apt/lists/* && \
    apk add musl-dev
RUN go build -o main /app/cmd/api/main.go

# Run stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/main /app/cmd/api/
COPY ../../../env env

EXPOSE 5487
CMD ["/app/cmd/api/main"]