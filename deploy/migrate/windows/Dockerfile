# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY ../../.. .
RUN go build -o main /app/cmd/migrate/main.go

# Run stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/main /app/cmd/migrate/
COPY ../../../config config
COPY ../../../migrations migrations

CMD ["/app/cmd/migrate/main", "-up"]