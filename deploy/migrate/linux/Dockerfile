# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -ldflags="-s -w" -o main /app/cmd/migrate/main.go

# Run stage
FROM alpine
WORKDIR /app
COPY --from=builder /app/main /app/cmd/migrate/
COPY config config
COPY migrations migrations

CMD ["/app/cmd/migrate/main", "-up"]